package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Config struct {
	JWT    JWTConfig    `yaml:"jwt"`
	OAuth2 OAuth2Config `yaml:"oauth2"`
	APIKey APIKeyConfig `yaml:"api_key"`
}

type JWTConfig struct {
	Token  string `yaml:"token"`
	Header string `yaml:"header"`
	Scheme string `yaml:"scheme"`
}

type OAuth2Config struct {
	TokenURL     string `yaml:"token_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Scope        string `yaml:"scope"`
}

type APIKeyConfig struct {
	Key        string `yaml:"key"`
	Header     string `yaml:"header"`
	QueryParam string `yaml:"query_param"`
}

func (c Config) IsZero() bool {
	return c.JWT.Token == "" && c.OAuth2.TokenURL == "" && c.APIKey.Key == ""
}

type Middleware struct {
	cfg        Config
	mu         sync.Mutex
	o2token    string
	o2expiry   time.Time
	httpClient *http.Client
}

func New(cfg Config) (*Middleware, error) {
	m := &Middleware{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	if cfg.OAuth2.TokenURL != "" {
		if err := m.refreshOAuth2(context.Background()); err != nil {
			return nil, fmt.Errorf("auth: initial oauth2 token fetch failed: %w", err)
		}
	}
	return m, nil
}

func (m *Middleware) Inject(req *http.Request) error {
	switch {
	case m.cfg.JWT.Token != "":
		return m.injectJWT(req, &m.cfg.JWT)
	case m.cfg.OAuth2.TokenURL != "":
		return m.injectOAuth2(req)
	case m.cfg.APIKey.Key != "":
		return m.injectAPIKey(req, &m.cfg.APIKey)
	}
	return nil
}

func (m *Middleware) injectJWT(req *http.Request, cfg *JWTConfig) error {
	if cfg.Token == "" {
		return fmt.Errorf("auth: jwt.token is empty")
	}
	header := cfg.Header
	if header == "" {
		header = "Authorization"
	}
	scheme := cfg.Scheme
	if scheme == "" {
		scheme = "Bearer"
	}
	req.Header.Set(header, scheme+" "+cfg.Token)
	return nil
}

func (m *Middleware) injectOAuth2(req *http.Request) error {
	token, err := m.getOAuth2Token(req.Context())
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (m *Middleware) getOAuth2Token(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.o2token != "" && time.Now().Add(30*time.Second).Before(m.o2expiry) {
		return m.o2token, nil
	}
	if err := m.refreshOAuth2(ctx); err != nil {
		return "", err
	}
	return m.o2token, nil
}

func (m *Middleware) refreshOAuth2(ctx context.Context) error {
	cfg := m.cfg.OAuth2
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	if cfg.Scope != "" {
		data.Set("scope", cfg.Scope)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.TokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("oauth2 token request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("oauth2 token request returned %d: %s", resp.StatusCode, body)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("oauth2 token parse: %w", err)
	}
	if result.AccessToken == "" {
		return fmt.Errorf("oauth2: empty access_token in response")
	}

	m.o2token = result.AccessToken
	if result.ExpiresIn > 0 {
		m.o2expiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	} else {
		m.o2expiry = time.Now().Add(1 * time.Hour)
	}
	return nil
}

func (m *Middleware) injectAPIKey(req *http.Request, cfg *APIKeyConfig) error {
	if cfg.Key == "" {
		return fmt.Errorf("auth: api_key.key is empty")
	}
	if cfg.QueryParam != "" {
		q := req.URL.Query()
		q.Set(cfg.QueryParam, cfg.Key)
		req.URL.RawQuery = q.Encode()
		return nil
	}
	header := cfg.Header
	if header == "" {
		header = "X-API-Key"
	}
	req.Header.Set(header, cfg.Key)
	return nil
}
