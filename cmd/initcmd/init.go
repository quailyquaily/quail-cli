package initcmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/quailyquaily/quail-cli/util"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var apiKey string

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a sample config file",
		Run: func(cmd *cobra.Command, args []string) {
			configFile, exists, err := util.ConfigFileExists()
			if err != nil {
				slog.Error("failed to check config file", "error", err)
				return
			}
			if exists {
				fmt.Printf("Config file already exists: %s\n", configFile)
				return
			}

			key := strings.TrimSpace(apiKey)
			if key == "" {
				key, err = promptOptionalAPIKey()
				if err != nil {
					slog.Error("failed to read api key", "error", err)
					return
				}
			}

			configFile, err = util.WriteSampleConfig(key)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					fmt.Printf("Config file already exists: %s\n", configFile)
					return
				}
				slog.Error("failed to create config file", "error", err)
				return
			}
			fmt.Printf("Config file created: %s\n", configFile)
		},
	}
	cmd.Flags().StringVar(&apiKey, "api-key", "", "Set an API key in the generated config")
	cmd.Flags().Lookup("api-key").NoOptDefVal = ""
	return cmd
}

func promptOptionalAPIKey() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", nil
	}

	fmt.Print("API Key (optional, press Enter to skip): ")
	buf, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}
