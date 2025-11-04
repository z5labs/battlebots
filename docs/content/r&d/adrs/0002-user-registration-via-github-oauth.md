---
title: "[0002] User Registration via GitHub OAuth"
description: >
    How users will register and authenticate with the Battle Bots platform using GitHub OAuth and stateless JWT tokens
type: docs
weight: 0002
category: "user-journey"
status: "accepted"
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

Chosen option: "GitHub OAuth authentication with stateless JWT tokens", because it best meets our decision drivers:

* **User experience**: Developers already have GitHub accounts - minimal registration friction
* **Security**: Leverages GitHub's OAuth 2.0 without password management burden
* **Implementation complexity**: Single OAuth provider reduces development time
* **Timeline**: Fastest path to launch with proven technology
* **Trust**: GitHub is the natural identity provider for our developer audience
* **Scalability**: Stateless JWT tokens enable horizontal scaling without session synchronization

The implementation uses GitHub OAuth for initial authentication, then converts the GitHub access token to internal JWT tokens for stateless API authentication. This hybrid approach provides OAuth convenience with JWT scalability.

### Consequences

* Good, because no server-side session state enables horizontal scalability
* Good, because JWT tokens reduce database lookups (user info in token claims)
* Good, because stateless architecture simplifies microservices integration
* Good, because developers trust GitHub as identity provider
* Good, because refresh token rotation provides security without UX friction
* Neutral, because requires implementing JWT token service (RS256 signing/validation)
* Neutral, because GitHub OAuth is OAuth 2.0, not OIDC (must generate our own ID tokens)
* Bad, because vendor dependency on GitHub for initial authentication
* Bad, because requires token blacklist for immediate revocation (adds some state)
* Bad, because limits to users with GitHub accounts (acceptable for developer audience)

### Confirmation

Implementation compliance will be confirmed through:

1. **Security Testing**: Penetration testing validates XSS/CSRF protection, token validation, and PKCE implementation
2. **Integration Tests**: Automated tests verify complete OAuth flow and JWT token generation/validation
3. **Code Review**: Security-focused review of JWT signing, token storage, and refresh token rotation
4. **Load Testing**: Horizontal scalability validated with 10,000+ concurrent users across multiple servers
5. **Documentation Review**: Architecture diagrams and sequence diagrams accurately reflect stateless JWT implementation

<!-- This is an optional element. Feel free to remove. -->
## Pros and Cons of the Options

### GitHub OAuth authentication with OpenFGA authorization

Single OAuth provider (GitHub) for registration and authentication, combined with OpenFGA for fine-grained authorization.

**Authentication vs Authorization:**
- **Authentication (GitHub OAuth + JWT):** Verifies user identity - "Who are you?"
- **Authorization (OpenFGA):** Determines permissions - "What can you do?"

This option separates authentication concerns (handled by GitHub OAuth and JWT tokens) from authorization concerns (handled by OpenFGA). Users authenticate once via GitHub, receive JWT tokens for API access, and every protected operation checks permissions through OpenFGA based on relationships between users and resources (organizations, teams, bots, games).

* Good, because target audience (developers) already have GitHub accounts
* Good, because no password management or reset flows needed
* Good, because GitHub's OAuth is well-documented and reliable
* Good, because reduces implementation complexity and time to launch
* Good, because GitHub identity ties naturally to developer workflows
* Good, because OpenFGA provides fine-grained authorization (object-level permissions, not just roles)
* Good, because OpenFGA naturally models organizational hierarchy (organization → team → bot → game)
* Good, because OpenFGA scales to millions of authorization relationships and checks
* Good, because separation of authentication (GitHub OAuth) and authorization (OpenFGA) follows security best practices
* Good, because self-hosted OpenFGA has no per-user costs (important for platform growth)
* Good, because OpenFGA's ReBAC model supports complex permission inheritance and delegation
* Neutral, because limits to users with GitHub accounts (acceptable for developer audience)
* Neutral, because requires deploying and maintaining OpenFGA service infrastructure
* Neutral, because requires learning ReBAC authorization modeling concepts
* Bad, because vendor dependency on GitHub for authentication
* Bad, because no fallback if GitHub OAuth is unavailable
* Bad, because additional infrastructure complexity (OpenFGA service + database for authorization data)
* Bad, because network latency added for OpenFGA authorization checks on every protected operation

