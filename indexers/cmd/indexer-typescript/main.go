// Package main is the entry point for the TypeScript SCIP indexer service.
package main

import (
	"context"
	"log"
	"os"

	"github.com/zoobzio/vicky/indexers"
	"github.com/zoobzio/vicky/internal/otel"
)

func main() {
	ctx := context.Background()

	// Initialize observability
	obs, err := otel.NewSidecar(ctx, otel.SidecarConfig{
		ServiceName: "vicky-indexer-typescript",
	})
	if err != nil {
		log.Fatalf("create observability: %v", err)
	}
	defer obs.Shutdown(ctx)

	cfg := indexers.StorageConfig{
		Endpoint:  getEnv("INDEXER_STORAGE_ENDPOINT", "localhost:9000"),
		AccessKey: getEnv("INDEXER_STORAGE_ACCESS_KEY", "vicky"),
		SecretKey: getEnv("INDEXER_STORAGE_SECRET_KEY", "vickydev"),
		Bucket:    getEnv("INDEXER_STORAGE_BUCKET", "vicky"),
		UseSSL:    os.Getenv("INDEXER_STORAGE_USE_SSL") == "true",
	}

	storage, err := indexers.NewStorage(cfg)
	if err != nil {
		log.Fatalf("create storage: %v", err)
	}

	srv := indexers.NewServer(storage, &indexers.TypeScriptExecutor{}, "typescript")

	addr := getEnv("INDEXER_LISTEN_ADDR", ":9090")
	if err := indexers.ListenAndServe(ctx, addr, srv); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
