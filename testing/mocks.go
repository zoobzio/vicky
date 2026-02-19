//go:build testing

package testing

import (
	"context"

	"github.com/zoobzio/grub"
	"github.com/zoobzio/vicky/external/chunker"
	"github.com/zoobzio/vicky/external/github"
	"github.com/zoobzio/vicky/external/indexer"
	"github.com/zoobzio/vicky/models"
)

// MockVersions implements contracts.Versions with function-field overrides.
// Unset fields return sensible defaults.
type MockVersions struct {
	OnGet                 func(ctx context.Context, key string) (*models.Version, error)
	OnSet                 func(ctx context.Context, key string, version *models.Version) error
	OnListByUserAndRepo   func(ctx context.Context, userID int64, owner, repoName string) ([]*models.Version, error)
	OnGetByUserRepoAndTag func(ctx context.Context, userID int64, owner, repoName, tag string) (*models.Version, error)
	OnUpdateStatus        func(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error)
}

func (m *MockVersions) Get(ctx context.Context, key string) (*models.Version, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Version{}, nil
}

func (m *MockVersions) Set(ctx context.Context, key string, version *models.Version) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, version)
	}
	return nil
}

func (m *MockVersions) ListByUserAndRepo(ctx context.Context, userID int64, owner, repoName string) ([]*models.Version, error) {
	if m.OnListByUserAndRepo != nil {
		return m.OnListByUserAndRepo(ctx, userID, owner, repoName)
	}
	return nil, nil
}

func (m *MockVersions) GetByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) (*models.Version, error) {
	if m.OnGetByUserRepoAndTag != nil {
		return m.OnGetByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return &models.Version{}, nil
}

func (m *MockVersions) UpdateStatus(ctx context.Context, id int64, status models.VersionStatus, versionErr *string) (*models.Version, error) {
	if m.OnUpdateStatus != nil {
		return m.OnUpdateStatus(ctx, id, status, versionErr)
	}
	return &models.Version{ID: id, Status: status}, nil
}

// MockEmbedder implements contracts.Embedder with function-field overrides.
type MockEmbedder struct {
	OnEmbed      func(ctx context.Context, texts []string) ([][]float32, error)
	OnEmbedQuery func(ctx context.Context, texts []string) ([][]float32, error)
	OnDimensions func() int
}

func (m *MockEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if m.OnEmbed != nil {
		return m.OnEmbed(ctx, texts)
	}
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = make([]float32, 3)
	}
	return result, nil
}

func (m *MockEmbedder) EmbedQuery(ctx context.Context, texts []string) ([][]float32, error) {
	if m.OnEmbedQuery != nil {
		return m.OnEmbedQuery(ctx, texts)
	}
	result := make([][]float32, len(texts))
	for i := range result {
		result[i] = make([]float32, 3)
	}
	return result, nil
}

func (m *MockEmbedder) Dimensions() int {
	if m.OnDimensions != nil {
		return m.OnDimensions()
	}
	return 3
}

// MockChunks implements contracts.Chunks with function-field overrides.
type MockChunks struct {
	OnGet                      func(ctx context.Context, key string) (*models.Chunk, error)
	OnSet                      func(ctx context.Context, key string, chunk *models.Chunk) error
	OnListByUserRepoAndTag     func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error)
	OnListByUserRepoTagAndPath func(ctx context.Context, userID int64, owner, repoName, tag, path string) ([]*models.Chunk, error)
	OnSearch                   func(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Chunk, error)
	OnSearchByKind             func(ctx context.Context, userID int64, owner, repoName, tag string, kind models.ChunkKind, vector []float32, limit int) ([]*models.Chunk, error)
}

func (m *MockChunks) Get(ctx context.Context, key string) (*models.Chunk, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Chunk{}, nil
}

func (m *MockChunks) Set(ctx context.Context, key string, chunk *models.Chunk) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, chunk)
	}
	return nil
}

func (m *MockChunks) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Chunk, error) {
	if m.OnListByUserRepoAndTag != nil {
		return m.OnListByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockChunks) ListByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) ([]*models.Chunk, error) {
	if m.OnListByUserRepoTagAndPath != nil {
		return m.OnListByUserRepoTagAndPath(ctx, userID, owner, repoName, tag, path)
	}
	return nil, nil
}

