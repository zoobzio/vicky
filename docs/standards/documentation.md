# Documentation

All packages maintain documentation in a `docs/` directory with consistent structure and frontmatter for site publishing.

## Directory Structure

Documentation files use numbered prefixes to control ordering:

```
docs/
├── 1.overview.md              # Package overview
├── 2.learn/                   # Learning materials
│   ├── 1.quickstart.md
│   ├── 2.concepts.md
│   └── ...
├── 3.guides/                  # How-to guides
│   ├── 1.[topic].md
│   └── ...
├── 4.cookbook/                # Recipes and patterns
│   └── ...
└── 5.reference/               # API reference
    ├── 1.api.md
    └── ...
```

## Frontmatter

Documentation files include frontmatter for site publishing:

```yaml
---
title: Article Title
description: One-line description
author: zoobzio
published: YYYY-MM-DD
updated: YYYY-MM-DD
tags:
  - Tag1
  - Tag2
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
