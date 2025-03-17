package login

import (
	"log/slog"

	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/quailyquaily/quail-cli/util"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to Quail using OAuth",
		Run: func(cmd *cobra.Command, args []string) {
			authBase := cmd.Context().Value(common.CTX_AUTH_BASE{}).(string)
			apiBase := cmd.Context().Value(common.CTX_API_BASE{}).(string)
			if _, err := util.Login(authBase, apiBase); err != nil {
				slog.Error("failed to login", "error", err)
				return
			}
		},
	}
}
