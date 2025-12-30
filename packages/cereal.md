# Cereal

Content-type aware marshaling with sanitization for Go.

[GitHub](https://github.com/zoobzio/cereal)

## Vision

Unified serialization with security built in. Struct tags control field protection - encryption, hashing, masking, redaction. Cereal handles format. Serializer handles safety. Marshal once, sanitize automatically.

## Design Decisions

**Two-layer architecture**
Cereal interface for format (JSON/XML/YAML/MessagePack). Serializer[T] for sanitization. Layers compose. Use Cereal alone for simple cases, Serializer when field protection needed.

**Tag-based field plans**
Sentinel scans type once at Serializer creation. Builds cached field plans. Zero reflection on hot path. Plans describe which fields need what operations.

**Clone before sanitize**
`Cloner[T]` constraint ensures deep copy before mutation. Original never touched. Sanitization operates on the clone.

**Precedence order**
`encrypt > hash > redact > mask`. Only one operation per field. Highest precedence wins when multiple tags present.

**Reversibility distinction**
Encryption is reversible (decrypt on unmarshal). Hash, redact, mask are one-way. Restore only processes encrypt fields.

**Contextual encryption**
`encrypt:"pii"` references registered encryptor by name. Multiple contexts (e.g., "pii", "secrets") with different keys. Context passed through to encryptor.

**Escape hatches**
`Sanitizable` and `Desanitizable` interfaces bypass reflection entirely. For code-generated or performance-critical paths.

## Sanitization Tags

| Tag | Example | Effect | Reversible |
|-----|---------|--------|------------|
| `encrypt` | `encrypt:"pii"` | AES-GCM/RSA encryption | Yes |
| `hash` | `hash:"sha256"` | Hex-encoded hash | No |
| `redact` | `redact:"***"` | Replace with literal | No |
| `mask` | `mask:"email"` | Content-aware partial masking | No |

## Built-in Maskers

| Type | Example |
|------|---------|
| ssn | `123-45-6789` → `***-**-6789` |
| email | `alice@example.com` → `a***@example.com` |
| phone | `(555) 123-4567` → `(***) ***-4567` |
| card | `4111111111111111` → `************1111` |
| ip | `192.168.1.100` → `192.168.xxx.xxx` |
| uuid | `550e8400-...` → `550e8400-****-****-****-************` |
| iban | `GB82WEST12345698765432` → `GB82************5432` |
| name | `John Smith` → `J*** S****` |

## Providers

| Provider | Content-Type |
|----------|--------------|
| `pkg/json` | application/json |
| `pkg/xml` | application/xml |
| `pkg/yaml` | application/yaml |
| `pkg/msgpack` | application/msgpack |

## Internal Architecture

```
Marshal:
Original → Clone() → Sanitize (tag-based) → Cereal.Marshal → []byte

Unmarshal:
[]byte → Cereal.Unmarshal → Restore (decrypt only) → *T
```

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Package documentation |
| `cereal.go` | Cereal interface |
| `serializer.go` | Serializer[T] with sanitization |
| `encrypt.go` | AES, RSA, Envelope encryptors |
| `hash.go` | SHA256, SHA512 hashers |
| `mask.go` | Content-aware maskers (8 types) |
| `sanitize.go` | Tag parsing, field plans, Apply/Restore |
| `clone.go` | Cloner[T] interface |
| `options.go` | WithEncryptor, WithHasher |
| `signals.go` | Capitan events |

## Current State / Direction

Stable. Core sanitization complete.

Future considerations:
- Additional mask patterns
- Streaming serialization

## Framework Context

**Dependencies**: sentinel (type metadata), capitan (events).

**Role**: Security-first serialization. Struct tags declare intent, Serializer enforces protection. Intended for use by herald for secure message transport (not yet integrated).
