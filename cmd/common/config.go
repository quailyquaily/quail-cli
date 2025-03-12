package common

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/quail-ink/quail-cli/client"
	"github.com/quail-ink/quail-cli/oauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	CTX_CLIENT    struct{}
	CTX_API_BASE  struct{}
	CTX_AUTH_BASE struct{}
	CTX_FORMAT    struct{}
	CTX_VERSION   struct{}
)

const (
	FORMAT_JSON  = "json"
	FORMAT_HUMAN = "human"
)

func GetConfigFilePath() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	fullpath := filepath.Join(home, ".config", "quail-cli")
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		os.MkdirAll(fullpath, 0755)
	}

	return fullpath
}

func Login(authBase, apiBase string) (err error) {
	token, err := oauth.Login(authBase, apiBase)
	if err != nil {
		slog.Error("failed to login", "error", err)
		return
	}

	viper.Set("app.access_token", token.AccessToken)
	viper.Set("app.refresh_token", token.RefreshToken)
	viper.Set("app.token_type", token.TokenType)
	viper.Set("app.expiry", token.Expiry)

	cl := client.New(token.AccessToken, apiBase)
	result, err := cl.GetMe()
	if err != nil {
		slog.Error("failed to get me", "error", err)
		return
	}
	viper.Set("app.user.id", result.Data.ID)
	viper.Set("app.user.name", result.Data.Name)
	viper.Set("app.user.bio", result.Data.Bio)

	fullpath := GetConfigFilePath()

	// if the config file doesn't exist, create it first
	err = viper.WriteConfigAs(path.Join(fullpath, "config.yaml"))
	if err != nil {
		slog.Error("failed to save config", "error", err, "dir", fullpath)
		return
	}

	fmt.Printf("Login successful. Access token saved to %s\n", fullpath)
	return nil
}