func (m *MockChunks) Search(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Chunk, error) {
	if m.OnSearch != nil {
		return m.OnSearch(ctx, userID, owner, repoName, tag, vector, limit)
	}
	return nil, nil
}

func (m *MockChunks) SearchByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.ChunkKind, vector []float32, limit int) ([]*models.Chunk, error) {
	if m.OnSearchByKind != nil {
		return m.OnSearchByKind(ctx, userID, owner, repoName, tag, kind, vector, limit)
	}
	return nil, nil
}

// MockUsers implements contracts.Users with function-field overrides.
type MockUsers struct {
	OnGet        func(ctx context.Context, key string) (*models.User, error)
	OnSet        func(ctx context.Context, key string, user *models.User) error
	OnDelete     func(ctx context.Context, key string) error
	OnGetByLogin func(ctx context.Context, login string) (*models.User, error)
	OnList       func(ctx context.Context, filter interface{}, limit, offset int) ([]*models.User, error)
	OnCount      func(ctx context.Context, filter interface{}) (int, error)
}

func (m *MockUsers) Get(ctx context.Context, key string) (*models.User, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.User{ID: 1000, Login: "testuser", AccessToken: "test-token"}, nil
}

func (m *MockUsers) Set(ctx context.Context, key string, user *models.User) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, user)
	}
	return nil
}

func (m *MockUsers) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

func (m *MockUsers) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	if m.OnGetByLogin != nil {
		return m.OnGetByLogin(ctx, login)
	}
	return &models.User{ID: 1000, Login: login, AccessToken: "test-token"}, nil
}

func (m *MockUsers) List(ctx context.Context, filter interface{}, limit, offset int) ([]*models.User, error) {
	if m.OnList != nil {
		return m.OnList(ctx, filter, limit, offset)
	}
	return nil, nil
}

func (m *MockUsers) Count(ctx context.Context, filter interface{}) (int, error) {
	if m.OnCount != nil {
		return m.OnCount(ctx, filter)
	}
	return 0, nil
}

// MockGitHub implements contracts.GitHub with function-field overrides.
type MockGitHub struct {
	OnGetTree             func(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error)
	OnGetFileContent      func(ctx context.Context, token, owner, repo, path, ref string) (*github.FileContent, error)
	OnGetFileContentBatch func(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*github.FileContent, error)
}

func (m *MockGitHub) GetTree(ctx context.Context, token, owner, repo, ref string) ([]github.TreeEntry, error) {
	if m.OnGetTree != nil {
		return m.OnGetTree(ctx, token, owner, repo, ref)
	}
	return nil, nil
}

func (m *MockGitHub) GetFileContent(ctx context.Context, token, owner, repo, path, ref string) (*github.FileContent, error) {
	if m.OnGetFileContent != nil {
		return m.OnGetFileContent(ctx, token, owner, repo, path, ref)
	}
	return &github.FileContent{Path: path, Content: []byte("content")}, nil
}

func (m *MockGitHub) GetFileContentBatch(ctx context.Context, token, owner, repo, ref string, paths []string) ([]*github.FileContent, error) {
	if m.OnGetFileContentBatch != nil {
		return m.OnGetFileContentBatch(ctx, token, owner, repo, ref, paths)
	}
	result := make([]*github.FileContent, len(paths))
	for i, p := range paths {
		result[i] = &github.FileContent{Path: p, Content: []byte("content")}
	}
	return result, nil
}

// MockIngestionConfigs implements contracts.IngestionConfigs with function-field overrides.
type MockIngestionConfigs struct {
	OnGet                func(ctx context.Context, key string) (*models.IngestionConfig, error)
	OnSet                func(ctx context.Context, key string, config *models.IngestionConfig) error
	OnGetByRepositoryID  func(ctx context.Context, repositoryID int64) (*models.IngestionConfig, error)
}

func (m *MockIngestionConfigs) Get(ctx context.Context, key string) (*models.IngestionConfig, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.IngestionConfig{
		ID:          1,
		Language:    models.LanguageGo,
		IncludeDocs: true,
		MaxFileSize: models.DefaultMaxFileSize,
	}, nil
}

func (m *MockIngestionConfigs) Set(ctx context.Context, key string, config *models.IngestionConfig) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, config)
	}
	return nil
}

