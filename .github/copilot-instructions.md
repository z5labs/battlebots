# Copilot Instructions for Battle Bots

## Project Overview

Battle Bots is a PVP game platform where humans implement autonomous "bots" that battle each other in a 2D battle space. This repository is currently in R&D phase, containing documentation and research artifacts. The game implementation is in early proof-of-concept with two architectural approaches being evaluated: client/server and peer-to-peer.

## Documentation-Driven Development

This project uses Hugo with Docsy theme for documentation at https://z5labs.github.io/battlebots/. All code development will be informed by structured R&D documents.

### Hugo Development Commands

```bash
cd docs
npm install -D postcss postcss-cli autoprefixer  # First time only
hugo server                                       # Dev server
hugo --gc --minify --baseURL "https://z5labs.github.io/battlebots/"  # Production build
hugo mod tidy && hugo mod get -u                  # Module management
```

## R&D Document Workflow

The project follows a three-tier documentation approach before implementation:

### 1. Architecture Decision Records (ADRs)

- **Location**: `docs/content/research_and_development/adrs/`
- **Format**: MADR 4.0.0 standard
- **Naming**: `NNNN-title-with-dashes.md` (e.g., `0001-use-madr.md`)
- **Categories**:
  - `strategic`: Framework choices, authentication strategies, cross-cutting patterns
  - `user-journey`: Feature implementation approaches for specific user flows
  - `api-design`: API endpoint design decisions (pagination, filtering, etc.)
- **Status values**: `proposed` | `accepted` | `rejected` | `deprecated` | `superseded by ADR-XXXX`
- **Creation**: Use `/new-adr` slash command in Claude, which prompts for title, category, and generates sequential numbering

### 2. User Journeys

- **Location**: `docs/content/research_and_development/user-journeys/`
- **Purpose**: Define user flows with Mermaid diagrams and prioritized technical requirements (P0/P1/P2)
- **Status values**: `draft` | `in-review` | `approved` | `implemented` | `deprecated`
- **Creation**: Use `/new-user-journey` slash command in Claude
- **Requirement IDs**: Format `REQ-[CATEGORY]-NNN` (e.g., `REQ-AC-001` for Access Control)

### 3. Analysis Documents

- **Location**: `docs/content/research_and_development/analysis/`
- **Purpose**: Technical research exploring technologies (e.g., `observability/otel-collector/`, `observability/logs/loki/`)
- **Structure**: Organized by topic with comprehensive overviews that inform ADRs

## Git Workflow Conventions

### Branch Naming
```
story/issue-{number}/{description}
```
Example: `story/issue-148/observability-stack-adr`

### Commit Messages
```
type(issue-{number}): description
```
Types: `docs`, `chore`, `story`, `refactor`
Example: `docs(issue-148): research loki`

## Key Architectural Decisions

From existing ADRs:

1. **Observability SDK**: OpenTelemetry (ADR-0002) for vendor-neutral instrumentation
2. **Observability Stack**: Tempo (traces), Mimir (metrics), Loki (logs), Grafana (visualization) - ADR-0003
3. **Language-Agnostic Bot Interface**: Bots run in containers, independent of implementation language
4. **Architecture Evaluation**: Client/server vs P2P approaches still being prototyped

## Development Environment

- **Devcontainer**: Configured with Go 1.24, Hugo 0.140.2 (extended), and protoc
- **Linting**: golangci-lint for Go formatting and linting
- **License**: MIT License (Z5Labs and Contributors)

## CI/CD

- **Docs Deployment**: Pushes to `main` affecting `docs/**` deploy to GitHub Pages
- **PR Previews**: PRs changing `docs/**` get preview deployments at `https://z5labs.github.io/battlebots/pr-preview/pr-{number}/`

## Code Implementation Guidelines

When implementing code (future):

1. **Consult ADRs first**: Check `docs/content/research_and_development/adrs/` for decisions on frameworks, patterns, and approaches
2. **Reference User Journeys**: Understand the user flow and requirements in `docs/content/research_and_development/user-journeys/`
3. **Follow established patterns**: Look at analysis documents for implementation guidance
4. **Document new decisions**: Create ADRs for significant technical choices not already covered
5. **Use OpenTelemetry**: All components should use OpenTelemetry SDK for observability

## Important File Locations

- **ADR Template**: `.claude/commands/new-adr.md`
- **User Journey Template**: `.claude/commands/new-user-journey.md`
- **Documentation Content**: `docs/content/research_and_development/`
- **Hugo Config**: `docs/hugo.yaml`
- **CI Workflows**: `.github/workflows/docs.yaml`, `.github/workflows/docs-preview.yaml`