#### Implementation Visualization

**Architecture Diagram:**

```mermaid
graph LR
    User[User Browser] --> WebApp[Battle Bots Web App]
    WebApp --> GitHub[GitHub OAuth]
    WebApp --> DB[(Application Database)]
    WebApp --> JWT[JWT Token Service]
    WebApp --> Redis[(Redis - Token Blacklist)]
    WebApp --> AuthZ[OpenFGA Service]
    AuthZ --> FGA_DB[(OpenFGA Database)]
    GitHub --> User

    style WebApp fill:#e1f5ff
    style GitHub fill:#f0f0f0
    style DB fill:#fff4e1
    style JWT fill:#f3e5f5
    style Redis fill:#ffe0e0
    style AuthZ fill:#e8f5e9
    style FGA_DB fill:#c8e6c9
```

**Architecture Notes:**
- **Authentication Layer**: GitHub OAuth validates user identity, JWT Token Service generates stateless tokens
- **Authorization Layer**: OpenFGA Service handles fine-grained permissions checks for all protected operations
- **Stateless Design**: JWT tokens eliminate session storage; Redis is optional for token revocation blacklist
- **Separation of Concerns**: Authentication database (user accounts) separate from authorization database (relationships/permissions)

**REST API Endpoints:**

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `GET` | `/auth/github/login` | No | Initiates GitHub OAuth flow with PKCE by generating code_challenge and CSRF state token, redirecting to GitHub authorization page |
| `GET` | `/auth/github/callback` | No | Handles OAuth callback from GitHub, exchanges auth code for GitHub access token, fetches user profile, creates/updates account, **generates internal JWT access token + refresh token**, sets httpOnly cookies |
| `POST` | `/auth/terms/accept` | JWT (Cookie) | Accepts terms of service for new user accounts (called before account creation) |
| `POST` | `/auth/refresh` | Refresh Token (Cookie) | Exchanges valid refresh token for new JWT access token + new refresh token (rotation), updates httpOnly cookies |
| `GET` | `/auth/session` | JWT (Cookie) | Returns current authenticated user information from JWT claims (no database lookup) |
| `POST` | `/auth/logout` | JWT (Cookie) | Revokes refresh token in database, optionally blacklists JWT, clears authentication cookies |

**Protected Resource Endpoints (with OpenFGA Authorization):**

Protected endpoints follow a two-step security model: JWT authentication validates identity, then OpenFGA authorization checks permissions.

| Method | Endpoint | Authentication | Authorization Check | Purpose |
|--------|----------|----------------|---------------------|---------|
| `POST` | `/bots/:id/deploy` | JWT (Cookie) | OpenFGA: `user:github\|username` has `can_deploy` on `bot:id` | Deploy bot to game server |
| `PUT` | `/bots/:id` | JWT (Cookie) | OpenFGA: `user:github\|username` has `can_edit` on `bot:id` | Update bot configuration |
| `DELETE` | `/bots/:id` | JWT (Cookie) | OpenFGA: `user:github\|username` has `can_delete` on `bot:id` | Delete bot |
| `POST` | `/teams/:id/members` | JWT (Cookie) | OpenFGA: `user:github\|username` has `can_manage_members` on `team:id` | Add member to team |
| `POST` | `/organizations/:id/teams` | JWT (Cookie) | OpenFGA: `user:github\|username` has `can_create_teams` on `organization:id` | Create team in organization |

**Authorization Flow Pattern:**
1. Middleware validates JWT and extracts user identity (e.g., `github|username`)
2. Middleware maps HTTP method + route to OpenFGA relation (e.g., `POST /bots/:id/deploy` → `can_deploy`)
3. Middleware calls OpenFGA Check API: `check(user, relation, object)`
4. If OpenFGA returns `allowed: true`, request proceeds to handler
5. If OpenFGA returns `allowed: false`, return `403 Forbidden`