func (m *MockIngestionConfigs) GetByRepositoryID(ctx context.Context, repositoryID int64) (*models.IngestionConfig, error) {
	if m.OnGetByRepositoryID != nil {
		return m.OnGetByRepositoryID(ctx, repositoryID)
	}
	return &models.IngestionConfig{
		ID:           1,
		RepositoryID: repositoryID,
		Language:     models.LanguageGo,
		IncludeDocs:  true,
		MaxFileSize:  models.DefaultMaxFileSize,
	}, nil
}

// MockBlobs implements contracts.Blobs with function-field overrides.
type MockBlobs struct {
	OnGetByPath      func(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error)
	OnPutBlob        func(ctx context.Context, userID int64, blob *models.Blob) error
	OnDeleteByPath   func(ctx context.Context, userID int64, owner, repo, tag, path string) error
	OnListByVersion  func(ctx context.Context, userID int64, owner, repo, tag string, limit int) ([]grub.ObjectInfo, error)
	OnListByRepo     func(ctx context.Context, userID int64, owner, repo string, limit int) ([]grub.ObjectInfo, error)
}

func (m *MockBlobs) GetByPath(ctx context.Context, userID int64, owner, repo, tag, path string) (*grub.Object[models.Blob], error) {
	if m.OnGetByPath != nil {
		return m.OnGetByPath(ctx, userID, owner, repo, tag, path)
	}
	return &grub.Object[models.Blob]{
		Key:  path,
		Data: models.Blob{Path: path, Content: "content", Owner: owner, Repo: repo, Tag: tag},
	}, nil
}

func (m *MockBlobs) PutBlob(ctx context.Context, userID int64, blob *models.Blob) error {
	if m.OnPutBlob != nil {
		return m.OnPutBlob(ctx, userID, blob)
	}
	return nil
}

func (m *MockBlobs) DeleteByPath(ctx context.Context, userID int64, owner, repo, tag, path string) error {
	if m.OnDeleteByPath != nil {
		return m.OnDeleteByPath(ctx, userID, owner, repo, tag, path)
	}
	return nil
}

func (m *MockBlobs) ListByVersion(ctx context.Context, userID int64, owner, repo, tag string, limit int) ([]grub.ObjectInfo, error) {
	if m.OnListByVersion != nil {
		return m.OnListByVersion(ctx, userID, owner, repo, tag, limit)
	}
	return nil, nil
}

func (m *MockBlobs) ListByRepo(ctx context.Context, userID int64, owner, repo string, limit int) ([]grub.ObjectInfo, error) {
	if m.OnListByRepo != nil {
		return m.OnListByRepo(ctx, userID, owner, repo, limit)
	}
	return nil, nil
}

// MockDocuments implements contracts.Documents with function-field overrides.
type MockDocuments struct {
	OnGet                      func(ctx context.Context, key string) (*models.Document, error)
	OnSet                      func(ctx context.Context, key string, doc *models.Document) error
	OnListByUserRepoAndTag     func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error)
	OnGetByUserRepoTagAndPath  func(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error)
	OnFindSimilar              func(ctx context.Context, userID int64, vector []float32, limit int) ([]*models.Document, error)
	OnFindSimilarInVersion     func(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Document, error)
}

func (m *MockDocuments) Get(ctx context.Context, key string) (*models.Document, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Document{}, nil
}

func (m *MockDocuments) Set(ctx context.Context, key string, doc *models.Document) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, doc)
	}
	return nil
}

