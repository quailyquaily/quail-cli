package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Client) GetPost(listIDOrSlug string, postIDOrSlug string) (*PostResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/lists/%s/posts/%s", c.APIBase, listIDOrSlug, postIDOrSlug), nil)
	if err != nil {
		return nil, err
	}
	pr := &PostResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (c *Client) GetPostContent(listIDOrSlug string, postIDOrSlug string) (*PostContentResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/lists/%s/posts/%s/content", c.APIBase, listIDOrSlug, postIDOrSlug), nil)
	if err != nil {
		return nil, err
	}
	pr := &PostContentResponse{}
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

// GetListPosts retrieves posts from a specific list
func (c *Client) GetListPosts(listID uint64, offset, limit int) (*SearchResponse, error) {
	if listID == 0 {
		return nil, errors.New("list ID is required")
	}

	url := fmt.Sprintf("%s/lists/%d/posts?offset=%d&limit=%d", c.APIBase, listID, offset, limit)

	resp, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	sr := &SearchResponse{}
	if err := json.Unmarshal(resp, sr); err != nil {
		return nil, err
	}

	return sr, nil
}
