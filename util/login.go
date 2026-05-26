package util

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/oauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
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

func LoginAPIKey(apiKey string) error {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		fmt.Print("API Key: ")
		buf, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}
		apiKey = strings.TrimSpace(string(buf))
	}
	if apiKey == "" {
		return fmt.Errorf("api key is required")
	}
	if !strings.HasPrefix(apiKey, "QK-") {
		return fmt.Errorf("invalid api key")
	}

	viper.Set("app.api_key", apiKey)

	if err := writeConfig(); err != nil {
		slog.Error("failed to save config", "error", err)
		return err
	}

	fmt.Println("API key saved.")
	return nil
}

func Login(authBase, apiBase string) (authCodeURL string, err error) {
	token, authCodeURL, err := oauth.Login(authBase, apiBase)
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

	if err = writeConfig(); err != nil {
		slog.Error("failed to save config", "error", err)
		return
	}

	fmt.Println("Login successful.")
	return
}

func writeConfig() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(GetConfigFilePath(), "config.yaml")
		viper.SetConfigFile(configFile)
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
		return err
	}
	if _, err := os.Stat(configFile); err == nil {
		return viper.WriteConfig()
	}
	return viper.WriteConfigAs(configFile)
}
