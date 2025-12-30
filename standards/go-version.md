# Go Version

All packages must target Go 1.24 as the minimum supported version.

## go.mod

```go
module github.com/zoobzio/[package]

go 1.24
```

## Rationale

A consistent Go version floor ensures:

- Uniform language features across all packages
- Predictable dependency resolution
- Simplified CI matrix configuration
