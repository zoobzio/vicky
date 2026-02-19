package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// ProviderEvent is emitted for embedding provider configuration events.
type ProviderEvent struct {
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
	Previous   string `json:"previous,omitempty"`
	Error      string `json:"error,omitempty"`
}

// Embedding provider signals.
var (
	EmbeddingProviderChangedSignal = capitan.NewSignal("vicky.embedding.provider.changed", "Embedding provider configuration changed")
	EmbeddingProviderErrorSignal   = capitan.NewSignal("vicky.embedding.provider.error", "Embedding provider error")
)

// Embedding provides access to embedding configuration events.
// Note: Runtime embedding events (request/complete/fail) are emitted by vex directly.
var Embedding = struct {
	ProviderChanged sum.Event[ProviderEvent]
	ProviderError   sum.Event[ProviderEvent]
}{
	ProviderChanged: sum.NewInfoEvent[ProviderEvent](EmbeddingProviderChangedSignal),
	ProviderError:   sum.NewErrorEvent[ProviderEvent](EmbeddingProviderErrorSignal),
}
