package stores

import (
	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/grub"
)

// Stores provides access to all data stores.
type Stores struct {
	Users             *Users
	Repositories      *Repositories
	IngestionConfigs  *IngestionConfigs
	Versions          *Versions
	Jobs              *Jobs
	Documents         *Documents
	Chunks            *Chunks
	Symbols           *Symbols
	SCIPSymbols       *SCIPSymbols
	SCIPOccurrences   *SCIPOccurrences
	SCIPRelationships *SCIPRelationships
	Sessions          *Sessions
	Blobs             *Blobs
	Keys              *Keys
}

// New creates all stores with the given database connection.
func New(db *sqlx.DB, renderer astql.Renderer, bucket grub.BucketProvider) (*Stores, error) {
	users, err := NewUsers(db, renderer)
	if err != nil {
		return nil, err
	}

	repositories, err := NewRepositories(db, renderer)
	if err != nil {
		return nil, err
	}

	ingestionConfigs, err := NewIngestionConfigs(db, renderer)
	if err != nil {
		return nil, err
	}

	versions, err := NewVersions(db, renderer)
	if err != nil {
		return nil, err
	}

	jobs, err := NewJobs(db, renderer)
	if err != nil {
		return nil, err
	}

	documents, err := NewDocuments(db, renderer)
	if err != nil {
		return nil, err
	}

	chunks, err := NewChunks(db, renderer)
	if err != nil {
		return nil, err
	}

	symbols, err := NewSymbols(db, renderer)
	if err != nil {
		return nil, err
	}

	scipSymbols, err := NewSCIPSymbols(db, renderer)
	if err != nil {
		return nil, err
	}

	scipOccurrences, err := NewSCIPOccurrences(db, renderer)
	if err != nil {
		return nil, err
	}

	scipRelationships, err := NewSCIPRelationships(db, renderer)
	if err != nil {
		return nil, err
	}

	sessions := NewSessions(db)
	blobs := NewBlobs(bucket)

	keys, err := NewKeys(db, renderer)
	if err != nil {
		return nil, err
	}

	return &Stores{
		Users:             users,
		Repositories:      repositories,
		IngestionConfigs:  ingestionConfigs,
		Versions:          versions,
		Jobs:              jobs,
		Documents:         documents,
		Chunks:            chunks,
		Symbols:           symbols,
		SCIPSymbols:       scipSymbols,
		SCIPOccurrences:   scipOccurrences,
		SCIPRelationships: scipRelationships,
		Sessions:          sessions,
		Blobs:             blobs,
		Keys:              keys,
	}, nil
}
