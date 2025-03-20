package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/quailyquaily/quail-cli/util"
)

func handleLoginTool(pctx context.Context, cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		authBase := pctx.Value(common.CTX_AUTH_BASE{}).(string)
		apiBase := pctx.Value(common.CTX_API_BASE{}).(string)

		var msg string
		authCodeURL, err := util.Login(authBase, apiBase)
		if err != nil {
			msg = fmt.Sprintf("failed to login. error=%v, auth_url=%s. If the browser does not open automatically, please open auth_url manually.", err, authCodeURL)
		} else {
			msg = fmt.Sprintf("login successfully. auth_url=%s. If the browser does not open automatically, please open auth_url manually.", authCodeURL)
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

func LoginTool(ctx context.Context, cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_login",
		mcp.WithDescription("Login to quaily.com. If the browser does not open automatically, please open auth_url manually."),
	)

	return tool, handleLoginTool(ctx, cl), nil
}
