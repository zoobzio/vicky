// Package main is the entry point for the vicky admin service.
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zoobzio/astql/postgres"
	"github.com/zoobzio/cereal"
	"github.com/zoobzio/rocco/session"
	"github.com/zoobzio/sum"
	admincontracts "github.com/zoobzio/vicky/admin/contracts"
	adminhandlers "github.com/zoobzio/vicky/admin/handlers"
	"github.com/zoobzio/vicky/config"
	"github.com/zoobzio/vicky/internal/auth"
	vickyotel "github.com/zoobzio/vicky/internal/otel"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/stores"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("vicky-admin starting...")
	ctx := context.Background()

	// Initialize sum service and registry
	svc := sum.New()
	k := sum.Start()

	// Load configs
	if err := sum.Config[config.Admin](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load admin config: %w", err)
	}
	if err := sum.Config[config.Database](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}
	if err := sum.Config[config.Encryption](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load encryption config: %w", err)
	}
	if err := sum.Config[config.Observability](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load observability config: %w", err)
	}
	if err := sum.Config[config.GitHub](ctx, k, nil); err != nil {
		return fmt.Errorf("failed to load github config: %w", err)
	}

	// Retrieve configs
	adminCfg := sum.MustUse[config.Admin](ctx)
	dbCfg := sum.MustUse[config.Database](ctx)
	encCfg := sum.MustUse[config.Encryption](ctx)
	obsCfg := sum.MustUse[config.Observability](ctx)
	ghCfg := sum.MustUse[config.GitHub](ctx)

	// Register AES encryptor (needed for User model)
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

	// Create users store
	usersStore, err := stores.NewUsers(db, postgres.New())
	if err != nil {
		return fmt.Errorf("failed to create users store: %w", err)
	}

	// Create repositories store
	reposStore, err := stores.NewRepositories(db, postgres.New())
	if err != nil {
		return fmt.Errorf("failed to create repositories store: %w", err)
	}

	// Create jobs store
	jobsStore, err := stores.NewJobs(db, postgres.New())
	if err != nil {
		return fmt.Errorf("failed to create jobs store: %w", err)
	}

	// Register stores against admin contracts
	sum.Register[admincontracts.Users](k, usersStore)
	sum.Register[admincontracts.Repositories](k, reposStore)
	sum.Register[admincontracts.Jobs](k, jobsStore)

	// Register model boundaries (User needs encryption/decryption)
	if _, err := sum.NewBoundary[models.User](k); err != nil {
		return fmt.Errorf("failed to register user boundary: %w", err)
	}

	sum.Freeze(k)

	// Initialize OTEL
	otelProviders, err := vickyotel.New(ctx, vickyotel.Config{
		Endpoint:    obsCfg.OTLPEndpoint,
		ServiceName: "vicky-admin",
	})
	if err != nil {
		return fmt.Errorf("failed to create otel providers: %w", err)
	}
	defer func() { _ = otelProviders.Shutdown(ctx) }()
	log.Println("observability initialized")

	// Create sessions store for admin
	sessionsStore := stores.NewSessions(db)

	// Create admin OAuth service
	oauthSvc, err := auth.NewAdminOAuthService(adminCfg, ghCfg)
	if err != nil {
		return fmt.Errorf("failed to create admin oauth service: %w", err)
	}

	// Build session config
	sessionCfg := session.Config{
		OAuth:       oauthSvc.OAuthConfig(),
		Store:       sessionsStore,
		Cookie:      session.CookieConfig{SignKey: []byte(adminCfg.SessionSignKey)},
		Resolve:     oauthSvc.Resolve,
		RedirectURL: "/admin",
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
	logoutHandler, err := session.NewLogoutHandler("/auth/logout", sessionCfg, "/admin")
	if err != nil {
		return fmt.Errorf("failed to create logout handler: %w", err)
	}

	// Configure engine with authenticator
	svc.Engine().
		WithTag("Authentication", "Admin OAuth login").
		WithTag("Admin", "Administrative operations").
		WithAuthenticator(session.Extractor(sessionsStore, sessionCfg.Cookie))
	svc.Handle(loginHandler, callbackHandler, logoutHandler)

	// Register admin handlers (all require authentication)
	svc.Handle(adminhandlers.All()...)
	log.Println("github authentication configured")

	log.Printf("starting admin server on port %d...", adminCfg.Port)
	return svc.Run("", adminCfg.Port)
}
