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

// GetDefinition returns the definition location(s) for a symbol.
var GetDefinition = rocco.GET("/intel/{owner}/{repo}/{tag}/definition", func(req *rocco.Request[rocco.NoBody]) (wire.DefinitionResponse, error) {
	occurrences := sum.MustUse[contracts.SCIPOccurrences](req.Context)
	documents := sum.MustUse[contracts.Documents](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.DefinitionResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]
	symbol := req.Params.Query["symbol"]

	if symbol == "" {
		return wire.DefinitionResponse{}, ErrMissingSymbol
	}

	definitions, err := occurrences.ListDefinitions(req.Context, userID, owner, repoName, tag, symbol)
	if err != nil {
		return wire.DefinitionResponse{}, err
	}

	// Build path resolver
	pathCache := make(map[int64]string)
	pathResolver := func(docID int64) string {
		if path, ok := pathCache[docID]; ok {
			return path
		}
		doc, err := documents.Get(req.Context, strconv.FormatInt(docID, 10))
		if err != nil || doc == nil {
			return ""
		}
		pathCache[docID] = doc.Path
		return doc.Path
	}

	return transformers.DefinitionsToResponse(symbol, definitions, pathResolver), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("symbol").
	WithSummary("Get definition").
	WithDescription("Returns the definition location(s) for a symbol.").
	WithTags("Code Intelligence").
	WithErrors(ErrMissingSymbol).
	WithAuthentication()

// FindReferences returns all references to a symbol.
var FindReferences = rocco.GET("/intel/{owner}/{repo}/{tag}/references", func(req *rocco.Request[rocco.NoBody]) (wire.ReferencesResponse, error) {
	occurrences := sum.MustUse[contracts.SCIPOccurrences](req.Context)
	documents := sum.MustUse[contracts.Documents](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.ReferencesResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]
	symbol := req.Params.Query["symbol"]

	if symbol == "" {
		return wire.ReferencesResponse{}, ErrMissingSymbol
	}

	includeDefinition := req.Params.Query["include_definition"] == "true"

	var refs []*models.SCIPOccurrence
	if includeDefinition {
		refs, err = occurrences.ListBySymbol(req.Context, userID, owner, repoName, tag, symbol)
	} else {
		refs, err = occurrences.ListReferences(req.Context, userID, owner, repoName, tag, symbol)
	}
	if err != nil {
		return wire.ReferencesResponse{}, err
	}

	// Build path resolver
	pathCache := make(map[int64]string)
	pathResolver := func(docID int64) string {
		if path, ok := pathCache[docID]; ok {
			return path
		}
		doc, err := documents.Get(req.Context, strconv.FormatInt(docID, 10))
		if err != nil || doc == nil {
			return ""
		}
		pathCache[docID] = doc.Path
		return doc.Path
	}

	return transformers.ReferencesToResponse(symbol, refs, pathResolver), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("symbol", "include_definition").
	WithSummary("Find references").
	WithDescription("Returns all references to a symbol.").
	WithTags("Code Intelligence").
	WithErrors(ErrMissingSymbol).
	WithAuthentication()

// FindImplementations returns types implementing an interface.
var FindImplementations = rocco.GET("/intel/{owner}/{repo}/{tag}/implementations", func(req *rocco.Request[rocco.NoBody]) (wire.ImplementationsResponse, error) {
	scipSymbols := sum.MustUse[contracts.SCIPSymbols](req.Context)
	relationships := sum.MustUse[contracts.SCIPRelationships](req.Context)
	occurrences := sum.MustUse[contracts.SCIPOccurrences](req.Context)
	documents := sum.MustUse[contracts.Documents](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.ImplementationsResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]
	symbol := req.Params.Query["symbol"]

	if symbol == "" {
		return wire.ImplementationsResponse{}, ErrMissingSymbol
	}

	// Find the symbol to get its ID
	sym, err := scipSymbols.GetBySymbol(req.Context, userID, owner, repoName, tag, symbol)
	if err != nil {
		return wire.ImplementationsResponse{}, err
	}
	if sym == nil {
		return wire.ImplementationsResponse{Symbol: symbol, Implementations: []wire.ImplementationInfo{}, Total: 0}, nil
	}

	// Get relationships for this symbol
	rels, err := relationships.ListImplementations(req.Context, sym.ID)
	if err != nil {
		return wire.ImplementationsResponse{}, err
	}

	// Build resolver for target symbols
	symbolResolver := func(targetSymbol string) (*models.SCIPSymbol, *wire.Location) {
		targetSym, err := scipSymbols.GetBySymbol(req.Context, userID, owner, repoName, tag, targetSymbol)
		if err != nil || targetSym == nil {
			return nil, nil
		}

		// Find definition location
		defs, err := occurrences.ListDefinitions(req.Context, userID, owner, repoName, tag, targetSymbol)
		if err != nil || len(defs) == 0 {
			return targetSym, nil
		}

		doc, err := documents.Get(req.Context, strconv.FormatInt(defs[0].DocumentID, 10))
		if err != nil || doc == nil {
			return targetSym, nil
		}

		loc := &wire.Location{
			Path:      doc.Path,
			StartLine: defs[0].StartLine,
			StartCol:  defs[0].StartCol,
			EndLine:   defs[0].EndLine,
			EndCol:    defs[0].EndCol,
		}
		return targetSym, loc
	}

	return transformers.RelationshipsToImplementationsResponse(symbol, rels, symbolResolver), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("symbol").
	WithSummary("Find implementations").
	WithDescription("Returns types implementing an interface or trait.").
	WithTags("Code Intelligence").
	WithErrors(ErrMissingSymbol).
	WithAuthentication()

// ListSymbols returns symbols in a file or version.
var ListSymbols = rocco.GET("/intel/{owner}/{repo}/{tag}/symbols", func(req *rocco.Request[rocco.NoBody]) (wire.SymbolListResponse, error) {
	scipSymbols := sum.MustUse[contracts.SCIPSymbols](req.Context)
	documents := sum.MustUse[contracts.Documents](req.Context)
	occurrences := sum.MustUse[contracts.SCIPOccurrences](req.Context)

	userID, err := strconv.ParseInt(req.Identity.ID(), 10, 64)
	if err != nil {
		return wire.SymbolListResponse{}, err
	}

	owner := req.Params.Path["owner"]
	repoName := req.Params.Path["repo"]
	tag := req.Params.Path["tag"]
	path := req.Params.Query["path"]

	var symbols []*models.SCIPSymbol

	if path != "" {
		// Get symbols for a specific file
		doc, err := documents.GetByUserRepoTagAndPath(req.Context, userID, owner, repoName, tag, path)
		if err != nil {
			return wire.SymbolListResponse{}, err
		}
		if doc == nil {
			return wire.SymbolListResponse{Symbols: []wire.SCIPSymbolInfo{}, Total: 0}, nil
		}
		symbols, err = scipSymbols.ListByDocument(req.Context, doc.ID)
		if err != nil {
			return wire.SymbolListResponse{}, err
		}
	} else {
		// Get all symbols for the version
		symbols, err = scipSymbols.ListByUserRepoAndTag(req.Context, userID, owner, repoName, tag)
		if err != nil {
			return wire.SymbolListResponse{}, err
		}
	}

	// Build location resolver
	locationResolver := func(sym *models.SCIPSymbol) *wire.Location {
		defs, err := occurrences.ListDefinitions(req.Context, userID, owner, repoName, tag, sym.Symbol)
		if err != nil || len(defs) == 0 {
			return nil
		}

		doc, err := documents.Get(req.Context, strconv.FormatInt(defs[0].DocumentID, 10))
		if err != nil || doc == nil {
			return nil
		}

		return &wire.Location{
			Path:      doc.Path,
			StartLine: defs[0].StartLine,
			StartCol:  defs[0].StartCol,
			EndLine:   defs[0].EndLine,
			EndCol:    defs[0].EndCol,
		}
	}

	return transformers.SCIPSymbolsToListResponse(symbols, locationResolver), nil
}).WithPathParams("owner", "repo", "tag").
	WithQueryParams("path", "kind").
	WithSummary("List symbols").
	WithDescription("Returns symbols in a file or version.").
	WithTags("Code Intelligence").
	WithErrors(ErrMissingSymbol).
	WithAuthentication()
