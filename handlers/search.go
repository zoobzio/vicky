package handlers

import (
	"strconv"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sum"
	"github.com/zoobzio/vicky/wire"
	"github.com/zoobzio/vicky/contracts"
	"github.com/zoobzio/vicky/models"
	"github.com/zoobzio/vicky/transformers"
)

// SearchChunks performs semantic search across chunks.
var SearchChunks = rocco.GET("/search/{owner}/{repo}/{tag}", func(req *rocco.Request[rocco.NoBody]) (wire.SearchResponse, error) {
	chunks := sum.MustUse[contracts.Chunks](req.Context)
	embedder := sum.MustUse[contracts.Embedder](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.SearchResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]

	query := req.Params.Query["q"]
	if query == "" {
		return wire.SearchResponse{}, ErrMissingQuery
	}

	limit := 10
	if l := req.Params.Query["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	vectors, err := embedder.EmbedQuery(req.Context, []string{query})
	if err != nil {
		return wire.SearchResponse{}, err
	}
	queryVector := vectors[0]

	results, err := chunks.Search(req.Context, userID, owner, repoName, tag, queryVector, limit)
	if err != nil {
		return wire.SearchResponse{}, err
	}

	return transformers.ChunksToSearchResponse(query, results), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("q", "limit", "kind").
	WithSummary("Search chunks").
	WithDescription("Performs semantic search across code and documentation chunks.").
	WithTags("Search").
	WithErrors(ErrMissingQuery).
	WithAuthentication()

// SearchSymbols finds symbols related to a query.
var SearchSymbols = rocco.GET("/search/{owner}/{repo}/{tag}/symbols", func(req *rocco.Request[rocco.NoBody]) (wire.SymbolSearchResponse, error) {
	symbols := sum.MustUse[contracts.Symbols](req.Context)
	embedder := sum.MustUse[contracts.Embedder](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.SymbolSearchResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]

	query := req.Params.Query["q"]
	if query == "" {
		return wire.SymbolSearchResponse{}, ErrMissingQuery
	}

	limit := 10
	if l := req.Params.Query["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	exportedOnly := req.Params.Query["exported"] == "true"

	vectors, err := embedder.EmbedQuery(req.Context, []string{query})
	if err != nil {
		return wire.SymbolSearchResponse{}, err
	}
	queryVector := vectors[0]

	var results []*models.Symbol
	if exportedOnly {
		results, err = symbols.FindRelatedExported(req.Context, userID, owner, repoName, tag, queryVector, limit)
	} else {
		results, err = symbols.FindRelated(req.Context, userID, owner, repoName, tag, queryVector, limit)
	}
	if err != nil {
		return wire.SymbolSearchResponse{}, err
	}

	return transformers.SymbolsToSearchResponse(query, results), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("q", "limit", "exported").
	WithSummary("Search symbols").
	WithDescription("Finds code symbols related to a query.").
	WithTags("Search").
	WithErrors(ErrMissingQuery).
	WithAuthentication()

// FindSimilarDocuments finds documents similar to a given document.
var FindSimilarDocuments = rocco.GET("/search/{owner}/{repo}/{tag}/similar", func(req *rocco.Request[rocco.NoBody]) (wire.SimilarDocumentsResponse, error) {
	documents := sum.MustUse[contracts.Documents](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.SimilarDocumentsResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]
	path := req.Params.Query["path"]

	limit := 10
	if l := req.Params.Query["limit"]; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Get the source document's vector
	sourceDoc, err := documents.GetByUserRepoTagAndPath(req.Context, userID, owner, repoName, tag, path)
	if err != nil {
		return wire.SimilarDocumentsResponse{}, err
	}

	// Find similar documents
	results, err := documents.FindSimilarInVersion(req.Context, userID, owner, repoName, tag, sourceDoc.Vector, limit)
	if err != nil {
		return wire.SimilarDocumentsResponse{}, err
	}

	return transformers.DocumentsToSimilarResponse(results, sourceDoc.ID), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("path", "limit").
	WithSummary("Find similar documents").
	WithDescription("Finds documents similar to a given document path.").
	WithTags("Search").
	WithAuthentication()
