package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestWriteSampleConfigWithAPIKey(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	viper.SetConfigFile(configFile)

	got, err := WriteSampleConfig(" QK-test ")
	if err != nil {
		t.Fatalf("WriteSampleConfig() error = %v", err)
	}
	if got != configFile {
		t.Fatalf("WriteSampleConfig() path = %q, want %q", got, configFile)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), `api_key: "QK-test"`) {
		t.Fatalf("generated config does not contain trimmed API key:\n%s", data)
	}

	got, exists, err := ConfigFileExists()
	if err != nil {
		t.Fatalf("ConfigFileExists() error = %v", err)
	}
	if got != configFile || !exists {
		t.Fatalf("ConfigFileExists() = %q, %v; want %q, true", got, exists, configFile)
	}

	_, err = WriteSampleConfig("QK-test")
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("WriteSampleConfig() second error = %v, want os.ErrExist", err)
	}
}

func TestWriteSampleConfigRejectsInvalidAPIKey(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	viper.SetConfigFile(configFile)

	if _, err := WriteSampleConfig("not-a-quaily-key"); err == nil {
		t.Fatal("WriteSampleConfig() error = nil, want invalid API key error")
	}
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Fatalf("Stat() error = %v, want os.IsNotExist", err)
	}
}