**Sequence Diagram - Registration/Login Flow:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant GitHub as GitHub OAuth
    participant DB as Database
    participant JWT as JWT Token Service

    User->>WebApp: GET /auth/github/login
    WebApp->>WebApp: Generate code_verifier + code_challenge (PKCE)
    WebApp->>WebApp: Generate CSRF state token
    WebApp->>User: 302 Redirect to GitHub OAuth (with code_challenge)
    User->>GitHub: Authorize Battle Bots application

    alt Authorization Successful
        GitHub->>WebApp: GET /auth/github/callback?code=xxx&state=yyy
        WebApp->>WebApp: Validate state token (CSRF protection)
        WebApp->>GitHub: POST /login/oauth/access_token<br/>(exchange code + code_verifier)
        GitHub-->>WebApp: Return GitHub access token
        WebApp->>GitHub: GET /user (fetch profile)
        GitHub-->>WebApp: Return user data (ID, username, email)

        alt User Exists
            WebApp->>DB: Update user profile
            DB-->>WebApp: Profile updated
        else New User
            WebApp->>User: Show Terms of Service page
            User->>WebApp: POST /auth/terms/accept (with JWT cookie)
            WebApp->>DB: Create user account
            DB-->>WebApp: Account created
        end

        WebApp->>JWT: Generate JWT access token (RS256, 15min expiry)
        JWT-->>WebApp: Signed JWT with claims (sub, iss, aud, exp, github_id, ...)
        WebApp->>WebApp: Generate refresh token (random 32-byte)
        WebApp->>DB: Store refresh token hash (user_id, token_hash, expires_at)
        DB-->>WebApp: Token stored
        WebApp->>User: 302 Redirect to dashboard<br/>(Set httpOnly cookies: access_token, refresh_token, csrf_token)

    else Authorization Failed/Cancelled
        GitHub->>WebApp: GET /auth/github/callback?error=xxx
        WebApp->>User: Show error message
    end
```

**Sequence Diagram - Protected Operation with OpenFGA Authorization:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant AuthMiddleware as Auth Middleware
    participant AuthZMiddleware as AuthZ Middleware
    participant JWT as JWT Token Service
    participant OpenFGA as OpenFGA Service
    participant Handler as Request Handler
    participant DB as Database

    Note over User,DB: User already authenticated with JWT cookie

    User->>WebApp: POST /bots/battle-bot-1/deploy<br/>(JWT in httpOnly cookie)
    WebApp->>AuthMiddleware: Validate authentication
    AuthMiddleware->>JWT: Verify JWT signature & expiry
    JWT-->>AuthMiddleware: JWT valid, claims extracted

    alt JWT Invalid or Expired
        AuthMiddleware-->>WebApp: 401 Unauthorized
        WebApp->>User: Redirect to login or show error
    else JWT Valid
        AuthMiddleware->>AuthZMiddleware: user_id: github|alice, operation: deploy, resource: bot:battle-bot-1
        AuthZMiddleware->>OpenFGA: Check(user: "user:github|alice",<br/>relation: "can_deploy",<br/>object: "bot:battle-bot-1")
        OpenFGA->>OpenFGA: Evaluate authorization model<br/>(check relationships and rules)

        alt User Not Authorized
            OpenFGA-->>AuthZMiddleware: allowed: false
            AuthZMiddleware-->>WebApp: 403 Forbidden
            WebApp->>User: Show "You don't have permission" error
        else User Authorized
            OpenFGA-->>AuthZMiddleware: allowed: true
            AuthZMiddleware->>Handler: Forward request
            Handler->>DB: Execute deployment logic
            DB-->>Handler: Deployment successful
            Handler-->>WebApp: 200 OK {deployment_status}
            WebApp->>User: Show success message
        end
    end
```

**Authorization Decision Flow:**

The authorization middleware performs the following steps for every protected endpoint:

1. **Extract User Identity**: Get user ID from validated JWT claims (e.g., `github|alice`)
2. **Map Operation to Relation**: Translate HTTP method + endpoint to OpenFGA relation
   - `POST /bots/:id/deploy` → `can_deploy`
   - `PUT /bots/:id` → `can_edit`
   - `DELETE /bots/:id` → `can_delete`
