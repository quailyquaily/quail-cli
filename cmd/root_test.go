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
