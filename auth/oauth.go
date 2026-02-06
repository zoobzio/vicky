// Package auth provides GitHub OAuth integration for vicky.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zoobzio/rocco/oauth"
	"github.com/zoobzio/rocco/session"
	"github.com/zoobzio/vicky/config"
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/models"
)

// OAuthService manages GitHub user resolution during OAuth login.
type OAuthService struct {
	cfg        config.GitHub
	users      contracts.Users
	oauthCfg   oauth.Config
	httpClient *http.Client
}

// NewOAuthService creates a new OAuth service for GitHub authentication.
func NewOAuthService(cfg config.GitHub) (*OAuthService, error) {
	if cfg.ClientID == "" {
		return nil, errors.New("auth: GITHUB_CLIENT_ID is required")
	}
	if cfg.ClientSecret == "" {
		return nil, errors.New("auth: GITHUB_CLIENT_SECRET is required")
	}
	if cfg.RedirectURI == "" {
		return nil, errors.New("auth: GITHUB_REDIRECT_URI is required")
	}

	oauthCfg := oauth.GitHub()
	oauthCfg.ClientID = cfg.ClientID
	oauthCfg.ClientSecret = cfg.ClientSecret
	oauthCfg.RedirectURI = cfg.RedirectURI
	oauthCfg.Scopes = []string{"read:user", "user:email"}

	return &OAuthService{
		cfg:        cfg,
		oauthCfg:   oauthCfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// SetUsers sets the users contract. Called after service registry is frozen.
func (s *OAuthService) SetUsers(users contracts.Users) {
	s.users = users
}

// OAuthConfig returns the OAuth config for use with session.Config.
func (s *OAuthService) OAuthConfig() oauth.Config {
	return s.oauthCfg
}

// Resolve maps OAuth tokens to session data. It fetches GitHub user info,
// validates org membership, upserts the user record, and returns session data.
func (s *OAuthService) Resolve(ctx context.Context, tokens *oauth.TokenResponse) (*session.Data, error) {
	ghUser, err := s.fetchGitHubUser(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Upsert user
	user := &models.User{
		ID:          ghUser.ID,
		Login:       ghUser.Login,
		Email:       ghUser.Email,
		Name:        strPtr(ghUser.Name),
		AvatarURL:   strPtr(ghUser.AvatarURL),
		AccessToken: tokens.AccessToken,
		LastLoginAt: time.Now(),
	}

	userID := strconv.FormatInt(ghUser.ID, 10)
	if err := s.users.Set(ctx, userID, user); err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	return &session.Data{
		UserID: userID,
		Email:  ghUser.Email,
		Meta: map[string]any{
			"login":      ghUser.Login,
			"avatar_url": ghUser.AvatarURL,
		},
	}, nil
}

// GitHubUser represents the response from GET /user.
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// fetchGitHubUser fetches user info from GitHub API.
func (s *OAuthService) fetchGitHubUser(ctx context.Context, token string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// strPtr returns a pointer to s, or nil if s is empty.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
