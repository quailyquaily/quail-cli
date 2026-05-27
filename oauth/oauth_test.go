package oauth

import (
	"runtime"
	"strings"
	"testing"
)

func TestOAuthRedirectURL(t *testing.T) {
	got := oauthRedirectURL("https://quaily.com/")
	want := "https://quaily.com/oauth/code"
	if got != want {
		t.Fatalf("oauthRedirectURL() = %q, want %q", got, want)
	}
}

func TestReadAuthorizationCode(t *testing.T) {
	code, err := readAuthorizationCode(strings.NewReader(" abc123 \n"))
	if err != nil {
		t.Fatalf("readAuthorizationCode() error = %v", err)
	}
	if code != "abc123" {
		t.Fatalf("readAuthorizationCode() = %q, want %q", code, "abc123")
	}
}

func TestShouldOpenBrowserDisabledByEnv(t *testing.T) {
	t.Setenv("QUAIL_CLI_NO_BROWSER", "1")
	if shouldOpenBrowser() {
		t.Fatal("expected browser auto-open to be disabled")
	}
}

func TestShouldOpenBrowserOnHeadlessLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-specific browser detection")
	}

	t.Setenv("QUAIL_CLI_NO_BROWSER", "")
	t.Setenv("DISPLAY", "")
	t.Setenv("WAYLAND_DISPLAY", "")
	t.Setenv("BROWSER", "")

	if shouldOpenBrowser() {
		t.Fatal("expected browser auto-open to be disabled without a graphical browser")
	}
}

func TestShouldOpenBrowserWithLinuxBrowserEnv(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-specific browser detection")
	}

	t.Setenv("QUAIL_CLI_NO_BROWSER", "")
	t.Setenv("DISPLAY", "")
	t.Setenv("WAYLAND_DISPLAY", "")
	t.Setenv("BROWSER", "firefox")

	if !shouldOpenBrowser() {
		t.Fatal("expected browser auto-open with BROWSER set")
	}
}
