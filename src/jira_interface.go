package main

import (
	"fmt"
	"io"
    "net/http"
	"strings"
	"time"
	"net/url"
	"encoding/json"
)

type IssueFilter struct {
	user        string
	sprint_name string
	status      string
	reporter    string
}

func jira_get(endpoint string) ([]byte, error) {
    url := strings.TrimRight(global_config.URL, "/") + endpoint

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    req.SetBasicAuth(global_config.email, global_config.token)
    req.Header.Set("Accept", "application/json")

    client := &http.Client{
        Timeout: 15 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf(
            "Jira API error: %s: %s",
            resp.Status,
            string(body),
        )
    }

    return body, nil
}

func merge_jql_to_endpoint(endpoint string, jql_str string) string {
	jql := url.QueryEscape(jql_str)
    return fmt.Sprintf("%s?jql=%s", endpoint, jql)
}

func get_sprint_tasks_nr(sprintID int) ([]byte, error) {
    endpoint := fmt.Sprintf("/rest/software/1.0/sprint/%d/issue", sprintID)
    endpoint = merge_jql_to_endpoint(endpoint, "assignee = currentUser()")

    return jira_get(endpoint)
}

func get_sprint_tasks_name(sprintName string) ([]byte, error) {
    jql := fmt.Sprintf(
        `sprint = "%s" AND assignee = currentUser()`,
        sprintName,
    )

    endpoint := merge_jql_to_endpoint("/rest/api/3/search/jql", jql)
    return jira_get(endpoint)
}

func get_issue(id string) (JiraIssue, error) {
    var issue JiraIssue

    data, err := jira_get("/rest/api/3/issue/" + id)
    if err != nil {
        return issue, err
    }

    err = json.Unmarshal(data, &issue)
    return issue, err
}

type ADFNode struct {
    Type    string    `json:"type"`
    Text    string    `json:"text"`
    Content []ADFNode `json:"content"`

    Attrs struct {
        State string `json:"state"`
        URL   string `json:"url"`
        ID    string `json:"id"`
        Type  string `json:"type"`
    } `json:"attrs"`
}

/*
func adfToText(data json.RawMessage) (string, error) {
    var root ADFNode
    if err := json.Unmarshal(data, &root); err != nil {
        return "", err
    }

    var b strings.Builder

    var walk func(ADFNode)
    walk = func(n ADFNode) {
        if n.Type == "text" {
            b.WriteString(n.Text)
        } else if n.Type == "hardBreak" {
            b.WriteByte('\n')
        }

        for _, child := range n.Content {
            walk(child)
        }

        if n.Type == "paragraph" || n.Type == "heading" {
            b.WriteByte('\n')
        }
    }

    walk(root)
    return strings.TrimRight(b.String(), "\n"), nil
}
*/
func adfToText(data json.RawMessage) (string, error) {
    var root ADFNode
    if err := json.Unmarshal(data, &root); err != nil {
        return "", err
    }

    var b strings.Builder

    var walk func(ADFNode)
    walk = func(n ADFNode) {
        switch n.Type {
        case "text":
            b.WriteString(n.Text)

        case "hardBreak":
            b.WriteByte('\n')

        case "taskItem":
            if n.Attrs.State == "DONE" {
                b.WriteString("- [x] - ")
            } else {
                b.WriteString("- [ ] - ")
            }
		case "media", "inlineCard":
			if n.Attrs.URL != "" {
				b.WriteString(n.Attrs.URL)
			} else if n.Type == "media" {
				b.WriteString("[image: " + n.Attrs.ID + "]")
			}
        }

        for _, child := range n.Content {
            walk(child)
        }

        switch n.Type {
        case "paragraph", "heading", "taskItem":
            if !strings.HasSuffix(b.String(), "\n") {
                b.WriteByte('\n')
            }
        }
    }

    walk(root)
    return strings.TrimRight(b.String(), "\n"), nil
}

func get_all_issue_sprint_name(sprint_name string) ([]JiraIssue, error) {
    jql := fmt.Sprintf(`sprint = "%s"`, sprint_name)
    endpoint := merge_jql_to_endpoint("/rest/api/3/search/jql", jql)

	endpoint += "&fields=summary,description,status,assignee,reporter,parent,fixVersions,priority,issuetype&maxResults=200"

    var issues []JiraIssue
    var token string

    for {
        pageURL := endpoint
        if token != "" {
            pageURL += "&nextPageToken=" + url.QueryEscape(token)
        }

        data, err := jira_get(pageURL)
        if err != nil {
            return nil, err
        }

        var response struct {
            Issues        []JiraIssue `json:"issues"`
            NextPageToken string      `json:"nextPageToken"`
        }

        if err := json.Unmarshal(data, &response); err != nil {
            return nil, err
        }

        issues = append(issues, response.Issues...)
        token = response.NextPageToken

        if token == "" {
            break
        }
    }

    return issues, nil
}

func get_all_issue_filter(filter IssueFilter) ([]JiraIssue, error) {
    var conditions []string

    // Escape values used in JQL strings
    quote := func(s string) string {
        s = strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s)
        return `"` + s + `"`
    }

	// reporter filter
	if filter.reporter == "" || filter.reporter == "all" {
		//
	} else {
		conditions = append(conditions, "reporter = "+quote(filter.reporter))
	}

    // Sprint filter
    if filter.sprint_name != "" {
        conditions = append(conditions, "sprint = "+quote(filter.sprint_name))
    } else {
        conditions = append(conditions, "sprint in openSprints()")
    }

    // User filter
    if filter.user == "all" {
		//
	} else if filter.user == "me" || filter.user == "" {
        conditions = append(conditions, "assignee = currentUser()")
    } else if filter.user != "" {
        conditions = append(conditions, "assignee = "+quote(filter.user))
    }

    // Status filter
    if filter.status != "" {
        conditions = append(conditions, "status = "+quote(filter.status))
    }

    jql := strings.Join(conditions, " AND ")
    endpoint := merge_jql_to_endpoint("/rest/api/3/search/jql", jql)

    endpoint += "&fields=summary,description,status,assignee,reporter,parent,fixVersions,priority,issuetype&maxResults=200"

    var issues []JiraIssue
    var token string

    for {
        pageURL := endpoint
        if token != "" {
            pageURL += "&nextPageToken=" + url.QueryEscape(token)
        }

        data, err := jira_get(pageURL)
        if err != nil {
            return nil, err
        }

        var response struct {
            Issues        []JiraIssue `json:"issues"`
            NextPageToken string      `json:"nextPageToken"`
        }

        if err := json.Unmarshal(data, &response); err != nil {
            return nil, err
        }

        issues = append(issues, response.Issues...)
        token = response.NextPageToken

        if token == "" {
            break
        }
    }

    return issues, nil
}
