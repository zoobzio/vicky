// Package auth provides authentication for vicky services.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/zoobzio/rocco/oauth"
	"github.com/zoobzio/rocco/session"
	"github.com/zoobzio/vicky/config"
)

// ErrNotAuthorized is returned when a user is not in the admin allowlist.
var ErrNotAuthorized = errors.New("user not authorized for admin access")

// AdminOAuthService manages GitHub OAuth for admin users.
type AdminOAuthService struct {
	adminCfg   config.Admin
	ghCfg      config.GitHub
	oauthCfg   oauth.Config
	httpClient *http.Client
}

// NewAdminOAuthService creates a new OAuth service for admin authentication.
func NewAdminOAuthService(adminCfg config.Admin, ghCfg config.GitHub) (*AdminOAuthService, error) {
	if ghCfg.ClientID == "" {
		return nil, errors.New("auth: GITHUB_CLIENT_ID is required")
	}
	if ghCfg.ClientSecret == "" {
		return nil, errors.New("auth: GITHUB_CLIENT_SECRET is required")
	}
	if ghCfg.RedirectURI == "" {
		return nil, errors.New("auth: GITHUB_REDIRECT_URI is required")
	}

	oauthCfg := oauth.GitHub()
	oauthCfg.ClientID = ghCfg.ClientID
	oauthCfg.ClientSecret = ghCfg.ClientSecret
	oauthCfg.RedirectURI = ghCfg.RedirectURI
	oauthCfg.Scopes = []string{"read:user", "user:email"}

	return &AdminOAuthService{
		adminCfg:   adminCfg,
		ghCfg:      ghCfg,
		oauthCfg:   oauthCfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// OAuthConfig returns the OAuth config for use with session.Config.
func (s *AdminOAuthService) OAuthConfig() oauth.Config {
	return s.oauthCfg
}

// Resolve maps OAuth tokens to session data. Returns ErrNotAuthorized if the
// user's GitHub login is not in the admin allowlist.
func (s *AdminOAuthService) Resolve(ctx context.Context, tokens *oauth.TokenResponse) (*session.Data, error) {
	ghUser, err := s.fetchGitHubUser(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Check if user is in admin allowlist
	if !s.adminCfg.IsUserAllowed(ghUser.Login) {
		return nil, ErrNotAuthorized
	}

	return &session.Data{
		UserID: fmt.Sprintf("%d", ghUser.ID),
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
func (s *AdminOAuthService) fetchGitHubUser(ctx context.Context, token string) (*GitHubUser, error) {
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