func (m *MockDocuments) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Document, error) {
	if m.OnListByUserRepoAndTag != nil {
		return m.OnListByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockDocuments) GetByUserRepoTagAndPath(ctx context.Context, userID int64, owner, repoName, tag, path string) (*models.Document, error) {
	if m.OnGetByUserRepoTagAndPath != nil {
		return m.OnGetByUserRepoTagAndPath(ctx, userID, owner, repoName, tag, path)
	}
	return &models.Document{}, nil
}

func (m *MockDocuments) FindSimilar(ctx context.Context, userID int64, vector []float32, limit int) ([]*models.Document, error) {
	if m.OnFindSimilar != nil {
		return m.OnFindSimilar(ctx, userID, vector, limit)
	}
	return nil, nil
}

func (m *MockDocuments) FindSimilarInVersion(ctx context.Context, userID int64, owner, repoName, tag string, vector []float32, limit int) ([]*models.Document, error) {
	if m.OnFindSimilarInVersion != nil {
		return m.OnFindSimilarInVersion(ctx, userID, owner, repoName, tag, vector, limit)
	}
	return nil, nil
}

// MockIndexer implements contracts.Indexer with function-field overrides.
type MockIndexer struct {
	OnIndex    func(ctx context.Context, req indexer.Request) (*indexer.Result, error)
	OnSupports func(language models.Language) bool
}

func (m *MockIndexer) Index(ctx context.Context, req indexer.Request) (*indexer.Result, error) {
	if m.OnIndex != nil {
		return m.OnIndex(ctx, req)
	}
	return &indexer.Result{JobID: req.JobID, VersionID: req.VersionID}, nil
}

func (m *MockIndexer) Supports(language models.Language) bool {
	if m.OnSupports != nil {
		return m.OnSupports(language)
	}
	return true
}

// MockSCIPSymbols implements contracts.SCIPSymbols with function-field overrides.
type MockSCIPSymbols struct {
	OnGet                    func(ctx context.Context, key string) (*models.SCIPSymbol, error)
	OnSet                    func(ctx context.Context, key string, symbol *models.SCIPSymbol) error
	OnListByUserRepoAndTag   func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPSymbol, error)
	OnListByDocument         func(ctx context.Context, documentID int64) ([]*models.SCIPSymbol, error)
	OnGetBySymbol            func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error)
	OnListByKind             func(ctx context.Context, userID int64, owner, repoName, tag string, kind models.SCIPSymbolKind) ([]*models.SCIPSymbol, error)
	OnListByEnclosingSymbol  func(ctx context.Context, userID int64, owner, repoName, tag, enclosingSymbol string) ([]*models.SCIPSymbol, error)
}

func (m *MockSCIPSymbols) Get(ctx context.Context, key string) (*models.SCIPSymbol, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.SCIPSymbol{}, nil
}

func (m *MockSCIPSymbols) Set(ctx context.Context, key string, symbol *models.SCIPSymbol) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, symbol)
	}
	return nil
}

func (m *MockSCIPSymbols) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPSymbol, error) {
	if m.OnListByUserRepoAndTag != nil {
		return m.OnListByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockSCIPSymbols) ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPSymbol, error) {
	if m.OnListByDocument != nil {
		return m.OnListByDocument(ctx, documentID)
	}
	return nil, nil
}

func (m *MockSCIPSymbols) GetBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) (*models.SCIPSymbol, error) {
	if m.OnGetBySymbol != nil {
		return m.OnGetBySymbol(ctx, userID, owner, repoName, tag, symbol)
	}
	return &models.SCIPSymbol{}, nil
}

func (m *MockSCIPSymbols) ListByKind(ctx context.Context, userID int64, owner, repoName, tag string, kind models.SCIPSymbolKind) ([]*models.SCIPSymbol, error) {
	if m.OnListByKind != nil {
		return m.OnListByKind(ctx, userID, owner, repoName, tag, kind)
	}
	return nil, nil
}

func (m *MockSCIPSymbols) ListByEnclosingSymbol(ctx context.Context, userID int64, owner, repoName, tag, enclosingSymbol string) ([]*models.SCIPSymbol, error) {
	if m.OnListByEnclosingSymbol != nil {
		return m.OnListByEnclosingSymbol(ctx, userID, owner, repoName, tag, enclosingSymbol)
	}
	return nil, nil
}

// MockSCIPOccurrences implements contracts.SCIPOccurrences with function-field overrides.
type MockSCIPOccurrences struct {
	OnGet                  func(ctx context.Context, key string) (*models.SCIPOccurrence, error)
	OnSet                  func(ctx context.Context, key string, occurrence *models.SCIPOccurrence) error
	OnListByUserRepoAndTag func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPOccurrence, error)
	OnListByDocument       func(ctx context.Context, documentID int64) ([]*models.SCIPOccurrence, error)
	OnListBySymbol         func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
	OnListDefinitions      func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
	OnListReferences       func(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error)
}

