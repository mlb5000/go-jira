package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// WebhookService handles webhooks for the JIRA instance / API.
//
// JIRA API docs: https://developer.atlassian.com/jiradev/jira-apis/webhooks
type WebhookService struct {
	client *Client
}

// Webhook represents a JIRA webhook.
type Webhook struct {
	Name                string   `json:"name,omitempty" structs:"name,omitempty"`
	Url                 string   `json:"url,omitempty" structs:"url,omitempty"`
	Events              []string `json:"events,omitempty" structs:"events,omitempty"`
	JqlFilter           string   `json:"jqlFilter,omitempty" structs:"jqlFilter,omitempty"`
	ExcludeIssueDetails bool     `json:"excludeIssueDetails,omitempty" structs:"excludeIssueDetails,omitempty"`
}

// Create creates a webhook in JIRA.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-createUser
func (s *WebhookService) Create(webhook *Webhook) (*Webhook, *Response, error) {
	apiEndpoint := "/rest/webhooks/1.0/webhook"
	req, err := s.client.NewRequest("POST", apiEndpoint, webhook)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseWebhook := new(Webhook)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	err = json.Unmarshal(data, responseWebhook)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return responseWebhook, resp, nil
}

// Gets all webhooks on the JIRA instance.
//
// JIRA API docs: https://developer.atlassian.com/jiradev/jira-apis/webhooks#Webhooks-Registeringawebhook
func (s *WebhookService) GetAll() (*[]Webhook, *Response, error) {
	apiEndpoint := "/rest/webhooks/1.0/webhook"
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseWebhook := make([]Webhook, 0)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	err = json.Unmarshal(data, &responseWebhook)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return &responseWebhook, resp, nil
}
