package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// RepositoryEvent is emitted for repository lifecycle events.
type RepositoryEvent struct {
	RepositoryID  int64  `json:"repository_id"`
	GitHubID      int64  `json:"github_id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
	UserID        int64  `json:"user_id"`
}

// VersionEvent is emitted for version lifecycle events.
type VersionEvent struct {
	VersionID      int64  `json:"version_id"`
	RepositoryID   int64  `json:"repository_id"`
	Tag            string `json:"tag"`
	CommitSHA      string `json:"commit_sha"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previous_status,omitempty"`
	UserID         int64  `json:"user_id,omitempty"`
	Error          string `json:"error,omitempty"`
}

// Repository signals.
var (
	RepositoryRegisteredSignal = capitan.NewSignal("vicky.repository.registered", "Repository registered")
	RepositoryUpdatedSignal    = capitan.NewSignal("vicky.repository.updated", "Repository updated")
	RepositoryDeletedSignal    = capitan.NewSignal("vicky.repository.deleted", "Repository deleted")
	RepositorySyncedSignal     = capitan.NewSignal("vicky.repository.synced", "Repository synced with GitHub")
)

// Version signals.
var (
	VersionCreatedSignal = capitan.NewSignal("vicky.version.created", "Version created")
	VersionUpdatedSignal = capitan.NewSignal("vicky.version.updated", "Version status updated")
	VersionDeletedSignal = capitan.NewSignal("vicky.version.deleted", "Version deleted")
	VersionReadySignal   = capitan.NewSignal("vicky.version.ready", "Version ready for search")
	VersionFailedSignal  = capitan.NewSignal("vicky.version.failed", "Version ingestion failed")
)

// Repository provides access to repository lifecycle events.
var Repository = struct {
	Registered sum.Event[RepositoryEvent]
	Updated    sum.Event[RepositoryEvent]
	Deleted    sum.Event[RepositoryEvent]
	Synced     sum.Event[RepositoryEvent]
}{
	Registered: sum.NewInfoEvent[RepositoryEvent](RepositoryRegisteredSignal),
	Updated:    sum.NewInfoEvent[RepositoryEvent](RepositoryUpdatedSignal),
	Deleted:    sum.NewInfoEvent[RepositoryEvent](RepositoryDeletedSignal),
	Synced:     sum.NewInfoEvent[RepositoryEvent](RepositorySyncedSignal),
}

// Version provides access to version lifecycle events.
var Version = struct {
	Created sum.Event[VersionEvent]
	Updated sum.Event[VersionEvent]
	Deleted sum.Event[VersionEvent]
	Ready   sum.Event[VersionEvent]
	Failed  sum.Event[VersionEvent]
}{
	Created: sum.NewInfoEvent[VersionEvent](VersionCreatedSignal),
	Updated: sum.NewInfoEvent[VersionEvent](VersionUpdatedSignal),
	Deleted: sum.NewInfoEvent[VersionEvent](VersionDeletedSignal),
	Ready:   sum.NewInfoEvent[VersionEvent](VersionReadySignal),
	Failed:  sum.NewErrorEvent[VersionEvent](VersionFailedSignal),
}
