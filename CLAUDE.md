# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Battle Bots is a PVP game platform for autonomous players. Users implement independent bots that battle other bots by reacting to state updates and performing actions. This repository contains the project's documentation and research & design (R&D) artifacts.

## Common Development Commands

### Documentation Site

The project uses Hugo (extended) with the Docsy theme for documentation.

```bash
# Build the documentation site
cd docs
hugo --gc --minify

# Run the documentation server locally
cd docs
hugo server

# Install documentation prerequisites (required before first build)
cd docs
npm install -D postcss postcss-cli autoprefixer
```

### Development Environment

The project includes a devcontainer configuration with Go 1.24, Hugo, and protoc. When working locally, ensure you have:
- Go 1.24+
- Hugo Extended 0.140.2+ (CI uses 0.152.2)
- golangci-lint for code formatting/linting

## R&D Process and Architecture

This repository follows a structured Research & Design workflow. **Always follow this process when designing new features:**

### 1. Document the User Journey First

Before designing solutions, create a user journey document:
- Use `/new-user-journey` to create standardized documentation
- Include Mermaid flow diagrams to visualize user interactions
- Define prioritized technical requirements (P0/P1/P2)
- Location: `docs/content/r&d/user-journeys/`
- Format: `NNNN-title-with-dashes.md`

### 2. Design Solutions with ADRs

After documenting user needs, create ADRs for architectural decisions:
- Use `/new-adr` to create MADR 4.0.0 format documents
- Location: `docs/content/r&d/adrs/`
- Format: `NNNN-title-with-dashes.md`

**ADR Categories:**
- **Strategic**: High-level decisions affecting the entire system (frameworks, authentication strategies, cross-cutting patterns)
- **User Journey**: Solutions for specific user journey problems (feature implementation approaches)
- **API Design**: API endpoint implementation decisions (pagination, filtering, bulk operations)

**ADR Status Values:**
- `proposed` - Under consideration
- `accepted` - Approved, should be implemented
- `rejected` - Considered but not approved
- `deprecated` - No longer relevant
- `superseded by ADR-XXXX` - Replaced by newer ADR

### 3. Document Component ADRs and APIs

For each major component or API identified:
- Create specific ADRs for technical decisions
- Document APIs with request/response schemas, authentication, business logic flows
- Use Mermaid diagrams for complex flows

### Documentation Structure

R&D documentation is organized as:
- **User Journeys** (`docs/content/r&d/user-journeys/`) - User experience flows with technical requirements
- **ADRs** (`docs/content/r&d/adrs/`) - Architectural Decision Records
- **APIs** (`docs/content/r&d/apis/`) - REST API documentation
- **Analysis** (`docs/content/r&d/analysis/`) - Research and analysis

## Project Conventions

### Issue Management

Use the "Story" issue template for new work:
- Title format: `story(subject): short description`
- Include description, acceptance criteria, and related issues
- Label: `story`

### Code Formatting

- **Go**: Use `golangci-lint fmt --stdin` for formatting
- **License**: MIT License, author is "Z5Labs and Contributors"

### Dependencies

- Renovate runs before 4am and targets the `main` branch
- Go modules (including indirect dependencies) are auto-updated
- Ignores: `dagger/**`, `docs/**`

## Working with Documentation

### Creating New Documents

**User Journeys:**
1. Run `/new-user-journey`
2. Include user personas, goals, and pain points
3. Create Mermaid flowcharts for the journey
4. Define requirements with REQ-[CATEGORY]-NNN format:
   - `AC` = Access Control
   - `AN` = Analytics
5. Set priorities: P0 (Must Have), P1 (Should Have), P2 (Nice to Have)

**ADRs:**
1. Run `/new-adr`
2. Select appropriate category (strategic/user-journey/api-design)
3. Document context, decision drivers, considered options
4. Include pros/cons for each option
5. Document decision outcome and consequences
6. Add confirmation method

### Branching Strategy

- Main branch: `main`
- Feature branches: Use descriptive names like `story/issue-132/create-adr-for-user-registration`

### Documentation Deployment

Documentation is automatically built and deployed to GitHub Pages when changes are pushed to `main` in the `docs/` directory.

## Architecture Notes

**Current State:** This repository is in the planning/design phase. The R&D process documents user journeys and architectural decisions before implementation begins.

**GitHub OAuth Authentication:** The platform plans to use GitHub OAuth for user registration and authentication (see User Journey 0001).

**Target Audience:** Developers and teams who want to deploy and compete with autonomous battlebots.

**Technology Stack Indicators:**
- Go as primary language (from devcontainer)
- Protocol Buffers for service definitions (protoc in devcontainer)
- Hugo + Docsy for documentation
- GitHub OAuth for authentication (planned)
