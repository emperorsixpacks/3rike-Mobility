package canton

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenProvider fetches and caches a Keycloak OIDC bearer token.
// It auto-refreshes the token 60 seconds before expiry.
type TokenProvider struct {
	tokenURL string
	clientID string
	username string
	password string
	scope    string

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// NewTokenProvider creates a provider using the password grant flow.
func NewTokenProvider(tokenURL, clientID, username, password string) *TokenProvider {
	return &TokenProvider{
		tokenURL: tokenURL,
		clientID: clientID,
		username: username,
		password: password,
		scope:    "openid daml_ledger_api offline_access",
	}
}

// Token returns a valid bearer token, fetching a new one if needed.
func (p *TokenProvider) Token() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Before(p.expiresAt.Add(-60*time.Second)) {
		return p.token, nil
	}
	return p.fetch()
}

func (p *TokenProvider) fetch() (string, error) {
	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {p.clientID},
		"username":   {p.username},
		"password":   {p.password},
		"scope":      {p.scope},
	}
	resp, err := http.Post(p.tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())) //nolint:noctx
	if err != nil {
		return "", fmt.Errorf("token fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("token fetch: HTTP %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("token decode: %w", err)
	}

	p.token = body.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(body.ExpiresIn) * time.Second)
	return p.token, nil
}
