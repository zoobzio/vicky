# Layer 2: Validation & Serialisation

Transform and validate data using schemas.

## Rule

Packages in this layer may depend on Layers 0-1.

## Packages

| Package | Role | Framework Deps |
|---------|------|----------------|
| [astql](../../packages/astql.md) | SQL query building | dbml |
| [vectql](../../packages/vectql.md) | Vector query building | vdml |
| [cereal](../../packages/cereal.md) | Sanitised marshaling | sentinel, capitan |

## Package Details

### astql

Type-safe SQL query builder. Schema-validated, injection-resistant.

**Key capabilities:**
- Three-stage pipeline: Validation → AST → Rendering
- DBML schema as allowlist for all identifiers
- Multi-layer security (schema validation, identifier rules, keyword blocking, parameterised queries)
- Four providers: PostgreSQL, MySQL, SQLite, SQL Server
- Panicking and error-returning variants (T/TryT, F/TryF)

**Framework dependencies:**
- dbml: Schema parsing and validation

**Why Layer 2:** Query validation layer. Generates SQL, doesn't execute. Used by soy.

### vectql

Query builder for vector databases. AST-based like astql.

**Key capabilities:**
- Two modes: Schemaless or VDML-validated
- Operations: Search, Upsert, Delete, Fetch, Update
- Rich filter system (conditions, groups, ranges, geo)
- Complexity limits (MaxFilterDepth, MaxBatchSize, MaxTopK)
- Four providers: Pinecone, Qdrant, Milvus, Weaviate

**Framework dependencies:**
- vdml: Optional schema validation

**Why Layer 2:** Vector database equivalent of astql. Provider-agnostic query construction.

### cereal

Content-type aware marshaling with sanitisation.

**Key capabilities:**
- Two-layer architecture: Cereal (format) + Serializer (security)
- Tag-based sanitisation: encrypt, hash, redact, mask
- Eight built-in maskers (SSN, email, phone, card, IP, UUID, IBAN, name)
- Clone before sanitise (original never touched)
- Four providers: JSON, XML, YAML, MessagePack

**Framework dependencies:**
- sentinel: Type metadata for field plans
- capitan: Signal emission

**Why Layer 2:** Security-first serialisation. Field tags declare intent. Used by herald for secure transport.

## Security Model

Layer 2 establishes defence in depth:

```
                    ┌─────────────────┐
User Input ────────►│  Schema Checks  │ astql: DBML allowlist
                    └────────┬────────┘ vectql: VDML allowlist
                             │
                             ▼
                    ┌─────────────────┐
                    │ Identifier Rules│ Alphanumeric only
                    └────────┬────────┘ Keyword blocking
                             │
                             ▼
                    ┌─────────────────┐
                    │   AST Building  │ Parameterised values
                    └────────┬────────┘ Complexity limits
                             │
                             ▼
                    ┌─────────────────┐
                    │  Sanitisation   │ cereal: encrypt/hash/mask
                    └────────┬────────┘
                             │
                             ▼
                    Safe Output
```

Invalid data never enters the system. Validation happens at construction, not execution.

## Architectural Significance

Layer 2 transforms the schema languages from Layer 0 into operational tools:

| Schema Language | Query Builder | Storage Layer |
|-----------------|---------------|---------------|
| dbml | astql | soy |
| vdml | vectql | (future) |

And provides secure serialisation for all transport:

| Transport | Serialiser |
|-----------|------------|
| herald (messaging) | cereal |
| rocco (HTTP) | cereal |