3. **Construct OpenFGA Check**: Build check request with user, relation, and object
4. **Evaluate Permission**: OpenFGA traverses relationship graph based on authorization model
5. **Enforce Decision**: Allow (200/201) or deny (403) based on OpenFGA response

**OpenFGA Authorization Model (Starter):**

This is the initial authorization model for Battle Bots MVP. It defines the organizational hierarchy and permission inheritance patterns.

```
model
  schema 1.1

type user

type organization
  relations
    define owner: [user]
    define member: [user]
    define can_create_teams: owner or member
    define can_invite_members: owner

type team
  relations
    define parent_org: [organization]
    define admin: [user]
    define member: [user] or admin or owner from parent_org
    define can_manage_bots: admin or owner from parent_org
    define can_manage_members: admin or owner from parent_org

type bot
  relations
    define parent_team: [team]
    define creator: [user]
    define can_view: member from parent_team
    define can_edit: creator or can_manage_bots from parent_team
    define can_deploy: creator or can_manage_bots from parent_team
    define can_delete: creator or admin from parent_team

type game
  relations
    define participant_bot: [bot]
    define creator: [user]
    define can_view: creator or can_view from participant_bot
    define can_cancel: creator
```

**Model Explanation:**

- **Organizations**: Top-level entities with owners and members. Owners can invite members; both can create teams.
- **Teams**: Belong to organizations. Organization owners automatically have admin rights on all teams in their org. Team admins can manage bots and members.
- **Bots**: Belong to teams. Bot creators can edit/deploy/delete their bots. Team admins can also manage all bots in their team. All team members can view bots.
- **Games**: Include participating bots. Game creator and anyone who can view participating bots can view the game.

**Permission Inheritance Examples:**
- Alice is owner of Org A → Alice is automatically admin of Team X (part of Org A) → Alice can manage all bots in Team X
- Bob is creator of Bot Y in Team X → Bob can edit, deploy, and delete Bot Y
- Carol is member of Team X → Carol can view all bots in Team X, but cannot edit/deploy them

**Relationship Tuple Examples:**
```
# Organization membership
user:github|alice, owner, organization:acme-corp
user:github|bob, member, organization:acme-corp

# Team structure
organization:acme-corp, parent_org, team:alpha-team
user:github|carol, admin, team:alpha-team

# Bot ownership
team:alpha-team, parent_team, bot:battle-bot-1
user:github|bob, creator, bot:battle-bot-1

# Game participation
bot:battle-bot-1, participant_bot, game:match-123
user:github|alice, creator, game:match-123
```

For complete OpenFGA deployment and configuration details, see [OpenFGA Analysis](/r&d/analysis/open-source-applications/openfga-analysis/).

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

#### Implementation Visualization

**Architecture Diagram:**

```mermaid
graph LR
    User[User Browser] --> WebApp[Battle Bots Web App]
    WebApp --> AuthService[Auth Service]
    AuthService --> DB[(Database)]
    AuthService --> JWT[JWT Token Service]
    AuthService --> Email[Email Service]
    Email --> SMTP[SMTP Server]

    style WebApp fill:#e1f5ff
    style AuthService fill:#e8f5e9
    style DB fill:#fff4e1
    style JWT fill:#f3e5f5
    style Email fill:#fff3e0
    style SMTP fill:#f0f0f0
```

**REST API Endpoints:**

*Authentication Endpoints:*

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `POST` | `/auth/register` | No | Register new user account with email, password, and username. Sends verification email. |
| `POST` | `/auth/login` | No | Authenticate with email/password credentials. Returns JWT access token (15min) and refresh token (7 day). |
| `POST` | `/auth/logout` | JWT | Invalidates refresh token and terminates user session. |
| `POST` | `/auth/refresh` | Refresh Token | Exchanges valid refresh token for new access token. |

*Email Verification Endpoints:*

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `GET` | `/auth/verify-email` | No | Verifies email address using token from verification email (query param: `token`). |
| `POST` | `/auth/resend-verification` | No | Resends verification email to user's registered email address. |

