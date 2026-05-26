package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetCommentsByPost(postID uint64, offset, limit int) (*CommentsResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/comments?post_id=%d&offset=%d&limit=%d", c.APIBase, postID, offset, limit), nil)
	if err != nil {
		return nil, err
	}
	cr := &CommentsResponse{}
	if err := json.Unmarshal(resp, cr); err != nil {
		return nil, err
	}
	return cr, nil
}

func (c *Client) CreateComment(postID uint64, content string) (*CommentResponse, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("%s/comments", c.APIBase), map[string]any{
		"post_id": postID,
		"content": content,
	})
	if err != nil {
		return nil, err
	}
	cr := &CommentResponse{}
	if err := json.Unmarshal(resp, cr); err != nil {
		return nil, err
	}
	return cr, nil
}
