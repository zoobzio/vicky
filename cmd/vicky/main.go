// Package main is the entry point for the vicky server.
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zoobzio/aperture"
	"github.com/zoobzio/astql/postgres"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/cereal"
	grubminio "github.com/zoobzio/grub/minio"
	"github.com/zoobzio/rocco/session"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/api/auth"
	"github.com/zoobzio/vicky/api/capacitors"
	"github.com/zoobzio/vicky/config"
	"github.com/zoobzio/vicky/api/contracts"
	chunkerclient "github.com/zoobzio/vicky/external/chunker"
	embeddingclient "github.com/zoobzio/vicky/external/embedding"
	"github.com/zoobzio/vicky/external/github"
	indexerclient "github.com/zoobzio/vicky/external/indexer"
	"github.com/zoobzio/vicky/api/events"
	"github.com/zoobzio/vicky/api/handlers"
	"github.com/zoobzio/vicky/api/ingest"
	vickyauth "github.com/zoobzio/vicky/internal/auth"
	vickyotel "github.com/zoobzio/vicky/internal/otel"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
	"github.com/zoobzio/vicky/api/wire"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("vicky starting...")
	ctx := context.Background()

	// Initialize sum service and registry
	svc := sum.New()
	k := sum.Start()

	// Load all configs via fig
	if err := sum.Config[config.App](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}
	if err := sum.Config[config.Database](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}
	if err := sum.Config[config.Encryption](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load encryption config: %w", err)
	}
	if err := sum.Config[config.Redis](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load redis config: %w", err)
	}
	if err := sum.Config[config.Storage](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load storage config: %w", err)
	}
	if err := sum.Config[config.Observability](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load observability config: %w", err)
	}
	if err := sum.Config[config.Indexer](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load indexer config: %w", err)
	}
	if err := sum.Config[config.Chunker](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load chunker config: %w", err)
	}
	if err := sum.Config[config.Embedding](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load embedding config: %w", err)
	}
	if err := sum.Config[config.GitHub](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load github config: %w", err)
	}

	// Retrieve configs from registry
	appCfg := sum.MustUse[config.App](ctx)
	ghCfg := sum.MustUse[config.GitHub](ctx)
	dbCfg := sum.MustUse[config.Database](ctx)
	encCfg := sum.MustUse[config.Encryption](ctx)

	// Register AES encryptor
	keyBytes, err := hex.DecodeString(encCfg.Key)
	if err != nil {
		return fmt.Errorf("invalid encryption key: %w", err)
	}
	enc, err := cereal.AES(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to create encryptor: %w", err)
	}
	svc.WithEncryptor(cereal.EncryptAES, enc)

	// Connect to database
	db, err := sqlx.Connect("postgres", dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() { _ = db.Close() }()
	log.Println("database connected")

	capitan.Emit(ctx, events.StartupDatabaseConnected)

	// Connect to object storage
	storageCfg := sum.MustUse[config.Storage](ctx)
	minioClient, err := minio.New(storageCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(storageCfg.AccessKey, storageCfg.SecretKey, ""),
		Secure: storageCfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create minio client: %w", err)
	}

	bucketProvider := grubminio.New(minioClient, storageCfg.Bucket)
	log.Println("storage connected")

	capitan.Emit(ctx, events.StartupStorageConnected)

	// Create all stores
	allStores, err := stores.New(db, postgres.New(), bucketProvider)
	if err != nil {
		return fmt.Errorf("failed to create stores: %w", err)
	}

	// Register stores with service registry
	sum.Register[contracts.Users](k, allStores.Users)
	sum.Register[contracts.Repositories](k, allStores.Repositories)
	sum.Register[contracts.IngestionConfigs](k, allStores.IngestionConfigs)
	sum.Register[contracts.Versions](k, allStores.Versions)
	sum.Register[contracts.Jobs](k, allStores.Jobs)
	sum.Register[contracts.Documents](k, allStores.Documents)
	sum.Register[contracts.Chunks](k, allStores.Chunks)
	sum.Register[contracts.Symbols](k, allStores.Symbols)
	sum.Register[contracts.SCIPSymbols](k, allStores.SCIPSymbols)
	sum.Register[contracts.SCIPOccurrences](k, allStores.SCIPOccurrences)
	sum.Register[contracts.SCIPRelationships](k, allStores.SCIPRelationships)
	sum.Register[contracts.Sessions](k, allStores.Sessions)
	sum.Register[contracts.Blobs](k, allStores.Blobs)
	sum.Register[contracts.Keys](k, allStores.Keys)

	// Register external services
	sum.Register[contracts.GitHub](k, github.NewClient())

	indexerCfg := sum.MustUse[config.Indexer](ctx)
	indexerClient := indexerclient.NewClient(indexerCfg.Addresses())
	defer func() { _ = indexerClient.Close() }()
	sum.Register[contracts.Indexer](k, indexerClient)
	chunkerCfg := sum.MustUse[config.Chunker](ctx)
	chunkerClient := chunkerclient.NewClient(chunkerCfg.Addr)
	defer func() { _ = chunkerClient.Close() }()
	sum.Register[contracts.Chunker](k, chunkerClient)

	embeddingCfg := sum.MustUse[config.Embedding](ctx)
	embeddingClient, err := embeddingclient.NewClient(
		embeddingCfg.Provider, embeddingCfg.Model, embeddingCfg.APIKey, embeddingCfg.Dimensions,
	)
	if err != nil {
		return fmt.Errorf("failed to create embedding client: %w", err)
	}
	sum.Register[contracts.Embedder](k, embeddingClient)

	// Register model boundaries
	if _, err := sum.NewBoundary[models.User](k); err != nil {
		return fmt.Errorf("failed to register user boundary: %w", err)
	}

	// Register wire boundaries
	if err := wire.RegisterBoundaries(k); err != nil {
		return fmt.Errorf("failed to register boundaries: %w", err)
	}

	sum.Freeze(k)

	capitan.Emit(ctx, events.StartupServicesReady)

	// Initialize OTEL providers
	obsCfg := sum.MustUse[config.Observability](ctx)
	otelProviders, err := vickyotel.New(ctx, vickyotel.Config{
		Endpoint:    obsCfg.OTLPEndpoint,
		ServiceName: obsCfg.ServiceName,
	})
	if err != nil {
		return fmt.Errorf("failed to create otel providers: %w", err)
	}
	defer func() { _ = otelProviders.Shutdown(ctx) }()
	log.Println("observability initialized")

	capitan.Emit(ctx, events.StartupOTELReady)

	// Initialize aperture
	ap, err := aperture.New(
		capitan.Default(),
		otelProviders.Log,
		otelProviders.Metric,
		otelProviders.Trace,
	)
	if err != nil {
		return fmt.Errorf("failed to create aperture: %w", err)
	}
	defer ap.Close()

	capitan.Emit(ctx, events.StartupApertureReady)

	// Initialize capacitors (hot-reload config)
	if err := capacitors.Init(ctx, db.DB, dbCfg.DSN(), ap); err != nil {
		return fmt.Errorf("failed to initialize capacitors: %w", err)
	}
	log.Println("capacitors initialized")

	capitan.Emit(ctx, events.StartupCapacitorsReady)

	// Start ingestion worker pool
	ingestWorker := ingest.NewWorker()
	ingestWorker.Start(ctx)
	defer func() { _ = ingestWorker.Stop() }()

	capitan.Emit(ctx, events.StartupWorkerReady)

	// Create OAuth service
	oauthSvc, err := auth.NewOAuthService(ghCfg)
	if err != nil {
		return fmt.Errorf("failed to create oauth service: %w", err)
	}
	oauthSvc.SetUsers(allStores.Users)

	// Build session config
	sessionCfg := session.Config{
		OAuth:       oauthSvc.OAuthConfig(),
		Store:       allStores.Sessions,
		Cookie:      session.CookieConfig{SignKey: []byte(appCfg.SessionSignKey)},
		Resolve:     oauthSvc.Resolve,
		RedirectURL: appCfg.SessionRedirectURL,
	}

	// Create session handlers
	loginHandler, err := session.NewLoginHandler("/auth/github", sessionCfg)
	if err != nil {
		return fmt.Errorf("failed to create login handler: %w", err)
	}

	callbackHandler, err := session.NewCallbackHandler("/auth/github/callback", sessionCfg)
	if err != nil {
		return fmt.Errorf("failed to create callback handler: %w", err)
	}

	logoutHandler, err := session.NewLogoutHandler("/auth/logout", sessionCfg, appCfg.SessionRedirectURL)
	if err != nil {
		return fmt.Errorf("failed to create logout handler: %w", err)
	}

	// Tag auth handlers
	loginHandler.WithTags("Authentication")
	callbackHandler.WithTags("Authentication")
	logoutHandler.WithTags("Authentication")

	// Configure service
	svc.Engine().
		WithTag("Authentication", "OAuth login, callback, and logout").
		WithTag("Users", "User profile management").
		WithTag("Repositories", "Repository registration and management").
		WithTag("Versions", "Version ingestion and status tracking").
		WithTag("Search", "Semantic search across code and documentation").
		WithTag("Code Intelligence", "SCIP-powered definitions, references, and symbol navigation").
		WithTag("API Keys", "Programmatic authentication via API keys").
		WithTagGroup("Identity", "Authentication", "Users", "API Keys").
		WithTagGroup("Resources", "Repositories", "Versions").
		WithTagGroup("Intelligence", "Search", "Code Intelligence").
		WithAuthenticator(vickyauth.KeyExtractor(allStores.Keys, session.Extractor(allStores.Sessions, sessionCfg.Cookie)))
	svc.Handle(loginHandler, callbackHandler, logoutHandler)
	svc.Handle(handlers.All()...)

	capitan.Emit(ctx, events.StartupServerListening, events.StartupPortKey.Field(appCfg.Port))

	log.Printf("starting server on port %d...", appCfg.Port)
	return svc.Run("", appCfg.Port)
}
