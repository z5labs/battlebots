# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Battle Bots is a PVP game platform where humans implement autonomous "bots" that battle each other. Each bot is independent software that reacts to state updates and performs actions in a 2D battle space.

**Current Status**: This repository currently contains documentation and research & development artifacts. The actual game implementation is in early proof-of-concept phase, with two architectural approaches being evaluated: client/server and peer-to-peer.

## Documentation Development

This repository uses Hugo with the Docsy theme for documentation hosted at https://z5labs.github.io/battlebots/

### Building and Running Documentation

```bash
# Navigate to docs directory
cd docs

# Install dependencies (first time only)
npm install -D postcss postcss-cli autoprefixer

# Run Hugo development server
hugo server

# Build production site
hugo --gc --minify --baseURL "https://z5labs.github.io/battlebots/"
```

### Hugo Module Management

```bash
cd docs

# Tidy Hugo modules
hugo mod tidy

# Update Hugo modules
hugo mod get -u
```

## Research & Development Process

The project follows a structured R&D approach with three types of documents:

### 1. Architecture Decision Records (ADRs)

ADRs use the MADR 4.0.0 format and are stored in [docs/content/research_and_development/adrs/](docs/content/research_and_development/adrs/).

**Creating ADRs**: Use the `/new-adr` slash command, which will:
- Prompt for the ADR title and category
- Generate sequential numbering (e.g., 0001, 0002)
- Create the file with the MADR template

**ADR Categories**:
- `strategic`: High-level architectural decisions (frameworks, authentication strategies, cross-cutting patterns)
- `user-journey`: Solutions for specific user journey problems (feature implementation approaches)
- `api-design`: API endpoint design decisions (pagination, filtering, bulk operations)

**Status Values**: `proposed` | `accepted` | `rejected` | `deprecated` | `superseded by ADR-XXXX`

### 2. User Journeys

User journey documents define user flows and technical requirements, stored in [docs/content/research_and_development/user-journeys/](docs/content/research_and_development/user-journeys/).

**Creating User Journeys**: Use the `/new-user-journey` slash command, which will:
- Prompt for the journey title
- Generate sequential numbering
- Create the file with the user journey template including Mermaid diagrams

**Status Values**: `draft` | `in-review` | `approved` | `implemented` | `deprecated`

### 3. Analysis Documents

Research and analysis documents exploring technologies and approaches, stored in [docs/content/research_and_development/analysis/](docs/content/research_and_development/analysis/).

These are structured by topic (e.g., observability/otel-collector, observability/logs/loki) and provide technical research that informs future ADRs.

## Git Workflow

### Branch Naming

Use the format: `story/issue-{number}/{description}`

Example: `story/issue-148/observability-stack-adr`

### Commit Messages

Format: `type(issue-{number}): description`

Types:
- `docs`: Documentation changes
- `chore`: Maintenance tasks
- `story`: User story implementation
- `refactor`: Code restructuring

Example: `docs(issue-148): research loki`

### CI/CD

- **Docs Deployment**: Pushes to `main` that change `docs/**` trigger deployment to GitHub Pages
- **PR Previews**: PRs that change `docs/**` get preview deployments at `https://z5labs.github.io/battlebots/pr-preview/pr-{number}/`

## Development Environment

A devcontainer is configured with:
- Go 1.24
- Hugo 0.140.2 (extended)
- protoc (Protocol Buffers compiler)

The devcontainer uses golangci-lint for Go formatting and linting.

## Documentation Structure

The documentation follows this hierarchy:

```
docs/content/
├── _index.md                        # Homepage
└── research_and_development/
    ├── adrs/                        # Architecture Decision Records
    ├── user-journeys/               # User journey documents
    └── analysis/                    # Technical research documents
        └── observability/           # Example: observability research
            ├── otel-collector/      # OpenTelemetry Collector analysis
            └── logs/loki/           # Loki log storage analysis
```

Each section includes an `_index.md` that provides overview and context.

## Key Architectural Considerations

From the POC user journey and ongoing research:

1. **Architecture Evaluation**: Both client/server and P2P implementations are being prototyped to determine the optimal approach
2. **Language-Agnostic Bot Interface**: The game logic must be independent of bot implementation language
3. **Containerization**: Bots run in containers for isolation and portability
4. **Observability**: Critical for real-time monitoring, debugging, and visualization
5. **Battle Visualization**: Required for displaying battle state and actions

## Future Development

The codebase will expand to include:
- Game server implementation (language TBD)
- Bot interface protocol (gRPC, HTTP, or custom - pending ADR)
- Observability stack (OpenTelemetry Collector + backends - pending ADR)
- Battle visualization frontend
- Bot SDKs for various languages
