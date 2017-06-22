package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
)

// UserService handles users for the JIRA instance / API.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user
type UserService struct {
	client *Client
}

// User represents a JIRA user.
type User struct {
	Self            string     `json:"self,omitempty" structs:"self,omitempty"`
	Name            string     `json:"name,omitempty" structs:"name,omitempty"`
	Password        string     `json:"-"`
	Key             string     `json:"key,omitempty" structs:"key,omitempty"`
	EmailAddress    string     `json:"emailAddress,omitempty" structs:"emailAddress,omitempty"`
	AvatarUrls      AvatarUrls `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	DisplayName     string     `json:"displayName,omitempty" structs:"displayName,omitempty"`
	Active          bool       `json:"active,omitempty" structs:"active,omitempty"`
	TimeZone        string     `json:"timeZone,omitempty" structs:"timeZone,omitempty"`
	ApplicationKeys []string   `json:"applicationKeys,omitempty" structs:"applicationKeys,omitempty"`
}

type UserPermissionSearch struct {
	Username    string `json:"username,omitempty"`
	Permissions string `json:"permissions,omitempty"`
	IssueKey    string `json:"issueKey,omitempty"`
	ProjectKey  string `json:"projectKey,omitempty"`
	StartAt     int    `json:"startAt,omitempty"`
	MaxResults  int    `json:"maxResults,omitempty"`
}

// Get gets user info from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-getUser
func (s *UserService) Get(username string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("/rest/api/2/user?username=%s", username)
	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Myself gets the current user from JIRA
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/myself-getUser
func (s *UserService) Myself() (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "/rest/api/2/myself", nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Create creates an user in JIRA.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-createUser
func (s *UserService) Create(user *User) (*User, *Response, error) {
	apiEndpoint := "/rest/api/2/user"
	req, err := s.client.NewRequest("POST", apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	responseUser := new(User)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not read the returned data")
	}
	err = json.Unmarshal(data, responseUser)
	if err != nil {
		return nil, resp, fmt.Errorf("Could not unmarshall the data into struct")
	}
	return responseUser, resp, nil
}

// Search for users based on permissions in JIRA.
//
// JIRA API docs: https://docs.atlassian.com/jira/REST/cloud/#api/2/user-findUsersWithAllPermissions
func (s *UserService) PermissionSearch(search UserPermissionSearch) (*[]User, *Response, error) {
	apiEndpoint := "/rest/api/2/user/permission/search"
	v := url.Values{}
	if search.IssueKey != "" {
		v.Set("issueKey", search.IssueKey)
	}
	if search.MaxResults != 0 {
		v.Set("maxResults", strconv.Itoa(search.MaxResults))
	} else {
		v.Set("maxResults", "1000")
	}
	if search.Permissions != "" {
		v.Set("permissions", search.Permissions)
	}
	if search.ProjectKey != "" {
		v.Set("projectKey", search.ProjectKey)
	}
	if search.StartAt != 0 {
		v.Set("startAt", strconv.Itoa(search.StartAt))
	}
	if search.Username != "" {
		v.Set("username", search.Username)
	}
	query := v.Encode()
	if query != "" {
		apiEndpoint = apiEndpoint + "?" + query
	}

	req, err := s.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	users := make([]User, 0)
	resp, err := s.client.Do(req, &users)
	if err != nil {
		return nil, resp, err
	}
	return &users, resp, nil
}
