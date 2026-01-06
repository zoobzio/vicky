# Repository Files

Common repository files that every package should include.

## Required Files

```
/
├── LICENSE
├── CONTRIBUTING.md
├── SECURITY.md
└── .gitignore
```

## LICENSE

All packages use the MIT license:

```
MIT License

Copyright (c) [year] zoobzio

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## CONTRIBUTING.md

Keep brief; link to development workflow:

```markdown
# Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make check` to verify
5. Submit a pull request

## Development

- `make help` — list available commands
- `make test` — run all tests
- `make lint` — run linters

## Code of Conduct

Be respectful. We're all here to build good software.
```

## SECURITY.md

Provide clear reporting instructions:

```markdown
# Security Policy

## Reporting a Vulnerability

Please report security vulnerabilities by emailing security@zoobzio.com.

Do not open a public issue for security vulnerabilities.

## Response Timeline

- Acknowledgment: within 48 hours
- Initial assessment: within 1 week
- Fix timeline: depends on severity
```

## .gitignore

Standard Go gitignore:

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test files
*.test
*.out
coverage.out
coverage.html
coverage.txt

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Build
dist/
vendor/
```

## Optional Files

### CHANGELOG.md

Not required — GoReleaser generates changelogs from commit messages. Only maintain manually if:
- You need curated release notes
- Commit history doesn't follow conventional commits

### CODE_OF_CONDUCT.md

Brief conduct statement in CONTRIBUTING.md is sufficient for most packages. Add a full CODE_OF_CONDUCT.md only if the project has a larger community.
