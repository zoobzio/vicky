package events

import "github.com/zoobzio/capitan"

// Startup signals for server lifecycle.
// These are direct capitan signals (not sum.Event) since they're
// operational events, not domain lifecycle events for consumers.
var (
	StartupDatabaseConnected = capitan.NewSignal("vicky.startup.database.connected", "Database connection established")
	StartupStorageConnected  = capitan.NewSignal("vicky.startup.storage.connected", "Object storage connection established")
	StartupServicesReady     = capitan.NewSignal("vicky.startup.services.ready", "All services registered")
	StartupOTELReady         = capitan.NewSignal("vicky.startup.otel.ready", "OpenTelemetry providers initialized")
	StartupApertureReady     = capitan.NewSignal("vicky.startup.aperture.ready", "Aperture observability bridge initialized")
	StartupCapacitorsReady   = capitan.NewSignal("vicky.startup.capacitors.ready", "Hot-reload capacitors initialized")
	StartupWorkerReady       = capitan.NewSignal("vicky.startup.worker.ready", "Ingestion worker pool started")
	StartupServerListening   = capitan.NewSignal("vicky.startup.server.listening", "HTTP server listening")
	StartupFailed            = capitan.NewSignal("vicky.startup.failed", "Server startup failed")
)

// Startup field keys for direct emission.
var (
	StartupPortKey    = capitan.NewIntKey("port")
	StartupWorkersKey = capitan.NewIntKey("workers")
	StartupErrorKey   = capitan.NewErrorKey("error")
)