*Password Management Endpoints:*

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `POST` | `/auth/forgot-password` | No | Requests password reset. Sends reset email if account exists (always returns success to prevent enumeration). |
| `POST` | `/auth/reset-password` | No | Resets password using token from reset email. Requires `token` and `new_password` in request body. |
| `POST` | `/auth/change-password` | JWT | Changes password for authenticated user. Requires `current_password` and `new_password`. |

**Sequence Diagram - Registration Flow:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Auth as Auth Service
    participant DB as Database
    participant Email as Email Service

    User->>WebApp: POST /auth/register<br/>{email, password, username}
    WebApp->>Auth: Validate input (email format, password strength)

    alt Validation Failed
        Auth-->>WebApp: 400 Bad Request {errors}
        WebApp->>User: Show validation error messages
    else Validation Passed
        Auth->>DB: Check if email exists
        DB-->>Auth: Email availability

        alt Email Already Exists
            Auth-->>WebApp: 409 Conflict
            WebApp->>User: Show "Email already registered" error
        else Email Available
            Auth->>Auth: Hash password (bcrypt/argon2)
            Auth->>DB: Create user account (unverified)
            DB-->>Auth: Account created
            Auth->>Auth: Generate email verification token
            Auth->>DB: Store verification token (with expiry)
            Auth->>Email: Send verification email with link:<br/>GET /auth/verify-email?token=xxx
            Email-->>User: Verification email delivered
            Auth-->>WebApp: 201 Created
            WebApp->>User: Show "Check your email" message
        end
    end

    Note over User,Email: User clicks verification link in email
    User->>WebApp: GET /auth/verify-email?token=xxx
    WebApp->>Auth: Validate verification token
    Auth->>DB: Mark email as verified
    WebApp->>User: 200 OK - Email verified, redirect to login
```

**Sequence Diagram - Login Flow:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Auth as Auth Service
    participant DB as Database
    participant JWT as JWT Service

    User->>WebApp: POST /auth/login<br/>{email, password}
    WebApp->>Auth: Authenticate credentials
    Auth->>DB: Fetch user by email
    DB-->>Auth: User record

    alt User Not Found
        Auth-->>WebApp: 401 Unauthorized
        WebApp->>User: Show "Invalid credentials" error
    else User Found
        Auth->>Auth: Verify password hash (bcrypt/argon2)

        alt Password Invalid
            Auth-->>WebApp: 401 Unauthorized
            WebApp->>User: Show "Invalid credentials" error
        else Password Valid
            alt Email Not Verified
                Auth-->>WebApp: 403 Forbidden {reason: "email_not_verified"}
                WebApp->>User: Show verification reminder +<br/>option to POST /auth/resend-verification
            else Email Verified
                Auth->>JWT: Generate access token (15min expiry)
                JWT-->>Auth: Access token (JWT)
                Auth->>JWT: Generate refresh token (7 day expiry)
                JWT-->>Auth: Refresh token
                Auth->>DB: Store refresh token hash
                Auth-->>WebApp: 200 OK {access_token, refresh_token, expires_in}
                WebApp->>User: Store tokens, redirect to dashboard
            end
        end
    end
```

