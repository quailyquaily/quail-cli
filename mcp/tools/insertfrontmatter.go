package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
)

func handleInsertFrontmatterTool() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var (
			ok    bool
			title string
		)

		if _, ok = request.Params.Arguments["title"]; !ok {
			return nil, fmt.Errorf("title is required")
		} else {
			title = request.Params.Arguments["title"].(string)
		}

		datetime := time.Now().Format("2006-01-02 15:04")

		result := fmt.Sprintf(`---
title: "%s"
slug: ""
datetime: "%s"
summary: ""
tags: []
theme: light
cover_image_url: ""
---

`, title, datetime)

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

func GetInsertFrontmatterTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_insert_frontmatter",
		mcp.WithDescription(`Insert the given frontmatter to the current file.
		The frontmatter should be a valid yaml frontmatter.
		The tool need to check the current file is a markdown file. If it's a markdown file and the frontmatter is not found, the tool will insert the frontmatter to current file.
		The tool will try to merge the new frontmatter with the existing frontmatter.
		The title could be the first heading in the file, or the file name.
	`),
		mcp.WithString("title", mcp.Description("Title of the post"), mcp.Required()),
	)

	return tool, handleInsertFrontmatterTool(), nil
}
