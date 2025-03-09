package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetPost(listIDOrSlug string, slug string) (*PostResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/lists/%s/posts/%s", c.APIBase, listIDOrSlug, slug), nil)
	if err != nil {
		return nil, err
	}
	pr := &PostResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (c *Client) CreatePost(listIDOrSlug string, payload map[string]any) (*PostResponse, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("%s/lists/%s/posts", c.APIBase, listIDOrSlug), payload)
	if err != nil {
		return nil, err
	}
	pr := &PostResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (c *Client) DeletePost(listIDOrSlug string, slug string) (*PostResponse, error) {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("%s/lists/%s/posts/%s", c.APIBase, listIDOrSlug, slug), nil)
	if err != nil {
		return nil, err
	}
	pr := &PostResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, err
}

func (c *Client) ModPost(listIDOrSlug, slug, op string) (*PostResponse, error) {
	resp, err := c.sendRequest("PUT", fmt.Sprintf("%s/lists/%s/posts/%s/%s", c.APIBase, listIDOrSlug, slug, op), nil)
	if err != nil {
		return nil, err
	}
	pr := &PostResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (c *Client) Search(query string) (*SearchResponse, error) {
	payload := make(map[string]any)
	payload["q"] = query
	resp, err := c.sendRequest("POST", fmt.Sprintf("%s/posts/search", c.APIBase), payload)
	if err != nil {
		return nil, err
	}
	pr := &SearchResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}
