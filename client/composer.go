package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GenerateMetadata(title, content string) (*GenerateMetadataResponse, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("%s/auxilia/composer/metadata?includes=slug,summary,tags", c.APIBase), map[string]any{
		"title":   title,
		"content": content,
	})
	if err != nil {
		return nil, err
	}
	pr := &GenerateMetadataResponse{}
	if err := json.Unmarshal(resp, pr); err != nil {
		return nil, err
	}
	return pr, nil
}
