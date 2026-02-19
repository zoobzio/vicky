package handlers

import "github.com/zoobzio/rocco"

// Handler errors using rocco's built-in error types.
var (
	ErrRepositoryNotFound = rocco.ErrNotFound.WithMessage("repository not found or not registered")
	ErrVersionNotFound    = rocco.ErrNotFound.WithMessage("version not found or not ingested")
	ErrMissingQuery       = rocco.ErrBadRequest.WithMessage("query parameter 'q' is required")
	ErrInvalidLimit       = rocco.ErrBadRequest.WithMessage("limit must be between 1 and 100")
	ErrMissingSymbol      = rocco.ErrBadRequest.WithMessage("query parameter 'symbol' is required")
	ErrKeyNotFound        = rocco.ErrNotFound.WithMessage("api key not found")
	ErrKeyForbidden       = rocco.ErrForbidden.WithMessage("api key belongs to another user")
)
