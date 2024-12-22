package common

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/quail-ink/quail-cli/oauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	CTX_CONFIG_FILE struct{}
	CTX_CLIENT      struct{}
	CTX_API_BASE    struct{}
	CTX_AUTH_BASE   struct{}
	CTX_FORMAT      struct{}
)

const (
	FORMAT_JSON  = "json"
	FORMAT_HUMAN = "human"
)

func ConfigViper(cfgFile string) string {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		fullpath := filepath.Join(home, ".config", "quail-cli")
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			os.MkdirAll(fullpath, 0755)
		}

		viper.AddConfigPath(fullpath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		cfgFile = filepath.Join(fullpath, "config.yaml")
	}

	viper.AutomaticEnv()

	return cfgFile
}

func Login(authBase, apiBase string) (err error) {
	token, err := oauth.Login(authBase, apiBase)
	if err != nil {
		slog.Error("failed to login", "error", err)
		return
	}

	fullpath := ConfigViper("")

	viper.Set("app.access_token", token.AccessToken)
	viper.Set("app.refresh_token", token.RefreshToken)
	viper.Set("app.token_type", token.TokenType)
	viper.Set("app.expiry", token.Expiry)

	// if the config file doesn't exist, create it first
	err = viper.WriteConfigAs(fullpath)
	if err != nil {
		slog.Error("failed to save config", "error", err, "config", fullpath)
		return
	}

	fmt.Printf("Login successful. Access token saved to %s\n", fullpath)
	return nil
}
