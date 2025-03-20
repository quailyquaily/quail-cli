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

func handleGenerateMetadataTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var (
			ok      bool
			title   string
			content string
		)

		if _, ok = request.Params.Arguments["title"]; !ok {
			return nil, fmt.Errorf("title is required")
		} else {
			title = request.Params.Arguments["title"].(string)
		}

		if _, ok = request.Params.Arguments["content"]; !ok {
			return nil, fmt.Errorf("content is required")
		} else {
			content = request.Params.Arguments["content"].(string)
		}

		result := ""
		ret, err := cl.GenerateMetadata(title, content)
		if err != nil {
			slog.Error("failed to generate metadata", "error", err)
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

func GetGenerateMetadataTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_generate_metadata",
		mcp.WithDescription(`Generate metadata (slug, summary, tags) for a post according to the given title and content.
	The payload could be found at the current file.
	It could be a markdown file which contains the post content and the title.
	The tool will try to parse the markdown file and extract the title and content from it.
	If the generation is successful, the tool will return the metadata of the post, including the slug, summary, and tags.
	The tool must update the slug, summary, and tags to the current file, more specifically, insert or update the frontmatter in the current markdown file.
	`),
		mcp.WithString("title", mcp.Description("Title of the post"), mcp.Required()),
		mcp.WithString("content", mcp.Description("Content of the post"), mcp.Required()),
	)

	return tool, handleGenerateMetadataTool(cl), nil
}
