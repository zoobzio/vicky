//go:build testing

package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zoobzio/rocco"
	rtesting "github.com/zoobzio/rocco/testing"
	"github.com/zoobzio/sum"
	sumtest "github.com/zoobzio/sum/testing"
	"github.com/zoobzio/vicky/api/contracts"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/api/wire"
)

// NewJob creates a test Job with sensible defaults.
func NewJob(t *testing.T) *models.Job {
	t.Helper()
	return &models.Job{
		ID:           1,
		VersionID:    10,
		RepositoryID: 100,
		UserID:       1000,
		Owner:        "testorg",
		RepoName:     "testrepo",
		Tag:          "v1.0.0",
		Stage:        models.JobStageFetch,
		Status:       models.JobStatusRunning,
	}
}

// NewChunk creates a test Chunk with the given ID and content.
func NewChunk(t *testing.T, id int64, content string) *models.Chunk {
	t.Helper()
	return &models.Chunk{
		ID:         id,
		DocumentID: 1,
		UserID:     1000,
		Owner:      "testorg",
		RepoName:   "testrepo",
		Tag:        "v1.0.0",
		Path:       "main.go",
		Kind:       models.ChunkKindFunction,
		StartLine:  1,
		EndLine:    10,
		Content:    content,
	}
}

// NewChunks creates n test Chunks with sequential IDs and generated content.
func NewChunks(t *testing.T, n int) []*models.Chunk {
	t.Helper()
	chunks := make([]*models.Chunk, n)
	for i := range chunks {
		chunks[i] = NewChunk(t, int64(i+1), fmt.Sprintf("chunk content %d", i+1))
	}
	return chunks
}

// NewVersion creates a test Version with sensible defaults.
func NewVersion(t *testing.T) *models.Version {
	t.Helper()
	return &models.Version{
		ID:           10,
		RepositoryID: 100,
		UserID:       1000,
		Owner:        "testorg",
		RepoName:     "testrepo",
		Tag:          "v1.0.0",
		CommitSHA:    "abc123",
		Status:       models.VersionStatusIngesting,
	}
}

// NewUser creates a test User with sensible defaults.
func NewUser(t *testing.T) *models.User {
	t.Helper()
	return &models.User{
		ID:          1000,
		Login:       "testuser",
		Email:       "test@example.com",
		AccessToken: "test-token",
	}
}

// NewIngestionConfig creates a test IngestionConfig with sensible defaults.
func NewIngestionConfig(t *testing.T) *models.IngestionConfig {
	t.Helper()
	return &models.IngestionConfig{
		ID:           1,
		RepositoryID: 100,
		UserID:       1000,
		Language:     models.LanguageGo,
		IncludeDocs:  true,
		MaxFileSize:  models.DefaultMaxFileSize,
	}
}

// NewDocument creates a test Document with the given ID and path.
func NewDocument(t *testing.T, id int64, path string) *models.Document {
	t.Helper()
	return &models.Document{
		ID:          id,
		VersionID:   10,
		UserID:      1000,
		Owner:       "testorg",
		RepoName:    "testrepo",
		Tag:         "v1.0.0",
		Path:        path,
		ContentType: models.ContentTypeCode,
		ContentHash: "abc123",
	}
}

// RegistryOption configures mock registrations for a test registry.
type RegistryOption func(k sum.Key)

// WithVersions registers a Versions implementation.
func WithVersions(v contracts.Versions) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Versions](k, v)
	}
}

// WithEmbedder registers an Embedder implementation.
func WithEmbedder(e contracts.Embedder) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Embedder](k, e)
	}
}

// WithChunks registers a Chunks implementation.
func WithChunks(c contracts.Chunks) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Chunks](k, c)
	}
}

// WithUsers registers a Users implementation.
func WithUsers(u contracts.Users) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Users](k, u)
	}
}

// WithGitHub registers a GitHub implementation.
func WithGitHub(g contracts.GitHub) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.GitHub](k, g)
	}
}

