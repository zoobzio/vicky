# Documentation

All packages maintain documentation in a `docs/` directory with consistent structure and frontmatter for site publishing.

## Directory Structure

Documentation files use numbered prefixes to control ordering:

```
docs/
├── 1.learn/                   # Learning materials
│   ├── 1.overview.md
│   ├── 2.quickstart.md
│   ├── 3.concepts.md
│   ├── 4.architecture.md
│   └── ...
├── 2.guides/                  # How-to guides
│   ├── 1.[topic].md
│   └── ...
├── 3.integrations/            # Ecosystem connections (optional)
│   ├── 1.[tool].md
│   └── ...
└── 4.reference/               # API reference
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

## File Specifications

### Overview (`1.learn/1.overview.md`)

The entry point. Lets a reader determine in 60 seconds whether this package suits their needs.

**Must answer:**

- What is this? (one sentence)
- What idea or question motivated it?
- What does it provide?
- What can you build with it?
- Where to go next?

**Must not contain:**

- Installation instructions
- API details
- Deep architectural explanation
- Step-by-step tutorials

**Structure:**

```markdown
# Overview

{One paragraph: what it is + core value proposition}

## The Idea

{The question or insight that motivated this package}

## The Implementation

{What the package provides to realize the idea - scannable list}

## What It Enables

{Concrete outcomes or ecosystem links - what you build with it}

## Next Steps

{Links to quickstart, concepts, reference}
```

### Quickstart (`1.learn/2.quickstart.md`)

Gets the reader from zero to productive. By the end, a user should understand enough to build with the package—not just have it installed.

Blends foundational concepts with core functionality. Concepts and Architecture serve power users who want deeper understanding; Guides provide detailed breakdowns of individual features. Quickstart ensures users know the basics before they reach either.

**Must contain:**

- Installation instructions
- Basic usage example (complete, copy-paste-run)
- Sections for foundational topics—brief explanation, minimal snippet, link to relevant guide
- Links to concepts, architecture, and reference for deeper understanding

**Must not contain:**

- Exhaustive API coverage
- Deep explanations of internals
- Edge cases or advanced configuration

**Structure:**

```markdown
# Quickstart

## Installation

{Package manager command, version requirements}

## Basic Usage

{Minimal complete example - imports, code, expected output}

## {Foundational Topic}

{Brief explanation of concept, minimal snippet, link to guide}

## {Foundational Topic}

{Brief explanation of concept, minimal snippet, link to guide}

## Next Steps

{Links to concepts, architecture, reference}
```

### Concepts (`1.learn/3.concepts.md`)

Defines the abstractions. Fleshes out ideas that Quickstart introduces, giving users the mental models to reason about the package.

**Must contain:**

- The core abstractions and what they represent
- Why the abstractions are shaped the way they are
- How abstractions relate to each other
- Cross-links to Reference sections for each abstraction that maps to an API element (use `#section` anchors for direct routing)

**Must not contain:**

- Complete type definitions (that's Reference)
- Implementation details (that's Architecture)
- Step-by-step usage (that's Quickstart/Guides)

**Structure:**

```markdown
# Concepts

{Brief intro - what mental models this page establishes}

## {Abstraction}

{What it represents, why it exists, how to think about it}

## {Abstraction}

{What it represents, why it exists, how to think about it}

## Next Steps

{Links to architecture, reference}
```

### Architecture (`1.learn/4.architecture.md`)

For power users and contributors who want to understand the implementation. Explains how the package works internally, not how to use it.

**Must contain:**

- Component structure and interactions
- Data flow and key algorithms
- Design rationale framed as Q&A (anticipating common questions, not fabricating tradeoffs)
- Performance characteristics

**Must not contain:**

- Usage examples (that's Quickstart/Guides)
- Concept explanations (link to Concepts instead)
- API documentation (link to Reference instead)

**Structure:**

```markdown
# Architecture

{Brief intro — who this is for, what it covers}

## Component Overview

{Diagram showing major components and their relationships}

## {Internal Mechanism}

{How it works, the algorithm, data flow}

## {Internal Mechanism}

{How it works, the algorithm, data flow}

## Design Q&A

{Common questions about why things work the way they do}

## Performance

{Characteristics, complexity, benchmarks}

## Next Steps

{Links to guides, reference}
```

### Guides (`2.guides/`)

Detailed breakdowns of individual features and workflows. Assumes the reader has completed the Quickstart.

**Required guides:**

- **Testing** (`testing.md`) — how to test code that uses this package (isolation, fixtures, mocking)
- **Troubleshooting** (`troubleshooting.md`) — common errors, edge cases, debugging steps

**Package-specific guides:**

Emerge from features introduced in Quickstart. If Quickstart gives a brief explanation with a "see guide for more" link, that guide should exist.

Criteria for a feature needing its own guide:
- Has multiple modes or options
- Requires understanding of edge cases
- Benefits from extended examples
- Users frequently ask "how do I..." about it

**Structure:**

```markdown
# {Feature} Guide

{Brief intro — what this guide covers, prerequisites}

## {Aspect}

{Explanation with examples}

## {Aspect}

{Explanation with examples}

## Next Steps

{Links to related guides, reference}
```

### Integrations (`3.integrations/`) — Optional

Documents how other packages consume this package's output. Highlights the unique value this package provides to the ecosystem.

**When to include:**

- Package produces structured output consumed by other tools
- Multiple downstream packages exist or are planned
- The integration pattern is non-obvious

**When to skip:**

- Package is a standalone tool
- No downstream consumers exist
- Integration is trivial (just import and call)

**Required files (when section exists):**

- **Per-tool pages** (`1.[tool].md`) — one page per downstream tool

The README lists integrations and explains the pattern. No separate overview needed.

**Per-tool pages must contain:**

- What the tool does (one sentence)
- The problem it solves (optional, when non-obvious)
- The pipeline: types → package extraction → tool transformation → output
- What this package provides (table mapping tool needs to package outputs)
- Link to the tool's own documentation

**Per-tool pages must not contain:**

- Full tool documentation (that belongs in the tool's repo)
- Extended tutorials
- Configuration options for the downstream tool

**Structure:**

```markdown
# {Tool}

{One sentence: what the tool does}

## The Problem (optional)

{What pain point does this integration address}

## The Pipeline

{Show the full flow: define types → extract metadata → transform → result}

{Include the actual metadata structure that gets extracted}

## What {Package} Provides

| {Tool} needs | {Package} provides |
|--------------|---------------------|
| ... | ... |

## Learn More

- [{Tool} repository]({link})
```

### Reference (`4.reference/`)

Complete API documentation. Authoritative source for signatures, types, and constants.

**Required files:**

- **API Reference** (`1.api.md`) — functions only
- **Types Reference** (`2.types.md`) — types, constants, enums

API Reference links to Types Reference for return types and parameters. No duplication between files.

**Function entries must have:**

- Signature (code block)
- Brief description
- Panic/error conditions
- Example

**Type entries must have:**

- Definition (code block)
- Field table (name, type, description)
- Usage notes where helpful

**Structure (API Reference):**

```markdown
# API Reference

## {Category}

### {Function}

\`\`\`go
func Signature() ReturnType
\`\`\`

{Description}

{Panic/error conditions}

\`\`\`go
// Example
\`\`\`
```

**Structure (Types Reference):**

```markdown
# Types Reference

## {Type}

{Brief description}

\`\`\`go
type Definition struct { ... }
\`\`\`

| Field | Type | Description |
|-------|------|-------------|
| ... | ... | ... |
```