func (m *MockSCIPOccurrences) Get(ctx context.Context, key string) (*models.SCIPOccurrence, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.SCIPOccurrence{}, nil
}

func (m *MockSCIPOccurrences) Set(ctx context.Context, key string, occurrence *models.SCIPOccurrence) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, occurrence)
	}
	return nil
}

func (m *MockSCIPOccurrences) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.SCIPOccurrence, error) {
	if m.OnListByUserRepoAndTag != nil {
		return m.OnListByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockSCIPOccurrences) ListByDocument(ctx context.Context, documentID int64) ([]*models.SCIPOccurrence, error) {
	if m.OnListByDocument != nil {
		return m.OnListByDocument(ctx, documentID)
	}
	return nil, nil
}

func (m *MockSCIPOccurrences) ListBySymbol(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	if m.OnListBySymbol != nil {
		return m.OnListBySymbol(ctx, userID, owner, repoName, tag, symbol)
	}
	return nil, nil
}

func (m *MockSCIPOccurrences) ListDefinitions(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	if m.OnListDefinitions != nil {
		return m.OnListDefinitions(ctx, userID, owner, repoName, tag, symbol)
	}
	return nil, nil
}

func (m *MockSCIPOccurrences) ListReferences(ctx context.Context, userID int64, owner, repoName, tag, symbol string) ([]*models.SCIPOccurrence, error) {
	if m.OnListReferences != nil {
		return m.OnListReferences(ctx, userID, owner, repoName, tag, symbol)
	}
	return nil, nil
}

// MockSCIPRelationships implements contracts.SCIPRelationships with function-field overrides.
type MockSCIPRelationships struct {
	OnGet                 func(ctx context.Context, key string) (*models.SCIPRelationship, error)
	OnSet                 func(ctx context.Context, key string, rel *models.SCIPRelationship) error
	OnListBySymbol        func(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error)
	OnListImplementations func(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error)
	OnListByTargetSymbol  func(ctx context.Context, targetSymbol string) ([]*models.SCIPRelationship, error)
}

func (m *MockSCIPRelationships) Get(ctx context.Context, key string) (*models.SCIPRelationship, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.SCIPRelationship{}, nil
}

func (m *MockSCIPRelationships) Set(ctx context.Context, key string, rel *models.SCIPRelationship) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, rel)
	}
	return nil
}

func (m *MockSCIPRelationships) ListBySymbol(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error) {
	if m.OnListBySymbol != nil {
		return m.OnListBySymbol(ctx, symbolID)
	}
	return nil, nil
}

func (m *MockSCIPRelationships) ListImplementations(ctx context.Context, symbolID int64) ([]*models.SCIPRelationship, error) {
	if m.OnListImplementations != nil {
		return m.OnListImplementations(ctx, symbolID)
	}
	return nil, nil
}

func (m *MockSCIPRelationships) ListByTargetSymbol(ctx context.Context, targetSymbol string) ([]*models.SCIPRelationship, error) {
	if m.OnListByTargetSymbol != nil {
		return m.OnListByTargetSymbol(ctx, targetSymbol)
	}
	return nil, nil
}

// MockChunker implements contracts.Chunker with function-field overrides.
type MockChunker struct {
	OnChunk    func(ctx context.Context, language string, filename string, content []byte) ([]chunker.Result, error)
	OnSupports func(language string) bool
}

func (m *MockChunker) Chunk(ctx context.Context, language string, filename string, content []byte) ([]chunker.Result, error) {
	if m.OnChunk != nil {
		return m.OnChunk(ctx, language, filename, content)
	}
	return nil, nil
}

func (m *MockChunker) Supports(language string) bool {
	if m.OnSupports != nil {
		return m.OnSupports(language)
	}
	return true
}

// MockJobs implements contracts.Jobs with function-field overrides.
type MockJobs struct {
	OnGet              func(ctx context.Context, key string) (*models.Job, error)
	OnSet              func(ctx context.Context, key string, job *models.Job) error
	OnListByVersionID  func(ctx context.Context, versionID int64) ([]*models.Job, error)
	OnLatestByVersionID func(ctx context.Context, versionID int64) (*models.Job, error)
	OnListByUser       func(ctx context.Context, userID int64) ([]*models.Job, error)
	OnListByStatus     func(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error)
	OnUpdateProgress   func(ctx context.Context, id int64, stage models.JobStage, progress int, itemsProcessed int) error
	OnStart            func(ctx context.Context, id int64) error
	OnMarkFailed       func(ctx context.Context, id int64, errMsg string) error
	OnMarkCompleted    func(ctx context.Context, id int64) error
	OnMarkCancelled    func(ctx context.Context, id int64) error
	OnIsCancelling     func(ctx context.Context, id int64) (bool, error)
}

