package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
)

func handlePostContentTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelSlug, ok := request.Params.Arguments["channel_slug"].(string)
		if !ok {
			channelSlug = ""
		}
		postSlug, ok := request.Params.Arguments["post_slug"].(string)
		if !ok {
			postSlug = ""
		}

		if channelSlug == "" || postSlug == "" {
			// try to get the channel slug and post slug from the URL
			url, ok := request.Params.Arguments["url"].(string)
			if !ok {
				url = ""
			}
			re := regexp.MustCompile(`https://quaily\.com/([^/]+)/([^/]+)`)
			matches := re.FindStringSubmatch(url)
			if len(matches) == 3 {
				channelSlug = matches[1]
				postSlug = matches[2]
			}
		}

		var msg string

		if channelSlug == "" || postSlug == "" {
			msg = "no channel slug or post slug or a valid URL provided"
		} else {
			resp, err := cl.GetPostContent(channelSlug, postSlug)
			if err != nil {
				msg = fmt.Sprintf("failed to get post. error=%v", err)
			} else {
				buf, err := json.Marshal(resp.Data)
				if err != nil {
					msg = fmt.Sprintf("failed to marshal post. error=%v", err)
				} else {
					msg = string(buf)
				}
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: msg,
				},
			},
		}, nil
	}
}

func GetPostContentTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("get_quaily_post_content",
		mcp.WithDescription(`Return the post of a given quaily channel and post slug. The tool will return both content and paid content if the user is logged in and paid for the post.
		The channel is specified by the channel slug.
		If there is an URL, the tool will accept the URL as well. The URL should be in the format of "https://quaily.com/{channel_slug}/{post_slug}/content".
		The tool returns the post data in JSON format. The content of the post is included in the field "content" and "paid_content" if the post is paid.
		`),
		mcp.WithString("channel_slug",
			mcp.Description("The slug of the channel"),
		),
		mcp.WithString("post_slug",
			mcp.Description("The slug of the post"),
		),
		mcp.WithString("url",
			mcp.Description("The URL of the post"),
		),
	)
	return tool, handlePostContentTool(cl), nil
}
