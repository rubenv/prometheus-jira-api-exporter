package jiraapiexporter

import (
	"context"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/prometheus/client_golang/prometheus"
)

func Update(ctx context.Context, host, user, pass string) error {
	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	client, err := jira.NewClient(tp.Client(), host)
	if err != nil {
		return err
	}

	projects, _, err := client.Project.GetList()
	if err != nil {
		return err
	}

	projectNames := make(map[string]string)
	for _, project := range *projects {
		projectNames[project.Key] = project.Name
	}

	// Fetch all issues
	issueCounts := NewIssueCounts()

	err = client.Issue.SearchPages("", &jira.SearchOptions{
		MaxResults: 100,
		Fields:     []string{"fixVersions", "issuetype"},
	}, func(issue jira.Issue) error {
		parts := strings.Split(issue.Key, "-")
		project := parts[0]
		issueType := issue.Fields.Type.Name
		if len(issue.Fields.FixVersions) == 0 {
			issueCounts.Count(project, "Backlog", issueType)
		} else {
			for _, fv := range issue.Fields.FixVersions {
				issueCounts.Count(project, fv.Name, issueType)
			}
		}
		return nil
	})

	for project, releases := range issueCounts {
		for release, issueTypes := range releases {
			for issueType, count := range issueTypes {
				issuesMetric.With(prometheus.Labels{
					"project":    projectNames[project],
					"projectkey": project,
					"release":    release,
					"type":       issueType,
				}).Set(count)
			}
		}
	}

	return nil
}
