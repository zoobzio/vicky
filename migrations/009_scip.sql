-- +goose Up
CREATE TABLE scip_symbols (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    symbol TEXT NOT NULL,
    kind INT NOT NULL,
    display_name TEXT,
    documentation TEXT[],
    enclosing_symbol TEXT,
    signature_documentation JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE scip_occurrences (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    owner TEXT NOT NULL,
    repo_name TEXT NOT NULL,
    tag TEXT NOT NULL,
    symbol TEXT NOT NULL,
    symbol_roles INT NOT NULL,
    start_line INT NOT NULL,
    start_col INT NOT NULL,
    end_line INT NOT NULL,
    end_col INT NOT NULL,
    syntax_kind INT,
    enclosing_range INT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE scip_relationships (
    id BIGSERIAL PRIMARY KEY,
    scip_symbol_id BIGINT NOT NULL REFERENCES scip_symbols(id) ON DELETE CASCADE,
    target_symbol TEXT NOT NULL,
    is_reference BOOLEAN NOT NULL DEFAULT false,
    is_implementation BOOLEAN NOT NULL DEFAULT false,
    is_type_definition BOOLEAN NOT NULL DEFAULT false,
    is_definition BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_scip_symbols_document_id ON scip_symbols(document_id);
CREATE INDEX idx_scip_symbols_user_id ON scip_symbols(user_id);
CREATE INDEX idx_scip_symbols_symbol ON scip_symbols(symbol);
CREATE INDEX idx_scip_symbols_kind ON scip_symbols(kind);
CREATE INDEX idx_scip_symbols_lookup ON scip_symbols(user_id, owner, repo_name, tag);
CREATE INDEX idx_scip_symbols_enclosing ON scip_symbols(enclosing_symbol);

CREATE INDEX idx_scip_occurrences_document_id ON scip_occurrences(document_id);
CREATE INDEX idx_scip_occurrences_user_id ON scip_occurrences(user_id);
CREATE INDEX idx_scip_occurrences_symbol ON scip_occurrences(symbol);
CREATE INDEX idx_scip_occurrences_lookup ON scip_occurrences(user_id, owner, repo_name, tag);

CREATE INDEX idx_scip_relationships_symbol_id ON scip_relationships(scip_symbol_id);
CREATE INDEX idx_scip_relationships_target ON scip_relationships(target_symbol);

-- +goose Down
DROP TABLE scip_relationships;
DROP TABLE scip_occurrences;
DROP TABLE scip_symbols;
