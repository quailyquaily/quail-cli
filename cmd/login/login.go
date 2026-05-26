package login

import (
	"log/slog"

	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/quailyquaily/quail-cli/util"
	"github.com/spf13/cobra"
)

var apiKey string

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Quail",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("api-key") {
				if err := util.LoginAPIKey(apiKey); err != nil {
					slog.Error("failed to save api key", "error", err)
					return
				}
				return
			}

			authBase := cmd.Context().Value(common.CTX_AUTH_BASE{}).(string)
			apiBase := cmd.Context().Value(common.CTX_API_BASE{}).(string)
			if _, err := util.Login(authBase, apiBase); err != nil {
				slog.Error("failed to login", "error", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "Save a Quaily API key instead of using OAuth")
	cmd.Flags().Lookup("api-key").NoOptDefVal = ""

	return cmd
}
