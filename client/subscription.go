package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetSubscriptions() (*SubscriptionsResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/subscriptions/", c.APIBase), nil)
	if err != nil {
		return nil, err
	}
	sr := &SubscriptionsResponse{}
	if err := json.Unmarshal(resp, sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (c *Client) GetSubscribedPosts(offset, limit int) (*SearchResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/posts/subscribed?offset=%d&limit=%d", c.APIBase, offset, limit), nil)
	if err != nil {
		return nil, err
	}
	sr := &SearchResponse{}
	if err := json.Unmarshal(resp, sr); err != nil {
		return nil, err
	}
	return sr, nil
}
