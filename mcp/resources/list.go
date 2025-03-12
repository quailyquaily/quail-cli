package resources

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quail-ink/quail-cli/client"
	"github.com/spf13/viper"
)

func handleListsResource(cl *client.Client) func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		userID := viper.GetInt64("app.user.id")

		lists, err := cl.GetUserLists(uint64(userID))
		if err != nil {
			slog.Error("failed to get user lists", "error", err)
			return nil, err
		}

		res := make([]mcp.ResourceContents, 0)
		for _, list := range lists {
			res = append(res, mcp.TextResourceContents{
				URI:      fmt.Sprintf("https://quaily.com/%s", list.Slug),
				MIMEType: "text/html",
				Text:     fmt.Sprintf("Title: %s, Description: %s, Tagline: %s", list.Title, list.Description, list.Tagline),
			})
		}

		return res, nil
	}
}

func GetListsResource(cl *client.Client) (mcp.Resource, mcps.ResourceHandlerFunc, error) {
	res := mcp.NewResource(
		"https://quaily.com/me/channels",
		"My Quaily Channels",
	)

	return res, handleListsResource(cl), nil
}