**Sequence Diagram - Password Reset Flow:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Auth as Auth Service
    participant DB as Database
    participant Email as Email Service

    Note over User,Email: Request Password Reset
    User->>WebApp: Click "Forgot Password" link
    User->>WebApp: POST /auth/forgot-password<br/>{email}
    WebApp->>Auth: Request password reset
    Auth->>DB: Check if email exists

    Note over Auth,Email: Always return success (prevent email enumeration)
    Auth-->>WebApp: 200 OK {message: "If account exists, email sent"}
    WebApp->>User: Show "Check your email" message

    opt Email Exists
        Auth->>Auth: Generate reset token (1 hour expiry)
        Auth->>DB: Store reset token hash
        Auth->>Email: Send reset email with link:<br/>GET /auth/reset-password?token=xxx
        Email-->>User: Reset email delivered
    end

    Note over User,DB: Reset Password
    User->>WebApp: GET /auth/reset-password?token=xxx<br/>(click link in email)
    WebApp->>Auth: Validate reset token
    Auth->>DB: Check token validity & expiry

    alt Token Invalid/Expired
        Auth-->>WebApp: 400 Bad Request {error: "invalid_token"}
        WebApp->>User: Show error, offer POST /auth/forgot-password
    else Token Valid
        Auth-->>WebApp: 200 OK
        WebApp->>User: Show new password form
        User->>WebApp: POST /auth/reset-password<br/>{token, new_password}
        WebApp->>Auth: Update password
        Auth->>Auth: Hash new password (bcrypt/argon2)
        Auth->>DB: Update password hash
        Auth->>DB: Invalidate reset token
        Auth->>Email: Send password changed notification email
        Auth-->>WebApp: 200 OK
        WebApp->>User: Show success, redirect to POST /auth/login
    end
```

### Google OAuth authentication

Single OAuth provider (Google) for registration and authentication.

* Good, because most users have Google accounts
* Good, because no password management needed
* Good, because Google OAuth is reliable and well-documented
* Neutral, because less aligned with developer-focused audience than GitHub
* Bad, because vendor dependency on Google
* Bad, because no fallback if Google OAuth is unavailable

#### Implementation Visualization

**Architecture Diagram:**

```mermaid
graph LR
    User[User Browser] --> WebApp[Battle Bots Web App]
    WebApp --> Google[Google OAuth]
    WebApp --> DB[(Database)]
    WebApp --> Session[Session Store]
    Google --> User

    style WebApp fill:#e1f5ff
    style Google fill:#f0f0f0
    style DB fill:#fff4e1
    style Session fill:#fff4e1
```

**REST API Endpoints:**

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `GET` | `/auth/google/login` | No | Initiates Google OAuth flow by generating CSRF state token and redirecting to Google authorization page |
| `GET` | `/auth/google/callback` | No | Handles OAuth callback from Google, exchanges auth code for access token, validates ID token, fetches user profile, creates/updates account |
| `POST` | `/auth/terms/accept` | Session | Accepts terms of service for new user accounts (called before account creation) |
| `GET` | `/auth/session` | Session | Returns current authenticated user information and session status |
| `POST` | `/auth/logout` | Session | Terminates user session and clears authentication cookies/tokens |

**Sequence Diagram - Registration/Login Flow:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Google as Google OAuth
    participant DB as Database
    participant Session as Session Store

    User->>WebApp: GET /auth/google/login
    WebApp->>WebApp: Generate CSRF state token
    WebApp->>User: 302 Redirect to Google OAuth
    User->>Google: Authorize Battle Bots application

    alt Authorization Successful
        Google->>WebApp: GET /auth/google/callback?code=xxx&state=yyy
        WebApp->>WebApp: Validate state token
        WebApp->>Google: POST /token (exchange code)
        Google-->>WebApp: Return access_token & id_token (JWT)
        WebApp->>WebApp: Validate ID token (JWT signature + claims)
        WebApp->>Google: GET /oauth2/v2/userinfo (optional)
        Google-->>WebApp: Return user data (sub, email, name, picture)

        alt User Exists
            WebApp->>DB: Update user profile
            DB-->>WebApp: Profile updated
        else New User
            WebApp->>User: Show Terms of Service page
            User->>WebApp: POST /auth/terms/accept
            WebApp->>DB: Create user account
            DB-->>WebApp: Account created
        end

        WebApp->>Session: Create authenticated session
        Session-->>WebApp: Session token
        WebApp->>User: 302 Redirect to dashboard (with session cookie)

    else Authorization Failed/Cancelled
        Google->>WebApp: GET /auth/google/callback?error=xxx
        WebApp->>User: Show error message
    end
```

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

#### Implementation Visualization

**Architecture Diagram:**

