# README

Package READMEs maintain consistency without prescription. The structure provides rhythm; the library's character dictates voice.

## Principles

1. **Essence over problem** — Lead with what makes this library *this library*, not a generic problem statement
2. **Show, don't frame** — Code speaks louder than "The Problem / The Solution"
3. **Voice from nature** — Let the library's character dictate tone
4. **Structure as scaffolding** — Sections provide rhythm, not script

## Structure

```
1. HEADER
   - Title
   - Badges (standardized)
   - Tagline (one line, what it *is*)
   - Supporting sentence (what you *do* with it)

2. ESSENCE (name varies)
   - The unique hook for this library
   - Minimal code showing the core insight
   - Section name should be specific, not generic
   - Examples: "Three Operations", "One Interface", "Scan Once"

3. INSTALL
   - go get
   - Requirements

4. QUICK START
   - Single complete example with imports
   - Runnable, demonstrates core usage
   - Comments explain, code shows

5. CAPABILITIES (optional, name varies)
   - Only if distinct modes/features worth highlighting
   - Skip if Quick Start already covers it

6. WHY [NAME]?
   - Bullet points
   - Concrete benefits, not marketing
   - Position *after* showing what it does

7. DOCUMENTATION
   - Categorized links (Learn / Guides / Cookbook / Reference)
   - Consistent hierarchy across packages

8. CONTRIBUTING
   - Brief
   - Link to CONTRIBUTING.md

9. LICENSE
   - One line
```

## Naming the Essence Section

The essence section name should complete: "[Library] is about ___."

| Library | Essence | Section Name |
|---------|---------|--------------|
| capitan | Three core functions | "Three Operations" |
| pipz | One composable interface | "One Interface" |
| sentinel | Scan once, cache forever | "Scan Once" |

Avoid: "The Problem", "The Solution", "Overview", "Introduction", "How It Works"

## Code Block Principles

1. **Essence block**: Minimal, conceptual, shows *the insight*
2. **Quick Start block**: Complete, runnable, shows *the workflow*
3. **Capability blocks**: Focused, shows *the variation*

## Length Target

- 100-200 lines for focused libraries
- 200-300 lines for libraries with multiple modes
- If exceeding 300 lines, content likely belongs in docs/

## Anti-Patterns

1. **Template voice**: Every README sounds the same
2. **Problem-first framing**: Generic complaints before showing value
3. **Feature laundry lists**: Long bullet points that could be docs
4. **Duplicate examples**: Quick Start vs Quick Example vs Real-World Example
5. **API reference in README**: Tables belong in docs/reference
