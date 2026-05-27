package oauth

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/lyricat/goutils/uuid"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

const (
	authPath     = "/oauth/authorize"
	tokenPath    = "/oauth/token"
	redirectPath = "/oauth/code"
	clientID     = "e9139b6e-298a-43e4-91f0-fc97960e281a"
	clientSecret = ""
)

func Login(authBase, apiBase string) (*oauth2.Token, string, error) {
	state := uuid.New()

	verifier := generateCodeVerifier()
	challenge := verifier

	authBase = strings.TrimRight(authBase, "/")
	apiBase = strings.TrimRight(apiBase, "/")
	authURL := authBase + authPath
	tokenURL := authBase + tokenPath
	redirectURL := oauthRedirectURL(authBase)

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user.full", "post.write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURL,
	}

	authCodeURL := conf.AuthCodeURL(state, oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"))

	fmt.Printf("Please visit this URL to authorize the application:\n%v\n", authCodeURL)
	fmt.Printf("After authorization, copy the code from the browser page and paste it here.\n")
	fmt.Printf("The browser page also shows state. It should match: %s\n", state)

	if shouldOpenBrowser() {
		if err := browser.OpenURL(authCodeURL); err != nil {
			slog.Warn("failed to open browser automatically; open the URL manually", "error", err)
		}
	} else {
		fmt.Println("No graphical browser detected. Open the URL manually.")
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, authCodeURL, fmt.Errorf("authorization code input requires an interactive terminal")
	}

	code, err := readAuthorizationCode(os.Stdin)
	if err != nil {
		return nil, authCodeURL, err
	}
	if code == "" {
		return nil, authCodeURL, fmt.Errorf("failed to get authorization code")
	}

	token, err := exchangeCodeForToken(apiBase, code, verifier, redirectURL)
	if err != nil {
		return nil, authCodeURL, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, authCodeURL, nil
}

func oauthRedirectURL(authBase string) string {
	return strings.TrimRight(authBase, "/") + redirectPath
}

func readAuthorizationCode(in io.Reader) (string, error) {
	fmt.Print("Authorization code: ")
	reader := bufio.NewReader(in)
	code, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func shouldOpenBrowser() bool {
	if strings.TrimSpace(os.Getenv("QUAIL_CLI_NO_BROWSER")) != "" {
		return false
	}
	if runtime.GOOS != "linux" {
		return true
	}
	return os.Getenv("DISPLAY") != "" ||
		os.Getenv("WAYLAND_DISPLAY") != "" ||
		os.Getenv("BROWSER") != ""
}

func RefreshToken(apiBase, refreshToken string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)

	tokenURL := fmt.Sprintf("%s%s", apiBase, tokenPath)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func exchangeCodeForToken(apiBase, code, verifier, redirectURL string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURL)
	data.Set("client_id", clientID)
	data.Set("code_verifier", verifier)

	tokenURL := strings.TrimRight(apiBase, "/") + tokenPath

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