```mermaid
graph TB
    User[User Browser] --> WebApp[Battle Bots Web App]
    WebApp --> OAuthGateway[OAuth Gateway/Strategy]
    OAuthGateway --> GitHub[GitHub OAuth Provider]
    OAuthGateway --> Google[Google OAuth Provider]
    OAuthGateway --> GitLab[GitLab OAuth Provider]
    WebApp --> DB[(Database)]
    WebApp --> Session[Session Store]
    DB --> UserAccounts[User Accounts Table]
    DB --> LinkedProviders[Linked Providers Table]

    style WebApp fill:#e1f5ff
    style OAuthGateway fill:#e8f5e9
    style GitHub fill:#f0f0f0
    style Google fill:#f0f0f0
    style GitLab fill:#f0f0f0
    style DB fill:#fff4e1
    style Session fill:#fff4e1
    style UserAccounts fill:#fff9c4
    style LinkedProviders fill:#fff9c4
```

**REST API Endpoints:**

*OAuth Authentication Endpoints (Provider-Agnostic):*

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `GET` | `/auth/:provider/login` | No | Initiates OAuth flow for specified provider (`:provider` = `github`, `google`, or `gitlab`). Generates CSRF state and redirects. |
| `GET` | `/auth/:provider/callback` | No | Handles OAuth callback from specified provider. Exchanges code for token, fetches profile, handles account creation/linking logic. |
| `POST` | `/auth/terms/accept` | Session | Accepts terms of service for new user accounts (called before account creation). |
| `GET` | `/auth/session` | Session | Returns current authenticated user information, session status, and list of linked providers. |
| `POST` | `/auth/logout` | Session | Terminates user session and clears authentication cookies/tokens. |

*Provider Management Endpoints:*

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| `GET` | `/auth/providers` | Session (optional) | Lists available OAuth providers. If authenticated, includes which providers are linked to current user's account. |
| `POST` | `/auth/link/:provider` | Session | Initiates OAuth flow to link an additional provider to the authenticated user's existing account. |
| `DELETE` | `/auth/unlink/:provider` | Session | Unlinks specified provider from authenticated user's account. Requires at least one provider to remain linked. |

**Sequence Diagram - Provider Selection and Registration:**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Gateway as OAuth Gateway
    participant Provider as Selected OAuth Provider
    participant DB as Database
    participant Session as Session Store

    User->>WebApp: GET /login (view login page)
    WebApp->>WebApp: GET /auth/providers (fetch available)
    WebApp->>User: Show provider buttons (GitHub, Google, GitLab)
    User->>WebApp: GET /auth/github/login (example: select GitHub)
    WebApp->>Gateway: Initiate OAuth with GitHub
    Gateway->>Gateway: Load provider-specific config & strategy
    Gateway->>WebApp: Generate CSRF state token
    WebApp->>User: 302 Redirect to github.com/login/oauth/authorize
    User->>Provider: Authorize Battle Bots application

    alt Authorization Successful
        Provider->>Gateway: GET /auth/github/callback?code=xxx&state=yyy
        Gateway->>Gateway: Validate state token
        Gateway->>Provider: Exchange code for access token
        Provider-->>Gateway: Return access token
        Gateway->>Provider: Fetch user profile
        Provider-->>Gateway: Return user data (provider_id, email, name)
        Gateway->>DB: Check if provider user ID exists

        alt Provider Account Already Linked
            DB-->>Gateway: Found existing account
            Gateway->>DB: Update profile from provider
            Gateway->>Session: Create authenticated session
            Session-->>Gateway: Session token
            Gateway->>WebApp: 200 OK
            WebApp->>User: 302 Redirect to dashboard
        else New Provider Account
            Gateway->>DB: Check email across all providers
            DB-->>Gateway: Email lookup results

            alt Email Exists with Different Provider
                DB-->>Gateway: Found account with same email
                Gateway->>WebApp: Show account linking confirmation UI
                WebApp->>User: "Link to existing account or create new?"

                opt User Chooses to Link
                    User->>WebApp: POST /auth/link/github {confirm: true}
                    WebApp->>Gateway: Link provider to existing account
                    Gateway->>DB: INSERT INTO linked_providers
                    Gateway->>Session: Create authenticated session
                    WebApp->>User: 302 Redirect to dashboard
                end

                opt User Chooses New Account
                    User->>WebApp: Create separate account
                    WebApp->>User: Show Terms of Service page
                    User->>WebApp: POST /auth/terms/accept
                    Gateway->>DB: INSERT INTO users (new account)
                    Gateway->>Session: Create authenticated session
                    WebApp->>User: 302 Redirect to onboarding
                end

            else Email Not Found
                WebApp->>User: Show Terms of Service page
                User->>WebApp: POST /auth/terms/accept
                Gateway->>DB: INSERT INTO users + linked_providers
                Gateway->>Session: Create authenticated session
                WebApp->>User: 302 Redirect to onboarding
            end
        end

    else Authorization Failed/Cancelled
        Provider->>Gateway: GET /auth/github/callback?error=access_denied
        Gateway->>WebApp: OAuth failed
        WebApp->>User: Show error, return to GET /login
    end
