// Package main is the entry point for the chunker service.
package main

import (
	"context"
	"log"
	"os"

	"github.com/zoobzio/chisel/golang"
	"github.com/zoobzio/chisel/markdown"
	"github.com/zoobzio/chisel/typescript"
	"github.com/zoobzio/vicky/chunkers"
	"github.com/zoobzio/vicky/internal/otel"
)

func main() {
	ctx := context.Background()

	// Initialize observability
	obs, err := otel.NewSidecar(ctx, otel.SidecarConfig{
		ServiceName: "vicky-chunker",
	})
	if err != nil {
		log.Fatalf("create observability: %v", err)
	}
	defer obs.Shutdown(ctx)

	srv := chunkers.NewServer(
		golang.New(),
		typescript.New(),
		typescript.NewJavaScript(),
		markdown.New(),
	)

	addr := getEnv("CHUNKER_LISTEN_ADDR", ":9091")
	if err := chunkers.ListenAndServe(ctx, addr, srv); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
