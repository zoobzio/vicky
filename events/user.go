package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// UserEvent is emitted for user lifecycle events.
type UserEvent struct {
	UserID    int64  `json:"user_id"`
	GitHubID  int64  `json:"github_id"`
	Login     string `json:"login"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	OrgCount  int    `json:"org_count,omitempty"`
	TeamCount int    `json:"team_count,omitempty"`
}

// TokenEvent is emitted for token-related events.
type TokenEvent struct {
	UserID int64    `json:"user_id"`
	Login  string   `json:"login"`
	Scopes []string `json:"scopes,omitempty"`
}

// User signals.
var (
	UserCreatedSignal        = capitan.NewSignal("vicky.user.created", "User created on first login")
	UserUpdatedSignal        = capitan.NewSignal("vicky.user.updated", "User profile updated")
	UserLoggedInSignal       = capitan.NewSignal("vicky.user.logged_in", "User authenticated")
	UserTokenRefreshedSignal = capitan.NewSignal("vicky.user.token_refreshed", "User access token refreshed")
)

// User provides access to user lifecycle events.
var User = struct {
	Created        sum.Event[UserEvent]
	Updated        sum.Event[UserEvent]
	LoggedIn       sum.Event[UserEvent]
	TokenRefreshed sum.Event[TokenEvent]
}{
	Created:        sum.NewInfoEvent[UserEvent](UserCreatedSignal),
	Updated:        sum.NewInfoEvent[UserEvent](UserUpdatedSignal),
	LoggedIn:       sum.NewInfoEvent[UserEvent](UserLoggedInSignal),
	TokenRefreshed: sum.NewInfoEvent[TokenEvent](UserTokenRefreshedSignal),
}
