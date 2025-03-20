package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
)

func handleUnpublishPostTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var (
			ok          bool
			slug        string
			channelSlug string
		)

		if _, ok = request.Params.Arguments["slug"]; !ok {
			slug = ""
		} else {
			slug = request.Params.Arguments["slug"].(string)
		}
		if _, ok = request.Params.Arguments["channel"]; !ok {
			return nil, fmt.Errorf("channel is required")
		} else {
			channelSlug = request.Params.Arguments["channel"].(string)
		}

		result := ""

		ret, err := cl.UnpublishPost(channelSlug, slug)
		if err != nil {
			slog.Error("failed to unpublish post", "error", err)
			result = err.Error()
		} else {
			buf, err := json.Marshal(ret.Data)
			if err != nil {
				slog.Error("failed to marshal post", "error", err)
				result = err.Error()
			} else {
				result = string(buf)
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: result,
				},
			},
		}, nil
	}
}

func GetUnpublishPostTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_unpublish_post",
		mcp.WithDescription(`Unpublish a post from quaily.com according to the given channel slug and post slug.
	If the post is unpublished successfully, the tool will return the metadata of the post.
	`),
		mcp.WithString("channel", mcp.Description("Channel slug to publish the post to"), mcp.Required()),
		mcp.WithString("slug", mcp.Description("Slug of the post"), mcp.Required()),
	)

	return tool, handleUnpublishPostTool(cl), nil
}
