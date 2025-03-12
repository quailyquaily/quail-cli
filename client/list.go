package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetUserLists(userID uint64) ([]List, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/users/%d/lists", c.APIBase, userID), nil)
	if err != nil {
		return nil, err
	}
	ltsp := &ListsResponse{}
	if err := json.Unmarshal(resp, ltsp); err != nil {
		return nil, err
	}
	return ltsp.Data, nil
}
