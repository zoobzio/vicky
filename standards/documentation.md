# Documentation

All packages maintain documentation in a `docs/` directory with consistent structure and frontmatter for site publishing.

## Directory Structure

```
docs/
├── index.md             # Package overview
├── getting-started.md   # Quick start guide
├── [topic].md           # Topic-specific documentation
└── ...
```

## Frontmatter

Documentation files include frontmatter for site publishing:

```yaml
---
title: Package Name
description: Brief description of the package
package: package-name
---
```

## Tone and Rhythm

Documentation should maintain a consistent tone across packages:

- Clear and direct
- Technical but accessible
- Example-driven where appropriate

Packages need not follow a rigid template but should feel cohesive when read alongside other package documentation.

## Content Guidelines

- Focus on usage and practical examples
- Explain the "why" alongside the "how"
- Keep examples minimal but complete
- Cross-reference related packages where relevant
