# Sctx

Certificate-based security contexts for Go.

[GitHub](https://github.com/zoobzio/sctx)

## Vision

Zero-trust authentication using existing PKI. Transform mTLS certificates into typed, permission-aware security contexts. No external identity providers. No JWTs. Cryptographically signed, delegatable tokens derived from certificates you already have.

## Design Decisions

**Singleton admin**
One AdminService per application. Single source of truth for token signing and context cache. Prevents competing authorities.

**Certificate-centric**
All authorization derived from certificate properties (CN, O, OU). Policy-driven transformation from cert fields to permissions. No external config drift.

**Assertion-based proof**
Require signed assertion before token generation. Proves client holds private key. Short-lived (1 minute), single-use with nonce validation. Prevents replay attacks.

**Generic contexts**
`Context[M]` carries type-safe application metadata. No `interface{}` gymnastics.

**Delegatable guards**
Guards are permission validators created by token holders. Services create guards for others without sharing tokens. Can only delegate permissions you possess.

**Dual crypto support**
Ed25519 (default, ~30% faster) and ECDSA P-256 (FIPS compliance). Choose based on requirements.

**Capitan events**
All operations emit structured events. Token lifecycle, guard validation, cache operations, assertions, revocation. Observable without changing API.

## Internal Architecture

```
Certificate → [Assertion validation] → [Policy execution] → Context[M]
                                                                 ↓
                                                          [Cache by fingerprint]
                                                                 ↓
                                                         [Sign token payload]
                                                                 ↓
                                                           SignedToken
```

Token format: `base64(payload):base64(signature)`

Cache uses SHA-256 certificate fingerprint as key. O(1) lookups. Background cleanup of expired entries.

## Code Organisation

| File | Responsibility |
|------|----------------|
| `api.go` | Public interfaces (Admin, Guard, ContextPolicy) |
| `admin.go` | AdminService, token generation, guard creation |
| `context.go` | Context[M], token payload serialisation |
| `assertion.go` | Assertion creation and validation |
| `guards.go` | Guard implementation, enrichment guards |
| `cache.go` | ContextCache interface, memory implementation |
| `crypto.go` | Ed25519/ECDSA signing and verification |
| `events.go` | Capitan event definitions |

## Security Model

| Attack | Mitigation |
|--------|------------|
| Token forgery | Signature verification |
| Replay | Nonce validation on assertions |
| Key theft | Assertions prove key possession |
| Context hijacking | Fingerprint binds token to certificate |
| Expiration bypass | Token expiry ≤ context expiry ≤ cert notAfter |
| Permission escalation | Guards prevent delegating unpossessed permissions |

## Current State / Direction

Stable. Core authentication flow complete. Revocation working.

Future considerations:
- Distributed cache adapters (Redis, etcd)
- Certificate rotation helpers

## Framework Context

**Dependencies**: capitan (events), testcontainers (testing only).

**Role**: Security layer for mTLS environments. Certificate-driven authentication without external identity providers. Integrates with existing PKI infrastructure.
