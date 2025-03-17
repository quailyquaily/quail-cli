package me

import (
	"log/slog"

	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Get current user information",
		Run: func(cmd *cobra.Command, args []string) {
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
			result, err := cl.GetMe()
			if err != nil {
				slog.Error("failed to get user information", "error", err)
				return
			}
			if format == common.FORMAT_JSON {
				client.PrettyPrintJSON(result)
			} else {
				client.PrettyPrintUser(result)
			}
		},
	}
}
