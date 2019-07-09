package jira

import "github.com/google/go-querystring/query"
import "fmt"

// FilterService handles fields for the JIRA instance / API.
//
// JIRA API docs: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-group-Filter
type FilterService struct {
	client *Client
}

// Filter represents a Filter in Jira
type Filter struct {
	Self             string        `json:"self"`
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Owner            User          `json:"owner"`
	Jql              string        `json:"jql"`
	ViewURL          string        `json:"viewUrl"`
	SearchURL        string        `json:"searchUrl"`
	Favourite        bool          `json:"favourite"`
	FavouritedCount  int           `json:"favouritedCount"`
	SharePermissions []interface{} `json:"sharePermissions"`
	Subscriptions    struct {
		Size       int           `json:"size"`
		Items      []interface{} `json:"items"`
		MaxResults int           `json:"max-results"`
		StartIndex int           `json:"start-index"`
		EndIndex   int           `json:"end-index"`
	} `json:"subscriptions"`
}

// GetList retrieves all filters from Jira
func (fs *FilterService) GetList() ([]*Filter, *Response, error) {

	options := &GetQueryOptions{}
	apiEndpoint := "rest/api/2/filter"
	req, err := fs.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	filters := []*Filter{}
	resp, err := fs.client.Do(req, &filters)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return filters, resp, err
}

// GetFavouriteList retrieves the user's favourited filters from Jira
func (fs *FilterService) GetFavouriteList() ([]*Filter, *Response, error) {
	apiEndpoint := "rest/api/2/filter/favourite"
	req, err := fs.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	filters := []*Filter{}
	resp, err := fs.client.Do(req, &filters)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}
	return filters, resp, err
}

// Get retrieves a single Filter from Jira
func (fs *FilterService) Get(filterID int) (*Filter, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/filter/%d", filterID)
	req, err := fs.client.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	filter := new(Filter)
	resp, err := fs.client.Do(req, filter)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return filter, resp, err
}
