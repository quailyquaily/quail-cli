package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
)

func handleListPostsTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		listID, ok := request.Params.Arguments["channel_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("channel id is required")
		}
		offset, ok := request.Params.Arguments["offset"].(float64)
		if !ok {
			offset = 0
		}
		limit, ok := request.Params.Arguments["limit"].(float64)
		if !ok {
			limit = 20
		}

		resp, err := cl.GetListPosts(uint64(listID), int(offset), int(limit))
		if err != nil {
			slog.Error("failed to get list posts", "error", err, "list_id", listID, "offset", offset, "limit", limit)
			return nil, err
		}
		res := make([]string, 0)
		for _, item := range resp.Data.Items {
			jsonItem := make(map[string]any)
			jsonItem["id"] = item.ID
			jsonItem["slug"] = item.Slug
			jsonItem["title"] = item.Title
			jsonItem["summary"] = item.Summary
			buf, err := json.Marshal(jsonItem)
			if err != nil {
				continue
			}
			res = append(res, string(buf))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: strings.Join(res, "\n\n"),
				},
			},
		}, nil
	}
}

func GetListPostsTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("get_my_channel_posts",
		mcp.WithDescription("Return the posts of a given channel. The channel is specified by the channel id."),
		mcp.WithNumber("channel_id",
			mcp.Description("The id of the channel"),
			mcp.Required(),
		),
		mcp.WithNumber("offset",
			mcp.Description("The offset of the posts to return"),
			mcp.DefaultNumber(0),
		),
		mcp.WithNumber("limit",
			mcp.Description("The limit of the posts to return"),
			mcp.DefaultNumber(20),
		),
	)
	return tool, handleListPostsTool(cl), nil
}
