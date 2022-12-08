package chglog

import (
	"fmt"

	agjira "github.com/andygrunwald/go-jira"
)

// JiraClient is an HTTP client for Jira
type JiraClient interface {
	GetJiraIssue(id string) (*agjira.Issue, error)
}

type jiraClient struct {
	username    string
	token       string
	bearerToken string
	url         string
}

// NewJiraClient returns an instance of JiraClient
func NewJiraClient(config *Config) JiraClient {
	return jiraClient{
		username:    config.Options.JiraUsername,
		token:       config.Options.JiraToken,
		bearerToken: config.Options.JiraBearerToken,
		url:         config.Options.JiraURL,
	}
}

func (jira jiraClient) client() (*agjira.Client, error) {
	if jira.bearerToken != "" {
		tp := agjira.BearerAuthTransport{
			Token: jira.bearerToken,
		}
		return agjira.NewClient(tp.Client(), jira.url)
	}

	tp := agjira.BasicAuthTransport{
		Username: jira.username,
		Password: jira.token,
	}
	return agjira.NewClient(tp.Client(), jira.url)
}

func (jira jiraClient) GetJiraIssue(id string) (*agjira.Issue, error) {
	client, err := jira.client()
	if err != nil {
		return nil, fmt.Errorf("cannot instantiate jira client, got:%v", err)
	}
	issue, _, err := client.Issue.Get(id, nil)
	return issue, err
}