// WithIngestionConfigs registers an IngestionConfigs implementation.
func WithIngestionConfigs(c contracts.IngestionConfigs) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.IngestionConfigs](k, c)
	}
}

// WithBlobs registers a Blobs implementation.
func WithBlobs(b contracts.Blobs) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Blobs](k, b)
	}
}

// WithDocuments registers a Documents implementation.
func WithDocuments(d contracts.Documents) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Documents](k, d)
	}
}

// WithIndexer registers an Indexer implementation.
func WithIndexer(i contracts.Indexer) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Indexer](k, i)
	}
}

// WithSCIPSymbols registers a SCIPSymbols implementation.
func WithSCIPSymbols(s contracts.SCIPSymbols) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.SCIPSymbols](k, s)
	}
}

// WithSCIPOccurrences registers a SCIPOccurrences implementation.
func WithSCIPOccurrences(o contracts.SCIPOccurrences) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.SCIPOccurrences](k, o)
	}
}

// WithSCIPRelationships registers a SCIPRelationships implementation.
func WithSCIPRelationships(r contracts.SCIPRelationships) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.SCIPRelationships](k, r)
	}
}

// WithChunker registers a Chunker implementation.
func WithChunker(c contracts.Chunker) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Chunker](k, c)
	}
}

// WithJobs registers a Jobs implementation.
func WithJobs(j contracts.Jobs) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Jobs](k, j)
	}
}

// WithRepositories registers a Repositories implementation.
func WithRepositories(r contracts.Repositories) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Repositories](k, r)
	}
}

// WithSymbols registers a Symbols implementation.
func WithSymbols(s contracts.Symbols) RegistryOption {
	return func(k sum.Key) {
		sum.Register[contracts.Symbols](k, s)
	}
}

// NewKey creates a test Key with sensible defaults.
// The KeyHash and KeyPrefix are set to plausible test values.
func NewKey(t *testing.T) *models.Key {
	t.Helper()
	return &models.Key{
		ID:        1,
		UserID:    1000,
		Name:      "Test Key",
		KeyHash:   "dGVzdGhhc2g=",
		KeyPrefix: "vky_test",
		Scopes:    []string{"search", "intel"},
	}
}

// WithKeys registers a Keys implementation.
func WithKeys(k contracts.Keys) RegistryOption {
	return func(sk sum.Key) {
		sum.Register[contracts.Keys](sk, k)
	}
}

// SetupRegistry resets the global sum registry, registers mocks via options,
// freezes the registry, and returns a test context.
func SetupRegistry(t *testing.T, opts ...RegistryOption) context.Context {
	t.Helper()
	sum.Reset()
	k := sum.Start()
	for _, opt := range opts {
		opt(k)
	}
	sum.Freeze(k)
	t.Cleanup(sum.Reset)
	return sumtest.TestContext(t)
}

// NewRepository creates a test Repository with sensible defaults.
func NewRepository(t *testing.T) *models.Repository {
	t.Helper()
	return &models.Repository{
		ID:            100,
		GitHubID:      999,
		UserID:        1000,
		Owner:         "testorg",
		Name:          "testrepo",
		FullName:      "testorg/testrepo",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/testorg/testrepo",
	}
}

// withBoundaries registers wire boundaries (needed by response OnSend hooks).
func withBoundaries() RegistryOption {
	return func(k sum.Key) {
		sum.New()
		_ = wire.RegisterBoundaries(k)
	}
}

// SetupHandlerTest sets up the sum registry, creates a rocco engine with mock
// authentication, and returns the engine. The mock identity uses ID "1000".
func SetupHandlerTest(t *testing.T, opts ...RegistryOption) *rocco.Engine {
	t.Helper()
	// Prepend boundary registration so wire types work
	allOpts := append([]RegistryOption{withBoundaries()}, opts...)
	_ = SetupRegistry(t, allOpts...)

	identity := rtesting.NewMockIdentity("1000")
	engine := rtesting.TestEngineWithAuth(func(_ context.Context, _ *http.Request) (rocco.Identity, error) {
		return identity, nil
	})
	return engine
}
