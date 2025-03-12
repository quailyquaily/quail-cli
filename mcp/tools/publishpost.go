package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quail-ink/quail-cli/client"
	"github.com/quail-ink/quail-cli/util"
)

func handlePublishPostTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("request.Params.Arguments: %+v\n", request.Params.Arguments)
		var (
			ok          bool
			title       string
			slug        string
			summary     string
			channelSlug string
			datetimeStr string
			content     string
			tags        string
		)

		if _, ok = request.Params.Arguments["title"]; !ok {
			return nil, fmt.Errorf("title is required")
		} else {
			title = request.Params.Arguments["title"].(string)
		}
		if _, ok = request.Params.Arguments["slug"]; !ok {
			slug = ""
		} else {
			slug = request.Params.Arguments["slug"].(string)
		}
		if _, ok = request.Params.Arguments["summary"]; !ok {
			summary = ""
		} else {
			summary = request.Params.Arguments["summary"].(string)
		}
		if _, ok = request.Params.Arguments["channel"]; !ok {
			return nil, fmt.Errorf("channel is required")
		} else {
			channelSlug = request.Params.Arguments["channel"].(string)
		}
		if _, ok = request.Params.Arguments["datetime"]; ok {
			datetimeStr = request.Params.Arguments["datetime"].(string)
		}
		datetimeStr = util.AnyDatetimeToRFC3339(datetimeStr)

		if _, ok = request.Params.Arguments["content"]; !ok {
			return nil, fmt.Errorf("content is required")
		} else {
			content = request.Params.Arguments["content"].(string)
		}
		payload := map[string]any{
			"slug":               slug,
			"title":              title,
			"summary":            summary,
			"channel":            channelSlug,
			"content":            content,
			"cover_image_url":    "https://quaily.com/images/default-cover.png",
			"datetime":           datetimeStr,
			"first_published_at": datetimeStr,
			"tags":               tags,
		}

		result := ""
		ret, err := cl.CreatePost(channelSlug, payload)
		if err != nil {
			slog.Error("failed to create post", "error", err)
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

func GetPublishPostTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("publish_post",
		mcp.WithDescription(`Publish a post to quaily according to the given payload.
	The payload could be found at the current context.
	It could be a markdown file which contains the post content and the frontmatter.
	The tool will try to parse the markdown file and extract the title, summary, channel slug, and content from it.
	The slug could be empty or omitted, in which case the tool will generate a slug for the post. The slug should only contain lowercase letters, numbers, and hyphens.
	The datetime could be empty or omitted, in which case the tool will use an empty string.
	The tags could be empty or omitted, in which case the tool will use an empty string.
	If the post is published successfully, the tool will return the metadata of the post.
	The tool must update the slug, datetime to the file, more specifically, the frontmatter.
	`),
		mcp.WithString("title", mcp.Description("Title of the post"), mcp.Required()),
		mcp.WithString("channel", mcp.Description("Channel slug to publish the post to"), mcp.Required()),
		mcp.WithString("content", mcp.Description("Content of the post"), mcp.Required()),
		mcp.WithString("slug", mcp.Description("Slug of the post"), mcp.DefaultString("")),
		mcp.WithString("summary", mcp.Description("Summary of the post"), mcp.DefaultString("")),
		mcp.WithString("datetime", mcp.Description("Datetime of the post"), mcp.DefaultString("")),
		mcp.WithString("tags", mcp.Description("Tags of the post"), mcp.DefaultString("")),
	)

	return tool, handlePublishPostTool(cl), nil
}
