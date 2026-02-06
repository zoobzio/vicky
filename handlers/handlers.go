// Package handlers provides HTTP handlers for the vicky API.
package handlers

import "github.com/zoobzio/rocco"

// All returns all API handlers for registration with rocco.
func All() []rocco.Endpoint {
	return []rocco.Endpoint{
		// Users
		GetMe,
		UpdateMe,

		// Repositories
		ListRepositories,
		RegisterRepository,
		GetRepository,

		// Versions
		ListVersions,
		GetVersion,
		TriggerIngest,

		// Search
		SearchChunks,
		SearchSymbols,
		FindSimilarDocuments,

		// Code Intelligence
		GetDefinition,
		FindReferences,
		FindImplementations,
		ListSymbols,
	}
}