```

**Sequence Diagram - Account Linking (Authenticated User):**

```mermaid
sequenceDiagram
    actor User
    participant WebApp as Battle Bots Web App
    participant Gateway as OAuth Gateway
    participant NewProvider as New OAuth Provider
    participant DB as Database

    Note over User,DB: User already authenticated, wants to link additional provider

    User->>WebApp: GET /account/settings
    WebApp->>WebApp: GET /auth/providers (with auth token)
    WebApp->>User: Show linked providers + "Link GitHub" button
    User->>WebApp: POST /auth/link/github (click "Link GitHub")
    WebApp->>Gateway: Initiate OAuth for linking
    Gateway->>Gateway: Generate state token with user_id embedded
    Gateway->>User: 302 Redirect to github.com/login/oauth/authorize
    User->>NewProvider: Authorize Battle Bots application

    alt Authorization Successful
        NewProvider->>Gateway: GET /auth/github/callback?code=xxx&state=yyy
        Gateway->>Gateway: Validate state & extract user_id from state
        Gateway->>NewProvider: POST /login/oauth/access_token (exchange code)
        NewProvider-->>Gateway: Return access token
        Gateway->>NewProvider: GET /user (fetch profile)
        NewProvider-->>Gateway: Return provider user data (github_id, email)
        Gateway->>DB: SELECT * FROM linked_providers<br/>WHERE provider='github'<br/>AND provider_user_id=github_id

        alt Provider Account Already Linked to Different User
            DB-->>Gateway: Conflict - belongs to user_id=999
            Gateway->>WebApp: 409 Conflict {error: "provider_already_linked"}
            WebApp->>User: Show "This GitHub account is already linked<br/>to a different Battle Bots account"
        else Provider Account Available
            Gateway->>DB: INSERT INTO linked_providers<br/>(user_id, provider, provider_user_id)
            DB-->>Gateway: Link created successfully
            Gateway->>WebApp: 200 OK {linked_providers: [...]}
            WebApp->>User: Show success + refresh GET /auth/providers
        end

    else Authorization Failed/Cancelled
        NewProvider->>Gateway: GET /auth/github/callback?error=access_denied
        Gateway->>WebApp: 400 Bad Request {error: "oauth_cancelled"}
        WebApp->>User: Show "Linking cancelled" message
    end
```

**Key Implementation Considerations:**

- **Provider Strategy Pattern**: Each OAuth provider (GitHub, Google, GitLab) has its own strategy implementing a common interface for authorization, token exchange, and profile fetching
- **Account Linking Logic**: Must handle cases where users authenticate with different providers but share the same email
- **Data Model**: Requires `linked_providers` table with columns: `user_id`, `provider`, `provider_user_id`, `provider_username`, `linked_at`
- **Testing Complexity**: Need integration tests for each provider plus account linking scenarios
- **User Experience**: Clear messaging when email conflicts occur across providers

<!-- This is an optional element. Feel free to remove. -->
## More Information

Related to User Journey 0001 (User Registration and Authentication).

This decision focuses on the initial launch strategy. Future ADRs may address:
- Adding additional OAuth providers based on user feedback
- Account migration strategies if switching providers
- Service account or API key authentication for bot deployments
