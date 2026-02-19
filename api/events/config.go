package events

import (
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/sum"
)

// ConfigEvent is emitted for configuration lifecycle events.
type ConfigEvent struct {
	Domain        string `json:"domain"`
	ChangedFields string `json:"changed_fields,omitempty"`
	Error         string `json:"error,omitempty"`
}

// Config signals.
var (
	ConfigLoadedSignal           = capitan.NewSignal("vicky.config.loaded", "Configuration loaded from database")
	ConfigChangedSignal          = capitan.NewSignal("vicky.config.changed", "Configuration changed via hot-reload")
	ConfigErrorSignal            = capitan.NewSignal("vicky.config.error", "Configuration error")
	ConfigValidationFailedSignal = capitan.NewSignal("vicky.config.validation_failed", "Configuration validation failed")
)

// Config provides access to configuration lifecycle events.
var Config = struct {
	Loaded           sum.Event[ConfigEvent]
	Changed          sum.Event[ConfigEvent]
	Error            sum.Event[ConfigEvent]
	ValidationFailed sum.Event[ConfigEvent]
}{
	Loaded:           sum.NewInfoEvent[ConfigEvent](ConfigLoadedSignal),
	Changed:          sum.NewInfoEvent[ConfigEvent](ConfigChangedSignal),
	Error:            sum.NewErrorEvent[ConfigEvent](ConfigErrorSignal),
	ValidationFailed: sum.NewWarnEvent[ConfigEvent](ConfigValidationFailedSignal),
}
