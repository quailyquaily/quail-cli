package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quail-ink/quail-cli/client"
)

func handleURLTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelSlug, ok := request.Params.Arguments["channel_slug"].(string)
		if !ok {
			channelSlug = ""
		}
		channelID, ok := request.Params.Arguments["channel_id"].(float64)
		if !ok {
			channelID = 0
		}
		postSlug, ok := request.Params.Arguments["post_slug"].(string)
		if !ok {
			postSlug = ""
		}
		postID, ok := request.Params.Arguments["post_id"].(float64)
		if !ok {
			postID = 0
		}
		if channelSlug == "" && channelID == 0 {
			return nil, fmt.Errorf("no channel slug or channel id provided")
		}
		mode := "channel"
		url := ""
		if postSlug != "" || postID != 0 {
			mode = "post"
		}
		if channelSlug == "" {
			resp, err := cl.GetList(uint64(channelID))
			if err != nil {
				return nil, err
			}
			channelSlug = resp.Data.Slug
		}
		if mode == "post" {
			if postSlug == "" {
				resp, err := cl.GetPost(channelSlug, fmt.Sprintf("%d", uint64(postID)))
				if err != nil {
					return nil, err
				}
				postSlug = resp.Data.Slug
			}
			url = fmt.Sprintf("https://quaily.com/%s/p/%s", channelSlug, postSlug)
		} else {
			url = fmt.Sprintf("https://quaily.com/%s", channelSlug)
		}

		result := make(map[string]string)
		result["url"] = url
		buf, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(buf),
				},
			},
		}, nil
	}
}

func GetURLTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("get_quaily_url",
		mcp.WithDescription(`Return the url of given channel or post.
		The format of an URL should be like:
		- for channel: https://quaily.com/{channel_slug}
		- for post: https://quaily.com/{channel_slug}/p/{post_slug}
		so, if the channel slug is given, the tool will return the url of the channel.
		if the post slug is given, the tool will return the url of the post.
		if the channel id is given, the tool will try to get the channel slug from the channel id, and then return the url of the channel.
		if the post id is given, the tool will try to get the post slug from the post id, and then return the url of the post.
		`),
		mcp.WithString("channel_slug",
			mcp.Description("The slug of the channel"),
			mcp.Required(),
		),
		mcp.WithNumber("channel_id",
			mcp.Description("The id of the channel"),
			mcp.DefaultNumber(0),
		),
		mcp.WithString("post_slug",
			mcp.Description("The slug of the post"),
			mcp.DefaultString(""),
		),
		mcp.WithNumber("post_id",
			mcp.Description("The id of the post"),
			mcp.DefaultNumber(0),
		),
	)
	return tool, handleURLTool(cl), nil
}
