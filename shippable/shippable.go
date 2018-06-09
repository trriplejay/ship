package shippable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DefaultURL is used if none is set
const DefaultURL = "https://api.shippable.com"

// Client provides methods for interacting with shippable's api
type Client struct {
	BaseURL    *url.URL
	Token      string
	httpClient *http.Client
}

// Subscription is an object that represents an organization
// It contains the organization name, a unique ID,
// and the total minion count
type Subscription struct {
	ID          string `json:"id"`
	OrgName     string `json:"orgName"`
	MinionCount int    `json:"minionCount"`
}

// Project is an object that represents a repository on shippable
type Project struct {
	ID                  string `json:"id"`
	FullName            string `json:"fullName"`
	Name                string `json:"name"`
	SubscriptionID      string `json:"subscriptionId"`
	BuilderAccountID    string `json:"builderAccountId"`
	RepositoryURL       string `json:"repositoryUrl"`
	SourceDefaultBranch string `json:"sourceDefaultBranch"`
	IsPrivateRepository bool   `json:"isPrivateRepository"`
}

// NewClient is a constructor for the shippable package
// it should be passed a host and an api token
// if no host is provided, api.shippable.com will be used by default
func NewClient(host string, token string) (*Client, error) {
	var u *url.URL
	var err error
	if host != "" {
		u, err = url.Parse(host)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = url.Parse(DefaultURL)
		if err != nil {
			return nil, err
		}
	}
	return &Client{BaseURL: u, Token: token, httpClient: &http.Client{}}, nil

}

// ListSubscriptions return an array of shippable subscription objects
// example usage:
//  t := os.Getenv("API_TOKEN")
// 	c, err := shippable.NewClient("", t)
// 	if err != nil {
// 		fmt.Println("Couldn't create client")
// 		os.Exit(1)
// 	}
// 	subs, err := c.ListSubscriptions("subscriptionOrgNames=trriplejay")
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	fmt.Printf("%+v\n", subs)
func (c *Client) ListSubscriptions(query string) ([]Subscription, error) {

	rel := &url.URL{Path: "/subscriptions"}
	u := c.BaseURL.ResolveReference(rel)
	u.RawQuery = query
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	fullAuth := fmt.Sprintf("apiToken %s", c.Token)
	req.Header.Set("Authorization", fullAuth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	s := resp.StatusCode
	if s != http.StatusOK {
		fmt.Printf("status is %s\n", resp.Status)
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}

	defer resp.Body.Close()
	var subscriptions []Subscription
	err = json.NewDecoder(resp.Body).Decode(&subscriptions)
	return subscriptions, err
}

// ListProjects returns an array of project objects
func (c *Client) ListProjects(query string) ([]Project, error) {
	rel := &url.URL{Path: "/projects"}
	u := c.BaseURL.ResolveReference(rel)
	u.RawQuery = query
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	fullAuth := fmt.Sprintf("apiToken %s", c.Token)
	req.Header.Set("Authorization", fullAuth)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	s := resp.StatusCode
	if s != http.StatusOK {
		fmt.Printf("status is %s\n", resp.Status)
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}

	defer resp.Body.Close()
	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	return projects, err
}
