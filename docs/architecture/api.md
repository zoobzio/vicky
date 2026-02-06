# Vicky API Architecture

Vicky is an embedding and vector search API for the zoobzio package ecosystem. It ingests public GitHub repositories, generates vector embeddings, and provides semantic search capabilities.

## Dependencies

| Package | Purpose |
|---------|---------|
| rocco | REST API framework |
| pipz | Ingest pipeline orchestration |
| grub | Data access (pgvector for vectors, postgres for metadata) |
| capitan | Domain event system |
| aperture | OTEL observability |
| zyn | Embedding generation (stubbed) |

## Package Structure

```
vicky/
├── api.go                 # Public interface, engine setup
├── handlers.go            # HTTP handler definitions
├── models.go              # Domain types
├── signals.go             # Capitan signals
├── ingest/
│   ├── pipeline.go        # Pipz ingest pipeline
│   ├── github.go          # GitHub fetcher stage
│   ├── parser.go          # Code parser stage
│   ├── chunker.go         # Content chunking stage
│   └── embedder.go        # Embedding generation stage (stubbed)
├── search/
│   ├── search.go          # Search logic
│   └── ranking.go         # Result ranking/filtering
├── store/
│   ├── vectors.go         # Grub Index[ChunkMeta] wrapper
│   └── packages.go        # Grub Database[Package] wrapper
├── internal/
│   └── github/            # GitHub API client
├── testing/
│   ├── helpers.go
│   └── integration/
├── Makefile
└── docker-compose.yml
```

## Domain Models

```go
// Package represents an ingested repository version
type Package struct {
    ID        string    `json:"id"`         // "{owner}/{repo}"
    Owner     string    `json:"owner"`
    Repo      string    `json:"repo"`
    Version   string    `json:"version"`    // semver tag
    Status    Status    `json:"status"`     // pending|ingesting|ready|failed
    ChunkCount int      `json:"chunk_count"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Status string

const (
    StatusPending   Status = "pending"
    StatusIngesting Status = "ingesting"
    StatusReady     Status = "ready"
    StatusFailed    Status = "failed"
)

// Chunk represents an embedded content segment
type Chunk struct {
    ID        string      `json:"id"`         // deterministic hash
    PackageID string      `json:"package_id"` // "{owner}/{repo}@{version}"
    Tier      Tier        `json:"tier"`       // package|file|symbol
    FilePath  string      `json:"file_path"`
    StartLine int         `json:"start_line"`
    EndLine   int         `json:"end_line"`
    Symbol    string      `json:"symbol"`     // function/type name if applicable
    Content   string      `json:"content"`    // raw text content
    Vector    []float32   `json:"-"`          // embedding (not serialised)
}

type Tier string

const (
    TierPackage Tier = "package"  // README, package doc
    TierFile    Tier = "file"     // file-level content
    TierSymbol  Tier = "symbol"   // function, type, method
)

// ChunkMeta is the grub Index metadata (stored alongside vector)
type ChunkMeta struct {
    PackageID string `json:"package_id"`
    Tier      Tier   `json:"tier"`
    FilePath  string `json:"file_path"`
    StartLine int    `json:"start_line"`
    EndLine   int    `json:"end_line"`
    Symbol    string `json:"symbol"`
}

// SearchResult returned to clients
type SearchResult struct {
    Score    float32 `json:"score"`
    Tier     Tier    `json:"tier"`
    FilePath string  `json:"file_path"`
    Symbol   string  `json:"symbol,omitempty"`
    Lines    string  `json:"lines,omitempty"`   // "L10-L25"
    URL      string  `json:"url"`               // GitHub permalink
    Preview  string  `json:"preview"`           // content snippet
}
```

## API Endpoints

### Ingest

```
POST /ingest/{owner}/{repo}/{version}
```

Initiates async ingestion of a GitHub repository at a specific version tag.

**Path Parameters:**
- `owner` - GitHub user/org
- `repo` - Repository name
- `version` - Semver tag (e.g., `v1.2.3`)

**Response:** `202 Accepted`
```json
{
  "package_id": "zoobzio/pipz@v1.0.0",
  "status": "pending"
}
```

**Errors:**
- `400` - Invalid version format
- `404` - Repository or tag not found
- `409` - Version already ingested

### Ingest Status

```
GET /ingest/{owner}/{repo}/{version}
```

Returns current ingestion status.

**Response:** `200 OK`
```json
{
  "package_id": "zoobzio/pipz@v1.0.0",
  "status": "ready",
  "chunk_count": 247,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:32:15Z"
}
```

### Search

```
GET /search/{owner}/{repo}/{version}
```

Semantic search within an ingested package.

**Path Parameters:**
- `owner` - GitHub user/org
- `repo` - Repository name
- `version` - Semver tag

**Query Parameters:**
- `q` (required) - Search query
- `scope` (optional) - `discovery` | `navigation` | `all` (default: `all`)
- `limit` (optional) - Max results (default: 10, max: 50)

**Response:** `200 OK`
```json
{
  "query": "circuit breaker pattern",
  "results": [
    {
      "score": 0.89,
      "tier": "symbol",
      "file_path": "circuit.go",
      "symbol": "NewCircuitBreaker",
      "lines": "L45-L78",
      "url": "https://github.com/zoobzio/pipz/blob/v1.0.0/circuit.go#L45-L78",
      "preview": "// NewCircuitBreaker creates a circuit breaker that opens after..."
    }
  ]
}
```

**Errors:**
- `400` - Missing query parameter
- `404` - Package not found or not ready

### List Packages

```
GET /packages
```

Lists all ingested packages.

**Query Parameters:**
- `owner` (optional) - Filter by owner
- `status` (optional) - Filter by status

**Response:** `200 OK`
```json
{
  "packages": [
    {
      "id": "zoobzio/pipz",
      "versions": ["v1.0.0", "v1.1.0", "v1.2.0"],
      "latest": "v1.2.0"
    }
  ]
}
```

## Ingest Pipeline

Built with pipz, the ingest pipeline processes repositories through discrete stages:

```
┌─────────────────────────────────────────────────────────────────┐
│                     Ingest Pipeline                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────────┐  │
│  │  Fetch   │──▶│  Parse   │──▶│  Chunk   │──▶│    Embed     │  │
│  │ (GitHub) │   │  (AST)   │   │ (Tiers)  │   │   (stubbed)  │  │
│  └──────────┘   └──────────┘   └──────────┘   └──────────────┘  │
│                                                      │           │
│                                               ┌──────▼──────┐    │
│                                               │    Store    │    │
│                                               │  (pgvector) │    │
│                                               └─────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Stages:**

