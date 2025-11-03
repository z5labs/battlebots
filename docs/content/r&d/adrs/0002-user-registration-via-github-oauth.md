---
title: "[0002] User Registration via GitHub OAuth"
description: >
    How users will register and authenticate with the Battle Bots platform
type: docs
weight: 0002
category: "user-journey"
status: "proposed"
date: 2025-11-02
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions (frameworks, auth strategies, cross-cutting patterns)
- user-journey: Solutions for specific user journey problems (feature implementation approaches)
- api-design: API endpoint design decisions (pagination, filtering, bulk operations)
-->

## Context and Problem Statement

Users need a way to register and authenticate with the Battle Bots platform to create and manage their autonomous bots. The registration process should be secure, user-friendly, and minimize friction for developers who are our target audience.

How should we implement user registration and authentication for the Battle Bots platform?

## Decision Drivers

* User experience: Minimize registration friction for developer audience
* Security: Ensure secure authentication without managing passwords
* Implementation complexity: Reduce development and maintenance burden
* Timeline: Need to launch quickly with minimal authentication infrastructure
* Trust: Leverage existing identity providers that developers already use
* Bot deployment: Need to tie bot ownership to verified user accounts

## Considered Options

* GitHub OAuth authentication
* Email/password registration with JWT
* Google OAuth authentication
* Support multiple OAuth providers (GitHub, Google, GitLab)

## Decision Outcome

Chosen option: "[option 1]", because [justification. e.g., only option, which meets k.o. criterion decision driver | which resolves force force | ... | comes out best (see below)].

<!-- This is an optional element. Feel free to remove. -->
### Consequences

* Good, because [positive consequence, e.g., improvement of one or more desired qualities, ...]
* Bad, because [negative consequence, e.g., compromising one or more desired qualities, ...]
* ...

<!-- This is an optional element. Feel free to remove. -->
### Confirmation

[Describe how the implementation of/compliance with the ADR is confirmed. E.g., by a review or an ArchUnit test.
 Although we classify this element as optional, it is included in most ADRs.]

<!-- This is an optional element. Feel free to remove. -->
## Pros and Cons of the Options

### GitHub OAuth authentication

Single OAuth provider (GitHub) for registration and authentication.

* Good, because target audience (developers) already have GitHub accounts
* Good, because no password management or reset flows needed
* Good, because GitHub's OAuth is well-documented and reliable
* Good, because reduces implementation complexity and time to launch
* Good, because GitHub identity ties naturally to developer workflows
* Neutral, because limits to users with GitHub accounts (acceptable for developer audience)
* Bad, because vendor dependency on GitHub
* Bad, because no fallback if GitHub OAuth is unavailable

### Email/password registration with JWT

Traditional registration with email/password and JWT tokens.

* Good, because no dependency on external OAuth providers
* Good, because works for users without third-party accounts
* Good, because full control over authentication flow
* Bad, because requires implementing password reset, email verification
* Bad, because need to securely store and hash passwords
* Bad, because higher implementation and maintenance complexity
* Bad, because more friction in registration process
* Bad, because security burden of password management

### Google OAuth authentication

Single OAuth provider (Google) for registration and authentication.

* Good, because most users have Google accounts
* Good, because no password management needed
* Good, because Google OAuth is reliable and well-documented
* Neutral, because less aligned with developer-focused audience than GitHub
* Bad, because vendor dependency on Google
* Bad, because no fallback if Google OAuth is unavailable

### Support multiple OAuth providers (GitHub, Google, GitLab)

Allow users to choose from multiple OAuth providers.

* Good, because provides user choice and flexibility
* Good, because reduces single vendor dependency
* Good, because accommodates different user preferences
* Bad, because significantly higher implementation complexity
* Bad, because need to handle account linking/merging
* Bad, because increases testing surface area
* Bad, because delays time to launch
* Bad, because more complex user experience (choice paralysis)

<!-- This is an optional element. Feel free to remove. -->
## More Information

Related to User Journey 0001 (User Registration and Authentication).

This decision focuses on the initial launch strategy. Future ADRs may address:
- Adding additional OAuth providers based on user feedback
- Account migration strategies if switching providers
- Service account or API key authentication for bot deployments
