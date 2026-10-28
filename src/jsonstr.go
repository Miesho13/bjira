package main

import (
	"encoding/json"
)

// /rest/api/3/myself
type UserJson struct {
	Self         string `json:"self"`
	AccountID    string `json:"accountId"`
	AccountType  string `json:"accountType"`
	EmailAddress string `json:"emailAddress"`
	AvatarUrls   struct {
		Size48 string `json:"48x48"`
		Size24 string `json:"24x24"`
		Size16 string `json:"16x16"`
		Size32 string `json:"32x32"`
	} `json:"avatarUrls"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	TimeZone    string `json:"timeZone"`
	Locale      string `json:"locale"`
	Groups      struct {
		Size  int   `json:"size"`
		Items []any `json:"items"`
	} `json:"groups"`
	ApplicationRoles struct {
		Size  int   `json:"size"`
		Items []any `json:"items"`
	} `json:"applicationRoles"`
	Expand string `json:"expand"`
}

type SprintIssuesJson struct {
    Issues []struct {
        ID string `json:"id"`
    } `json:"issues"`
    IsLast bool `json:"isLast"`
}

type JiraIssueUser struct {
    AccountID   string `json:"accountId"`
    DisplayName string `json:"displayName"`
    Email       string `json:"emailAddress"`
}

type JiraNamedField struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type JiraIssueRef struct {
    ID     string `json:"id"`
    Key    string `json:"key"`
    Fields struct {
        Summary string `json:"summary"`
    } `json:"fields"`
}

type JiraIssue struct {
    ID     string `json:"id"`
    Key    string `json:"key"`

    Fields struct {
        Summary     string          `json:"summary"`
        Description json.RawMessage `json:"description"`

        Status     JiraNamedField `json:"status"`
        Priority   *JiraNamedField `json:"priority"`
        IssueType  JiraNamedField `json:"issuetype"`
        Resolution *JiraNamedField `json:"resolution"`

        Assignee *JiraIssueUser `json:"assignee"`
        Reporter *JiraIssueUser `json:"reporter"`

        Parent   *JiraIssueRef  `json:"parent"`
        Subtasks []JiraIssueRef `json:"subtasks"`

        FixVersions []JiraNamedField `json:"fixVersions"`
        Components  []JiraNamedField `json:"components"`
        Labels      []string         `json:"labels"`

        Created string `json:"created"`
        Updated string `json:"updated"`
        DueDate string `json:"duedate"`

        Project struct {
            Key  string `json:"key"`
            Name string `json:"name"`
        } `json:"project"`
    } `json:"fields"`
}
