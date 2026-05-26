package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func WriteSampleConfig(apiKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey != "" && !strings.HasPrefix(apiKey, "QK-") {
		return "", fmt.Errorf("invalid api key")
	}

	configFile := ResolveConfigFile()

	if _, err := os.Stat(configFile); err == nil {
		return configFile, os.ErrExist
	} else if !os.IsNotExist(err) {
		return configFile, err
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
		return configFile, err
	}

	return configFile, os.WriteFile(configFile, []byte(sampleConfig(apiKey)), 0600)
}

func ResolveConfigFile() string {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(GetConfigFilePath(), "config.yaml")
		viper.SetConfigFile(configFile)
	}
	return configFile
}

func ConfigFileExists() (string, bool, error) {
	configFile := ResolveConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		return configFile, true, nil
	} else if !os.IsNotExist(err) {
		return configFile, false, err
	}
	return configFile, false, nil
}

func sampleConfig(apiKey string) string {
	return fmt.Sprintf(`# quail-cli configuration
# quail-cli stores API key and OAuth tokens in app.
app:
  api_key: %s
  access_token: ""
  expiry: ""
  refresh_token: ""
  token_type: ""
  user:
    id: 0
    name: ""
    bio: ""

post:
  # Map your Markdown frontmatter keys to Quaily post fields.
  # In this example, "featureImage" in frontmatter maps to "cover_image_url".
  frontmatter_mapping:
    cover_image_url: featureImage
`, strconv.Quote(apiKey))
}
