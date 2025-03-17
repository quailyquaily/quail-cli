/*
Copyright © 2024 lyric
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/quailyquaily/quail-cli/cmd/login"
	"github.com/quailyquaily/quail-cli/cmd/mcp"
	"github.com/quailyquaily/quail-cli/cmd/me"
	"github.com/quailyquaily/quail-cli/cmd/post"
	"github.com/quailyquaily/quail-cli/cmd/version"
	"github.com/quailyquaily/quail-cli/oauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	authBase    string
	apiBase     string
	accessToken string
	format      string
	cl          *client.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "quail-cli",
	Short: "A CLI tool for interacting with Quaily's API",
	Long:  `quail-cli is a command-line interface for sending requests to Quaily's API at https://api.quail.ink`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		ctx = context.WithValue(ctx, common.CTX_CLIENT{}, cl)
		ctx = context.WithValue(ctx, common.CTX_API_BASE{}, apiBase)
		ctx = context.WithValue(ctx, common.CTX_AUTH_BASE{}, authBase)
		ctx = context.WithValue(ctx, common.CTX_FORMAT{}, format)

		cmd.SetContext(ctx)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func ExecuteContext(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/quail-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiBase, "api-base", "https://api.quail.ink", "Quail API base URL")
	rootCmd.PersistentFlags().StringVar(&authBase, "auth-base", "https://quaily.com", "Quail Auth base URL")
	rootCmd.PersistentFlags().StringVar(&format, "format", "human", "the output format (human: human readable, json: JSON)")

	rootCmd.AddCommand(login.NewCmd())
	rootCmd.AddCommand(me.NewCmd())
	rootCmd.AddCommand(post.NewCmd())
	rootCmd.AddCommand(mcp.NewCmd())
	rootCmd.AddCommand(version.NewCmd())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {

		fullpath := common.GetConfigFilePath()

		viper.AddConfigPath(fullpath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		cfgFile = filepath.Join(fullpath, "config.yaml")

		viper.SetConfigFile(cfgFile)

		if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
			// if the config file does not exist, ask the user to login
			fmt.Println("Config file does not exist. Please login.")
			common.Login(authBase, apiBase)
			return
		}
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("failed to read config", err, "config", viper.ConfigFileUsed())
		return
	}

	accessToken = viper.GetString("app.access_token")
	expiry := viper.GetTime("app.expiry")

	if time.Now().After(expiry) {
		// if the access token has expired, try to get a new one using the refresh token
		fmt.Println("Access token has expired. Try to get new one.")
		refreshToken := viper.GetString("app.refresh_token")
		token, err := oauth.RefreshToken(apiBase, refreshToken)
		if err != nil {
			slog.Error("failed to refresh token", "error", err)
			return
		}
		// update the config file with the new access token
		viper.Set("app.access_token", token.AccessToken)
		viper.Set("app.expiry", token.Expiry)
		viper.Set("app.token_type", token.TokenType)
		viper.Set("app.refresh_token", token.RefreshToken)

		viper.WriteConfig()

		accessToken = token.AccessToken
	}

	cl = client.New(accessToken, apiBase)
}
