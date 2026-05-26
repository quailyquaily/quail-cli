package oauth

import (
	"runtime"
	"testing"
)

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
