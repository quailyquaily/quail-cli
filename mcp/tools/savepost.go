package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/util"
)

func handleSavePostTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var (
			ok            bool
			title         string
			slug          string
			summary       string
			channelSlug   string
			datetimeStr   string
			content       string
			tags          string
			coverImageURL string
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

		if _, ok = request.Params.Arguments["cover_image_url"]; ok {
			coverImageURL = request.Params.Arguments["cover_image_url"].(string)
		}
		if _, ok = request.Params.Arguments["tags"]; ok {
			tags = request.Params.Arguments["tags"].(string)
		}

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
			"cover_image_url":    coverImageURL,
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

func GetSavePostTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_save_post",
		mcp.WithDescription(`Save a post to quaily.com according to the given payload.
	The payload could be found at the current context.
	It could be a markdown file which contains the post content and the frontmatter.
	The tool will try to parse the markdown file's frontmatter and extract the title, summary, slug, datetime, tags, and cover_image_url from it.
	If no frontmatter is found, the tool will insert the frontmatter to the file.
	The content should EXCLUDE the frontmatter.
	The slug could be empty or omitted, in which case the tool will generate a slug for the post. The slug should only contain lowercase letters, numbers, and hyphens. Read it from the frontmatter.
	The datetime could be empty or omitted, in which case the tool will use an empty string. Read it from the frontmatter.
	The tags could be empty or omitted, in which case the tool will use an empty string. Read it from the frontmatter.
	The cover_image_url could be empty or omitted, in which case the tool will use an empty string. Read it from the frontmatter.
	If the post is saved successfully, the tool will return the metadata of the post,
	The tool must update the slug, datetime to the file, more specifically, the frontmatter in the markdown file.
	`),
		mcp.WithString("title", mcp.Description("Title of the post"), mcp.Required()),
		mcp.WithString("channel", mcp.Description("Channel slug to publish the post to"), mcp.Required()),
		mcp.WithString("content", mcp.Description("Content of the post."), mcp.Required()),
		mcp.WithString("slug", mcp.Description("Slug of the post"), mcp.DefaultString("")),
		mcp.WithString("summary", mcp.Description("Summary of the post"), mcp.DefaultString("")),
		mcp.WithString("datetime", mcp.Description("Datetime of the post"), mcp.DefaultString("")),
		mcp.WithString("tags", mcp.Description("Tags of the post"), mcp.DefaultString("")),
		mcp.WithString("cover_image_url", mcp.Description("Cover image url of the post"), mcp.DefaultString("")),
	)

	return tool, handleSavePostTool(cl), nil
}