func (m *MockJobs) Get(ctx context.Context, key string) (*models.Job, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Job{}, nil
}

func (m *MockJobs) Set(ctx context.Context, key string, job *models.Job) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, job)
	}
	return nil
}

func (m *MockJobs) ListByVersionID(ctx context.Context, versionID int64) ([]*models.Job, error) {
	if m.OnListByVersionID != nil {
		return m.OnListByVersionID(ctx, versionID)
	}
	return nil, nil
}

func (m *MockJobs) LatestByVersionID(ctx context.Context, versionID int64) (*models.Job, error) {
	if m.OnLatestByVersionID != nil {
		return m.OnLatestByVersionID(ctx, versionID)
	}
	return &models.Job{}, nil
}

func (m *MockJobs) ListByUser(ctx context.Context, userID int64) ([]*models.Job, error) {
	if m.OnListByUser != nil {
		return m.OnListByUser(ctx, userID)
	}
	return nil, nil
}

func (m *MockJobs) ListByStatus(ctx context.Context, userID int64, status models.JobStatus) ([]*models.Job, error) {
	if m.OnListByStatus != nil {
		return m.OnListByStatus(ctx, userID, status)
	}
	return nil, nil
}

func (m *MockJobs) UpdateProgress(ctx context.Context, id int64, stage models.JobStage, progress int, itemsProcessed int) error {
	if m.OnUpdateProgress != nil {
		return m.OnUpdateProgress(ctx, id, stage, progress, itemsProcessed)
	}
	return nil
}

func (m *MockJobs) Start(ctx context.Context, id int64) error {
	if m.OnStart != nil {
		return m.OnStart(ctx, id)
	}
	return nil
}

func (m *MockJobs) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	if m.OnMarkFailed != nil {
		return m.OnMarkFailed(ctx, id, errMsg)
	}
	return nil
}

func (m *MockJobs) MarkCompleted(ctx context.Context, id int64) error {
	if m.OnMarkCompleted != nil {
		return m.OnMarkCompleted(ctx, id)
	}
	return nil
}

func (m *MockJobs) MarkCancelled(ctx context.Context, id int64) error {
	if m.OnMarkCancelled != nil {
		return m.OnMarkCancelled(ctx, id)
	}
	return nil
}

func (m *MockJobs) IsCancelling(ctx context.Context, id int64) (bool, error) {
	if m.OnIsCancelling != nil {
		return m.OnIsCancelling(ctx, id)
	}
	return false, nil
}

// MockRepositories implements contracts.Repositories with function-field overrides.
type MockRepositories struct {
	OnGet                    func(ctx context.Context, key string) (*models.Repository, error)
	OnSet                    func(ctx context.Context, key string, repo *models.Repository) error
	OnListByUserID           func(ctx context.Context, userID int64) ([]*models.Repository, error)
	OnGetByUserAndGitHubID   func(ctx context.Context, userID, githubID int64) (*models.Repository, error)
	OnGetByUserOwnerAndName  func(ctx context.Context, userID int64, owner, name string) (*models.Repository, error)
}

func (m *MockRepositories) Get(ctx context.Context, key string) (*models.Repository, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Repository{}, nil
}

func (m *MockRepositories) Set(ctx context.Context, key string, repo *models.Repository) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, repo)
	}
	return nil
}

func (m *MockRepositories) ListByUserID(ctx context.Context, userID int64) ([]*models.Repository, error) {
	if m.OnListByUserID != nil {
		return m.OnListByUserID(ctx, userID)
	}
	return nil, nil
}

func (m *MockRepositories) GetByUserAndGitHubID(ctx context.Context, userID, githubID int64) (*models.Repository, error) {
	if m.OnGetByUserAndGitHubID != nil {
		return m.OnGetByUserAndGitHubID(ctx, userID, githubID)
	}
	return &models.Repository{}, nil
}

