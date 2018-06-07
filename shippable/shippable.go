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

// Subscription is a fundamental object in shippable
// it contains the organization name, a unique ID,
// and the total minion count
type Subscription struct {
	OrgName     string `json:"orgName"`
	ID          string `json:"id"`
	MinionCount int    `json:"minionCount"`
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
func (c *Client) ListSubscriptions(q string) ([]Subscription, error) {

	rel := &url.URL{Path: "/subscriptions"}
	u := c.BaseURL.ResolveReference(rel)
	u.RawQuery = q
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
