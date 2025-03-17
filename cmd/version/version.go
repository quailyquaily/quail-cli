package version

import (
	"fmt"

	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			version := ctx.Value(common.CTX_VERSION{}).(string)
			fmt.Printf("%s\n", version)
		},
	}

	return cmd
}