func (m *MockRepositories) GetByUserOwnerAndName(ctx context.Context, userID int64, owner, name string) (*models.Repository, error) {
	if m.OnGetByUserOwnerAndName != nil {
		return m.OnGetByUserOwnerAndName(ctx, userID, owner, name)
	}
	return &models.Repository{}, nil
}

// MockKeys implements contracts.Keys with function-field overrides.
type MockKeys struct {
	OnGet            func(ctx context.Context, key string) (*models.Key, error)
	OnSet            func(ctx context.Context, key string, apiKey *models.Key) error
	OnDelete         func(ctx context.Context, key string) error
	OnGetByKeyHash   func(ctx context.Context, hash string) (*models.Key, error)
	OnListByUserID   func(ctx context.Context, userID int64) ([]*models.Key, error)
	OnUpdateLastUsed func(ctx context.Context, id int64) error
}

func (m *MockKeys) Get(ctx context.Context, key string) (*models.Key, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Key{ID: 1, UserID: 1000, Name: "Test Key", KeyHash: "dGVzdGhhc2g=", KeyPrefix: "vky_test", Scopes: []string{"search"}}, nil
}

func (m *MockKeys) Set(ctx context.Context, key string, apiKey *models.Key) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, apiKey)
	}
	return nil
}

func (m *MockKeys) Delete(ctx context.Context, key string) error {
	if m.OnDelete != nil {
		return m.OnDelete(ctx, key)
	}
	return nil
}

func (m *MockKeys) GetByKeyHash(ctx context.Context, hash string) (*models.Key, error) {
	if m.OnGetByKeyHash != nil {
		return m.OnGetByKeyHash(ctx, hash)
	}
	return &models.Key{ID: 1, UserID: 1000, Name: "Test Key", KeyHash: hash, KeyPrefix: "vky_test", Scopes: []string{"search"}}, nil
}

func (m *MockKeys) ListByUserID(ctx context.Context, userID int64) ([]*models.Key, error) {
	if m.OnListByUserID != nil {
		return m.OnListByUserID(ctx, userID)
	}
	return nil, nil
}

func (m *MockKeys) UpdateLastUsed(ctx context.Context, id int64) error {
	if m.OnUpdateLastUsed != nil {
		return m.OnUpdateLastUsed(ctx, id)
	}
	return nil
}

// MockSymbols implements contracts.Symbols with function-field overrides.
type MockSymbols struct {
	OnGet                          func(ctx context.Context, key string) (*models.Symbol, error)
	OnSet                          func(ctx context.Context, key string, symbol *models.Symbol) error
	OnListByUserRepoAndTag         func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error)
	OnListExportedByUserRepoAndTag func(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error)
	OnFindRelated                  func(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error)
	OnFindRelatedExported          func(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error)
}

func (m *MockSymbols) Get(ctx context.Context, key string) (*models.Symbol, error) {
	if m.OnGet != nil {
		return m.OnGet(ctx, key)
	}
	return &models.Symbol{}, nil
}

func (m *MockSymbols) Set(ctx context.Context, key string, symbol *models.Symbol) error {
	if m.OnSet != nil {
		return m.OnSet(ctx, key, symbol)
	}
	return nil
}

func (m *MockSymbols) ListByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error) {
	if m.OnListByUserRepoAndTag != nil {
		return m.OnListByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockSymbols) ListExportedByUserRepoAndTag(ctx context.Context, userID int64, owner, repoName, tag string) ([]*models.Symbol, error) {
	if m.OnListExportedByUserRepoAndTag != nil {
		return m.OnListExportedByUserRepoAndTag(ctx, userID, owner, repoName, tag)
	}
	return nil, nil
}

func (m *MockSymbols) FindRelated(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
	if m.OnFindRelated != nil {
		return m.OnFindRelated(ctx, userID, owner, repoName, tag, docVector, limit)
	}
	return nil, nil
}

func (m *MockSymbols) FindRelatedExported(ctx context.Context, userID int64, owner, repoName, tag string, docVector []float32, limit int) ([]*models.Symbol, error) {
	if m.OnFindRelatedExported != nil {
		return m.OnFindRelatedExported(ctx, userID, owner, repoName, tag, docVector, limit)
	}
	return nil, nil
}
