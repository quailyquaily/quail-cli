package cmd

import (
	"os"
	"testing"
)

func TestIsSetupCommand(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "init command",
			args: []string{"quail-cli", "init"},
			want: true,
		},
		{
			name: "login command after config flag",
			args: []string{"quail-cli", "--config", "./config.yaml", "login"},
			want: true,
		},
		{
			name: "init command after config equals flag",
			args: []string{"quail-cli", "--config=./config.yaml", "init"},
			want: true,
		},
		{
			name: "version command",
			args: []string{"quail-cli", "--config", "./missing.yaml", "version"},
			want: false,
		},
		{
			name: "init as command argument",
			args: []string{"quail-cli", "post", "upsert", "init"},
			want: false,
		},
		{
			name: "init after double dash",
			args: []string{"quail-cli", "--", "init"},
			want: false,
		},
	}

	for _, tt := range tests {
		os.Args = tt.args
		if got := isSetupCommand(); got != tt.want {
			t.Fatalf("%s: isSetupCommand() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCommandName(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "version skips config flag",
			args: []string{"quail-cli", "--config", "./missing.yaml", "version"},
			want: "version",
		},
		{
			name: "version skips config equals flag",
			args: []string{"quail-cli", "--config=./missing.yaml", "version"},
			want: "version",
		},
		{
			name: "json before command",
			args: []string{"quail-cli", "--json", "reader", "subscriptions"},
			want: "reader",
		},
		{
			name: "double dash stops command parsing",
			args: []string{"quail-cli", "--", "version"},
			want: "",
		},
	}

	for _, tt := range tests {
		os.Args = tt.args
		if got := commandName(); got != tt.want {
			t.Fatalf("%s: commandName() = %q, want %q", tt.name, got, tt.want)
		}
	}
}
