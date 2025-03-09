package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/lyricat/goutils/uuid"
	"golang.org/x/oauth2"
)

const (
	authPath     = "/oauth/authorize"
	tokenPath    = "/oauth/token"
	redirectURL  = "http://localhost:63812/oauth/code"
	clientID     = "e9139b6e-298a-43e4-91f0-fc97960e281a"
	clientSecret = ""
)

func Login(authBase, apiBase string) (*oauth2.Token, error) {
	state := uuid.New()

	verifier := generateCodeVerifier()
	challenge := verifier

	authURL := fmt.Sprintf("%s%s", authBase, authPath)
	tokenURL := fmt.Sprintf("%s%s", authBase, tokenPath)

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

	fmt.Printf("Please visit this URL to authorize the application: %v\n", authCodeURL)

	// start a local server to handle the redirect
	codeChan := make(chan string)
	go func() {
		http.HandleFunc("/oauth/code", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			returnedState := r.URL.Query().Get("state")

			if returnedState != state {
				slog.Error("state mismatch", "expected", state, "got", returnedState)
				fmt.Fprintf(w, "Error: state mismatch")
				codeChan <- ""
				return
			}

			fmt.Fprintf(w, "Authorization successful! You can close this window.")
			codeChan <- code
		})

		http.ListenAndServe(":63812", nil)
	}()

	code := <-codeChan

	if code == "" {
		return nil, fmt.Errorf("failed to get authorization code")
	}

	token, err := exchangeCodeForToken(apiBase, code, verifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, nil
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

func exchangeCodeForToken(apiBase, code, verifier string) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURL)
	data.Set("client_id", clientID)
	data.Set("code_verifier", verifier)

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
