package contracts

import (
	"context"

	"github.com/zoobzio/rocco/session"
)

// Sessions defines the contract for session and OAuth state storage.
// Implements session.Store from rocco/session.
type Sessions interface {
	session.Store
	// Cleanup removes expired states and sessions.
	Cleanup(ctx context.Context) error
}
