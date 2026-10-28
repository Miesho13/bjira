package main

import(
	"fmt"
	"log"
	"strings"
	"flag"
	"os"
)

func user_name(user *JiraIssueUser) string {
    if user == nil {
        return "nil"
    }
    return user.DisplayName
}

func print_issues(issues []JiraIssue) {
	for  _, issue := range issues {
		fmt.Printf("%s\t%s\t[%s]\n", 
			color(color("id:" + issue.ID, Yellow), Bold), 
			color(issue.Key, Yellow), 
			color(issue.Fields.Status.Name, Green),
		)

		fmt.Printf(" - Assignee: %s\n", user_name(issue.Fields.Assignee))
		fmt.Printf(" - Reporter: %s\n", user_name(issue.Fields.Reporter))

		fmt.Printf("%s\n", color(color("# " + issue.Fields.Summary, Bold), Blue))
		text, err := adfToText(issue.Fields.Description);
		if err == nil {
			lines := strings.Split(text, "\n")
			for _, line := range lines {
				fmt.Printf("\t%s\n", line)
			}
		} else {
			fmt.Printf("random issue\n")
		}

		fmt.Printf("\n\n\n")
	}
}

func ls_command() error {
	flags := flag.NewFlagSet("ls", flag.ExitOnError)

	var filter IssueFilter

	flags.StringVar(&filter.user, "user", "", "Filter by user")
	flags.StringVar(&filter.status, "status", "", "Filter by status")
	flags.StringVar(&filter.sprint_name, "sprint", "", "Filter by sprint")
	flags.StringVar(&filter.reporter, "reporter", "", "Filter by reporter")
	flags.Parse(os.Args[2:])

	issues, err := get_all_issue_filter(filter)
	if err != nil {
		log.Panic(err)
	}
	print_issues(issues);

	return nil
}