1. **Fetch** - Clone repository at tag via GitHub API/git
2. **Parse** - Extract Go AST for symbols, markdown for docs
3. **Chunk** - Segment content into tiered chunks with metadata
4. **Embed** - Generate vector embeddings (stubbed interface)
5. **Store** - Persist vectors + metadata to pgvector via grub

**Resilience (via pipz):**
- Retry with backoff on GitHub rate limits
- Circuit breaker for embedding service
- Timeout per stage

## Signals (Capitan)

```go
var (
    // Ingest lifecycle
    SignalIngestStarted   = capitan.NewSignal("vicky.ingest.started", "Package ingestion initiated")
    SignalIngestCompleted = capitan.NewSignal("vicky.ingest.completed", "Package ingestion finished")
    SignalIngestFailed    = capitan.NewSignal("vicky.ingest.failed", "Package ingestion failed")

    // Pipeline stages
    SignalFetchStarted    = capitan.NewSignal("vicky.fetch.started", "GitHub fetch initiated")
    SignalFetchCompleted  = capitan.NewSignal("vicky.fetch.completed", "GitHub fetch finished")
    SignalParseCompleted  = capitan.NewSignal("vicky.parse.completed", "AST parsing finished")
    SignalChunkCompleted  = capitan.NewSignal("vicky.chunk.completed", "Content chunking finished")
    SignalEmbedCompleted  = capitan.NewSignal("vicky.embed.completed", "Embedding generation finished")

    // Search
    SignalSearchExecuted  = capitan.NewSignal("vicky.search.executed", "Search query processed")
)

// Keys
var (
    KeyPackageID  = capitan.NewStringKey("package_id")
    KeyVersion    = capitan.NewStringKey("version")
    KeyChunkCount = capitan.NewIntKey("chunk_count")
    KeyQuery      = capitan.NewStringKey("query")
    KeyResultCount = capitan.NewIntKey("result_count")
    KeyDuration   = capitan.NewDurationKey("duration")
    KeyError      = capitan.NewErrorKey("error")
)
```

## Database Schema

**packages table (grub Database):**
```sql
CREATE TABLE packages (
    id         TEXT PRIMARY KEY,  -- "owner/repo"
    owner      TEXT NOT NULL,
    repo       TEXT NOT NULL,
    version    TEXT NOT NULL,
    status     TEXT NOT NULL,
    chunk_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE(owner, repo, version)
);
```

**chunks table (grub Index with pgvector):**
```sql
CREATE TABLE chunks (
    id         TEXT PRIMARY KEY,
    vector     vector(1536),       -- dimension depends on embedding model
    metadata   JSONB NOT NULL,     -- ChunkMeta
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX chunks_vector_idx ON chunks
    USING ivfflat (vector vector_cosine_ops) WITH (lists = 100);
```

## Configuration

```yaml
# config.yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  name: vicky
  user: vicky
  password: ${VICKY_DB_PASSWORD}

github:
  token: ${GITHUB_TOKEN}  # optional, for higher rate limits

embedding:
  provider: stub  # openai|google|voyage|stub
  model: text-embedding-3-small
  dimensions: 1536

observability:
  otlp_endpoint: localhost:4317
```

## Docker Compose (Local Development)

```yaml
services:
  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_DB: vicky
      POSTGRES_USER: vicky
      POSTGRES_PASSWORD: vicky
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # OTLP gRPC

  vicky:
    build: .
    depends_on:
      - postgres
      - jaeger
    environment:
      VICKY_DB_PASSWORD: vicky
      OTEL_EXPORTER_OTLP_ENDPOINT: jaeger:4317
    ports:
      - "8080:8080"

volumes:
  pgdata:
```

## Embedding Interface (Stubbed)

```go
// Embedder generates vector embeddings for text content.
// Implementation is stubbed pending zyn extension.
type Embedder interface {
    // Embed generates embeddings for the given texts.
    Embed(ctx context.Context, texts []string) ([][]float32, error)

    // Dimensions returns the vector dimensionality.
    Dimensions() int
}

// StubEmbedder returns zero vectors for development/testing.
type StubEmbedder struct {
    dims int
}

func (s *StubEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
    result := make([][]float32, len(texts))
    for i := range texts {
        result[i] = make([]float32, s.dims)
    }
    return result, nil
}

func (s *StubEmbedder) Dimensions() int {
    return s.dims
}
```

## Open Questions

1. **Version lifecycle** - How long to retain old versions? Prune policy?
2. **Cross-package search** - Search across all packages vs single package?
3. **Incremental updates** - Re-ingest entire package or diff-based updates?
4. **Authentication** - Public API or require API keys?
5. **Rate limiting** - Per-IP, per-key, or none?
