---
title: OpenFGA Analysis
description: >
    Comprehensive analysis of OpenFGA fine-grained authorization system, including JWT integration, deployment options, policy definition, and dynamic data capabilities.
type: docs
weight: 1
date: 2025-11-02
---

## Executive Summary

This document analyzes OpenFGA (Open Fine-Grained Authorization), a high-performance authorization system inspired by Google Zanzibar, to determine its suitability for implementing fine-grained access control in the Battle Bots platform. OpenFGA provides Relationship-Based Access Control (ReBAC) that answers authorization questions like "Does user U have relation R with object O?" rather than traditional role-based approaches.

**Key Finding:** OpenFGA is a mature, production-ready authorization system that separates authentication from authorization, allowing it to work seamlessly with JWT-based authentication systems. It offers flexible deployment options (Docker, Kubernetes, binary), uses an intuitive DSL for policy definition, and supports dynamic authorization decisions through contextual tuples. For Battle Bots' planned GitHub OAuth authentication strategy, OpenFGA provides a robust path to implementing fine-grained permissions for bot management, organization access, and game resources.

## Overview of OpenFGA

### What is OpenFGA?

OpenFGA is a high-performance, flexible authorization/permission engine built for developers and inspired by Google's Zanzibar paper. It implements Relationship-Based Access Control (ReBAC), which models permissions as relationships between users and resources.

**Core Characteristics:**
- Open source (Apache 2.0 license)
- Inspired by Google Zanzibar (used by Google Drive, Calendar, Photos, etc.)
- Focuses exclusively on authorization (not authentication)
- Designed for fine-grained permissions at scale
- Supports both direct relationships and computed permissions
- Provides APIs for authorization checks, relationship management, and model queries

### Key Concepts

**Authorization Model:** Defines the types of objects in your system and how they relate to each other (e.g., users can be viewers, editors, or owners of documents).

**Relationship Tuples:** Store the actual relationships between users and objects (e.g., "anne is owner of document:1").

**Check API:** Evaluates whether a user has a specific relation/permission on an object based on the model and stored tuples.

**Contextual Tuples:** Temporary relationships provided at query time that exist only for that specific authorization check.

### ReBAC vs. RBAC

| Aspect | Traditional RBAC | OpenFGA ReBAC |
|--------|------------------|---------------|
| **Permission Model** | Users assigned to roles; roles have permissions | Users have relationships with specific objects |
| **Granularity** | Coarse-grained (role-level) | Fine-grained (object-level) |
| **Hierarchy** | Simple role hierarchies | Complex relationship graphs with inheritance |
| **Question Answered** | "Does user have role X?" | "Does user have relation R with object O?" |
| **Examples** | "Is user an admin?" | "Can user edit document:123?" |
| **Scalability** | Limited for complex scenarios | Designed for millions of relationships |

## JWT Verification and Integration

### Separation of Authentication and Authorization

OpenFGA focuses exclusively on authorization, making it designed to work alongside authentication systems like JWT verification rather than replacing them. The typical integration pattern separates concerns:

**Authentication (AuthN):** Verifies user identity through JWT token validation
**Authorization (AuthZ):** Determines what the authenticated user can do via OpenFGA

### Integration Pattern

The standard integration flow for JWT + OpenFGA:

```
1. Request arrives with JWT Bearer token
   ↓
2. Middleware validates JWT signature
   ↓
3. Middleware extracts user identity from JWT claims
   ↓
4. Application calls OpenFGA Check API
   - User: extracted from JWT (e.g., "user:anne")
   - Relation: mapped from operation (e.g., "can_edit")
   - Object: target resource (e.g., "document:123")
   ↓
5. OpenFGA evaluates authorization
   ↓
6. Application allows or denies request
```

### Framework-Specific Implementations

#### Node.js / Fastify Example

```javascript
// Step 1: JWT Authentication
app.decorate("authenticate", async function(request, reply) {
  try {
    await request.jwtVerify(); // Validates JWT signature
  } catch (err) {
    reply.send(err);
  }
});

// Step 2: Authorization Preparation
app.decorate("preauthorize", async function(request, reply) {
  const method = request.method;
  const documentName = request.params.documentName;

  // Map HTTP method to OpenFGA relation
  const relation = method === 'GET' ? 'reader' :
                   method === 'POST' ? 'writer' :
                   method === 'DELETE' ? 'owner' : null;

  // Extract user from JWT claims
  const user = `user:${request.user.username}`;

  // Format resource
  const object = `document:${documentName}`;

  request.authParams = { user, relation, object };
});

// Step 3: OpenFGA Check
app.addHook('preHandler', async (request, reply) => {
  const { user, relation, object } = request.authParams;

  const { allowed } = await fgaClient.check({
    user,
    relation,
    object
  });

  if (!allowed) {
    reply.code(403).send({ error: 'Forbidden' });
  }
});
```

#### Go / Fiber Example

```go
// JWT Middleware (Authentication)
app.Use(jwtware.New(jwtware.Config{
    SigningKey: []byte(os.Getenv("JWT_SECRET")),
}))

// Authorization Middleware
app.Use(func(c *fiber.Ctx) error {
    user := c.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    username := claims["username"].(string)

    // Map method to relation
    method := c.Method()
    relation := methodToRelation(method)

    // OpenFGA Check
    body, err := fgaClient.Check(context.Background()).
        Body(ClientCheckRequest{
            User:     fmt.Sprintf("user:%s", username),
            Relation: relation,
            Object:   fmt.Sprintf("document:%s", c.Params("name")),
        }).
        Execute()

    if err != nil || !body.GetAllowed() {
        return c.Status(403).JSON(fiber.Map{
            "error": "Forbidden",
        })
    }

    return c.Next()
})
```

### JWT Claims as Contextual Data

OpenFGA supports using JWT claims as contextual tuples for dynamic authorization. This pattern is useful for:

**Multi-Tenancy:** JWT contains current organization context
```javascript
// JWT payload
{
  "sub": "user:anne",
  "current_org": "org:acme"
}

// OpenFGA Check with contextual tuple
await fgaClient.check({
  user: "user:anne",
  relation: "can_manage_bots",
  object: "bot:battle-bot-1",
  contextualTuples: [
    {
      user: "user:anne",
      relation: "member",
      object: "org:acme"  // From JWT claim
    }
  ]
});
```

**Time-Based Access:** JWT contains role assignments with temporal context
**IP-Based Restrictions:** JWT contains location/IP information for access decisions

### OIDC Support

OpenFGA supports OIDC authentication for its own API access (not to be confused with application-level authentication):

- Client credentials flow with client ID and secret
- JWT tokens with `kid` header for key identification
- Integration with identity providers (Auth0, Keycloak, etc.)
- Audience and issuer validation

This enables secure API access for administrative operations and service-to-service communication.

### Best Practices for JWT + OpenFGA Integration

1. **Always verify JWT first:** Never call OpenFGA without valid authentication
2. **Extract user identity from JWT claims:** Use standardized claims (`sub`, `username`, etc.)
3. **Map operations to relations:** Create consistent mapping from HTTP methods/actions to OpenFGA relations
4. **Use contextual tuples sparingly:** Prefer storing relationships in OpenFGA when they're stable
5. **Cache authorization decisions:** Implement caching layer for frequently checked permissions
6. **Handle errors gracefully:** Distinguish between authentication failures and authorization denials

## Deployment Options

OpenFGA provides multiple deployment options suitable for different environments and requirements.

### Docker Deployment

#### Quick Start (Memory Storage)

The simplest deployment uses in-memory storage for development:

```bash
# Pull the latest image
docker pull openfga/openfga

# Run with memory storage (development only)
docker run -p 8080:8080 -p 8081:8081 -p 3000:3000 openfga/openfga run
```

**Exposed Ports:**
- **8080:** HTTP API
- **8081:** gRPC API
- **3000:** Playground UI (web-based testing interface)

**Warning:** Memory storage loses all data on restart. Not suitable for production.

#### Production with PostgreSQL

```bash
# Create Docker network
docker network create openfga

# Run PostgreSQL
docker run -d \
  --name postgres \
  --network openfga \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=openfga \
  postgres:17

# Run migrations (first time or after upgrades)
docker run --rm \
  --network openfga \
  openfga/openfga migrate \
  --datastore-engine postgres \
  --datastore-uri 'postgres://postgres:password@postgres:5432/openfga'

# Run OpenFGA server
docker run -d \
  --name openfga \
  --network openfga \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 3000:3000 \
  -e OPENFGA_DATASTORE_ENGINE=postgres \
  -e OPENFGA_DATASTORE_URI='postgres://postgres:password@postgres:5432/openfga' \
  openfga/openfga run
```

#### Production with MySQL

```bash
# Run MySQL
docker run -d \
  --name mysql \
  --network openfga \
  -e MYSQL_ROOT_PASSWORD=secret \
  -e MYSQL_DATABASE=openfga \
  mysql:8

# Run migrations
docker run --rm \
  --network openfga \
  openfga/openfga migrate \
  --datastore-engine mysql \
  --datastore-uri 'root:secret@tcp(mysql:3306)/openfga?parseTime=true'

# Run OpenFGA server
docker run -d \
  --name openfga \
  --network openfga \
  -p 8080:8080 \
  -e OPENFGA_DATASTORE_ENGINE=mysql \
  -e OPENFGA_DATASTORE_URI='root:secret@tcp(mysql:3306)/openfga?parseTime=true' \
  openfga/openfga run
```

#### Production with SQLite

```bash
# Create volume for persistence
docker volume create openfga-data

# Run OpenFGA with SQLite
docker run -d \
  --name openfga \
  -p 8080:8080 \
  -v openfga-data:/home/nonroot \
  -e OPENFGA_DATASTORE_ENGINE=sqlite \
  -e OPENFGA_DATASTORE_URI='/home/nonroot/openfga.db' \
  openfga/openfga run
```

#### Docker Compose Example

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:17
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: openfga
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - openfga

  openfga-migrate:
    image: openfga/openfga:latest
    command: migrate
    environment:
      OPENFGA_DATASTORE_ENGINE: postgres
      OPENFGA_DATASTORE_URI: 'postgres://postgres:password@postgres:5432/openfga'
    depends_on:
      - postgres
    networks:
      - openfga

  openfga:
    image: openfga/openfga:latest
    command: run
    environment:
      OPENFGA_DATASTORE_ENGINE: postgres
      OPENFGA_DATASTORE_URI: 'postgres://postgres:password@postgres:5432/openfga'
    ports:
      - "8080:8080"  # HTTP API
      - "8081:8081"  # gRPC API
      - "3000:3000"  # Playground
    depends_on:
      openfga-migrate:
        condition: service_completed_successfully
    networks:
      - openfga

volumes:
  postgres_data:

networks:
  openfga:
```

### Kubernetes Deployment

#### Helm Chart (Recommended)

Official Helm chart available at [ArtifactHub](https://artifacthub.io/packages/helm/openfga/openfga):

```bash
# Add Helm repository
helm repo add openfga https://openfga.github.io/helm-charts
helm repo update

# Install with default settings (memory storage)
helm install openfga openfga/openfga

# Install with PostgreSQL
helm install openfga openfga/openfga \
  --set datastore.engine=postgres \
  --set datastore.uri="postgres://user:pass@postgres:5432/openfga"
```

**Helm Chart Features:**
- **High Availability:** Defaults to 3 replicas for production (1 for memory storage)
- **Horizontal Scaling:** Support for multiple replicas with persistent datastores
- **Service Types:** ClusterIP, NodePort, LoadBalancer options
- **Health Probes:** Built-in liveness and readiness probes
- **Monitoring:** Integration with Prometheus and Grafana
- **Ingress:** Optional ingress configuration for external access

#### Kubernetes Operators

Multiple operator options available for managing OpenFGA deployments:

**1. ZEISS OpenFGA Operator**
- Repository: [github.com/ZEISS/openfga-operator](https://github.com/ZEISS/openfga-operator)
- Installation via Helm
- Manages OpenFGA lifecycle and configuration
- Custom Resource Definitions (CRDs) for OpenFGA instances

**2. Canonical OpenFGA Operator**
- Repository: [github.com/canonical/openfga-operator](https://github.com/canonical/openfga-operator)
- Juju Charm for Kubernetes
- Integration with Canonical Observability Stack (COS)
- Built-in Grafana dashboards
- Prometheus and Loki alert rules
- Available on [Charmhub](https://charmhub.io/openfga-k8s)

#### Sample Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openfga
spec:
  replicas: 3
  selector:
    matchLabels:
      app: openfga
  template:
    metadata:
      labels:
        app: openfga
    spec:
      containers:
      - name: openfga
        image: openfga/openfga:latest
        args: ["run"]
        env:
        - name: OPENFGA_DATASTORE_ENGINE
          value: "postgres"
        - name: OPENFGA_DATASTORE_URI
          valueFrom:
            secretKeyRef:
              name: openfga-db-secret
              key: database-uri
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 8081
          name: grpc
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: openfga
spec:
  selector:
    app: openfga
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: grpc
    port: 8081
    targetPort: 8081
  type: ClusterIP
```

### Binary Deployment

OpenFGA can run as a standalone binary:

```bash
# Download binary from GitHub releases
# https://github.com/openfga/openfga/releases

# Run with environment variables
export OPENFGA_DATASTORE_ENGINE=postgres
export OPENFGA_DATASTORE_URI='postgres://user:pass@localhost:5432/openfga'
./openfga run

# Or with command-line flags
./openfga run \
  --datastore-engine postgres \
  --datastore-uri 'postgres://user:pass@localhost:5432/openfga'

# Run migrations before first use or after upgrades
./openfga migrate \
  --datastore-engine postgres \
  --datastore-uri 'postgres://user:pass@localhost:5432/openfga'
```

### Storage Backend Options

| Storage Backend | Use Case | Production Ready | Data Persistence | Version Support |
|----------------|----------|------------------|------------------|-----------------|
| **Memory** | Development, testing | No | Lost on restart | All versions |
| **PostgreSQL** | Production | Yes | Full persistence | PostgreSQL 14+ |
| **MySQL** | Production | Yes | Full persistence | MySQL 8+ |
| **SQLite** | Small deployments, embedded | Beta | File-based persistence | SQLite 3+ |

**Production Recommendation:** Use PostgreSQL or MySQL with regular backups. Memory storage is only suitable for development and testing.

### Configuration Options

OpenFGA supports configuration via:

1. **Environment Variables:** Prefixed with `OPENFGA_` (e.g., `OPENFGA_DATASTORE_ENGINE`)
2. **Command-Line Flags:** Prefixed with `--` (e.g., `--datastore-engine`)
3. **Configuration File:** YAML-based configuration for complex setups

**Key Configuration Parameters:**
- `datastore.engine`: Storage backend (memory, postgres, mysql, sqlite)
- `datastore.uri`: Database connection string
- `http.addr`: HTTP API listen address (default: 0.0.0.0:8080)
- `grpc.addr`: gRPC API listen address (default: 0.0.0.0:8081)
- `playground.enabled`: Enable web playground (default: true)
- `playground.port`: Playground port (default: 3000)
- `authn.method`: Authentication method (none, preshared, oidc)
- `log.level`: Logging level (debug, info, warn, error)

### Security Configuration

#### Pre-Shared Keys

```bash
docker run -d \
  -e OPENFGA_AUTHN_METHOD=preshared \
  -e OPENFGA_AUTHN_PRESHARED_KEYS=key1,key2 \
  openfga/openfga run
```

#### OIDC Authentication

```bash
docker run -d \
  -e OPENFGA_AUTHN_METHOD=oidc \
  -e OPENFGA_AUTHN_OIDC_ISSUER=https://auth.example.com \
  -e OPENFGA_AUTHN_OIDC_AUDIENCE=openfga-api \
  openfga/openfga run
```

#### TLS Support

Mount certificate files and configure TLS:

```bash
docker run -d \
  -v /path/to/certs:/certs \
  -e OPENFGA_HTTP_TLS_ENABLED=true \
  -e OPENFGA_HTTP_TLS_CERT=/certs/server.crt \
  -e OPENFGA_HTTP_TLS_KEY=/certs/server.key \
  openfga/openfga run
```

### Monitoring and Profiling

```bash
# Enable profiling endpoint
docker run -d \
  -e OPENFGA_PROFILER_ENABLED=true \
  -e OPENFGA_PROFILER_ADDR=:3002 \
  -p 3002:3002 \
  openfga/openfga run
```

Access profiling data at `http://localhost:3002/debug/pprof/`

## Policy Definition

OpenFGA uses a configuration language to define authorization models. The language supports two formats: a human-friendly DSL and a JSON API format.

### Configuration Language Overview

**Purpose:** Define the types of objects in your system and how they relate to each other

**Two Formats:**
1. **DSL (Domain Specific Language):** User-friendly syntax for humans
2. **JSON:** API-compatible format for programmatic use

**Where DSL is Used:**
- Playground web interface
- OpenFGA CLI
- VS Code extension
- IntelliJ extension

**Where JSON is Used:**
- API calls
- SDK integration
- Programmatic model management

### Converting Between Formats

Use the OpenFGA CLI to convert:

```bash
# DSL to JSON
fga model transform --from dsl --to json

# JSON to DSL
fga model transform --from json --to dsl
```

### DSL Syntax

#### Basic Structure

```
model
  schema 1.1

type user

type document
  relations
    define owner: [user]
    define editor: [user]
    define viewer: [user]
    define can_view: viewer or editor or owner
    define can_edit: editor or owner
    define can_delete: owner
```

**Key Elements:**
- `model schema 1.1`: Schema version declaration
- `type`: Defines an object type
- `relations`: Section containing relationship definitions
- `define`: Declares a relation or computed permission
- `[user]`: Direct relationship type restrictions

#### Type Definitions

Types represent objects in your system:

```
type user                    # Users/subjects
type organization            # Organizations
type bot                     # Battle bots
type game                    # Games/matches
type document                # Documents
type folder                  # Folders
```

### Relationship Definitions

#### Direct Relationships

Specify which types can be directly assigned to a relation:

```
type document
  relations
    define owner: [user]
    define editor: [user, organization#member]
    define viewer: [user, user:*, organization#member]
```

**Formats:**
- `[user]`: Specific user objects
- `[user:*]`: Any user (public access)
- `[organization#member]`: Users with specific relation on another type

#### Union Operator (`or`)

User has permission if they match ANY condition:

```
type document
  relations
    define owner: [user]
    define editor: [user]
    define viewer: [user]
    define can_view: viewer or editor or owner
```

Example tuples:
```
user:anne, editor, document:1
```

Result: `anne` has `can_view` because she's an `editor`

#### Intersection Operator (`and`)

User has permission if they match ALL conditions:

```
type document
  relations
    define restricted: [user]
    define approved_viewer: [user]
    define can_view_restricted: approved_viewer and restricted
```

Example tuples:
```
user:anne, approved_viewer, document:1
user:anne, restricted, document:1
```

Result: `anne` has `can_view_restricted` only if both tuples exist

#### Exclusion Operator (`but not`)

Include users from base set while excluding others:

```
type document
  relations
    define viewer: [user]
    define blocked: [user]
    define can_view: viewer but not blocked
```

Example tuples:
```
user:anne, viewer, document:1
user:bob, viewer, document:1
user:bob, blocked, document:1
```

Result: `anne` has `can_view`, but `bob` does not (blocked)

**Use Cases:**
- Block lists
- Revoked access
- Temporary suspensions

#### Transitivity (`from`)

Inherit relationships through related objects:

```
type folder
  relations
    define owner: [user]
    define viewer: [user]

type document
  relations
    define parent: [folder]
    define owner: [user] or owner from parent
    define viewer: [user] or viewer from parent
```

Example tuples:
```
folder:docs, owner, user:anne
folder:docs, parent, document:readme
```

Result: `anne` is `owner` of `document:readme` (inherited from folder)

**Use Cases:**
- Hierarchical permissions (folder → documents)
- Organization membership (org → resources)
- Group-based access (group → members)

### Complex Authorization Models

#### Multi-Level Hierarchies

```
model
  schema 1.1

type user

type organization
  relations
    define member: [user]
    define admin: [user]

type team
  relations
    define parent_org: [organization]
    define member: [user] or admin from parent_org
    define admin: [user]

type bot
  relations
    define parent_team: [team]
    define owner: [user] or admin from parent_team
    define can_deploy: owner or admin from parent_team
    define can_configure: owner or admin from parent_team
    define can_view: member from parent_team or owner
```

This model enables:
- Organization admins automatically have admin rights on all teams
- Team admins can deploy and configure bots
- Team members can view bots
- Bot owners have full control

#### Battle Bots Example

```
model
  schema 1.1

type user

type organization
  relations
    define member: [user]
    define owner: [user]
    define admin: [user] or owner
    define can_invite_members: admin
    define can_create_teams: admin or member

type team
  relations
    define parent_org: [organization]
    define member: [user] or admin from parent_org
    define admin: [user]
    define can_manage_bots: admin or admin from parent_org

type bot
  relations
    define parent_team: [team]
    define creator: [user]
    define can_edit: creator or can_manage_bots from parent_team
    define can_delete: creator or admin from parent_team
    define can_deploy: creator or can_manage_bots from parent_team
    define can_view: member from parent_team

type game
  relations
    define participant_bot: [bot]
    define creator: [user]
    define can_view: creator or can_view from participant_bot
```

### JSON Format Example

The same model in JSON API format:

```json
{
  "schema_version": "1.1",
  "type_definitions": [
    {
      "type": "user"
    },
    {
      "type": "document",
      "relations": {
        "owner": {
          "this": {}
        },
        "editor": {
          "this": {}
        },
        "viewer": {
          "this": {}
        },
        "can_view": {
          "union": {
            "child": [
              {"this": {}},
              {"computedUserset": {"relation": "editor"}},
              {"computedUserset": {"relation": "owner"}}
            ]
          }
        }
      },
      "metadata": {
        "relations": {
          "owner": {
            "directly_related_user_types": [
              {"type": "user"}
            ]
          },
          "editor": {
            "directly_related_user_types": [
              {"type": "user"}
            ]
          },
          "viewer": {
            "directly_related_user_types": [
              {"type": "user"}
            ]
          }
        }
      }
    }
  ]
}
```

### Testing Authorization Models

#### Assertions

Define expected behavior to prevent regressions:

```yaml
# model_test.yaml
name: Document Permissions Test
model: |
  model
    schema 1.1

  type user

  type document
    relations
      define owner: [user]
      define viewer: [user]
      define can_view: viewer or owner

tuples:
  - user: user:anne
    relation: owner
    object: document:1
  - user: user:bob
    relation: viewer
    object: document:2

assertions:
  # Positive assertions
  - check:
      - user: user:anne
        relation: can_view
        object: document:1
        expected: true

  # Negative assertions
  - check:
      - user: user:bob
        relation: can_view
        object: document:1
        expected: false
```

Run tests with OpenFGA CLI:

```bash
fga model test --file model_test.yaml
```

### Modeling Best Practices

1. **Start Simple:** Begin with most critical feature, iterate to add complexity
2. **Use Meaningful Names:** Relation names should clearly indicate permission (e.g., `can_edit` not just `edit`)
3. **Leverage Inheritance:** Use `from` operator to avoid duplicating relationships
4. **Test Thoroughly:** Write assertions for expected and unexpected behaviors
5. **Version Models:** Track changes to authorization models like code
6. **Document Decisions:** Explain why relations are structured a certain way

### Tools for Modeling

- **Playground:** Web UI at `http://localhost:3000` for interactive modeling
- **CLI:** Command-line tool for model management and testing
- **VS Code Extension:** Syntax highlighting and validation
- **IntelliJ Extension:** IDE integration for JetBrains products

## Dynamic Data and Contextual Tuples

OpenFGA supports dynamic authorization decisions through contextual tuples, enabling runtime data to influence authorization without storing relationships permanently.

### What are Contextual Tuples?

**Definition:** Contextual tuples are temporary relationship tuples provided at query time that exist only for the duration of a specific authorization check.

**Key Characteristics:**
- Not persisted to the database
- Provided with each Check, BatchCheck, ListObjects, ListUsers, or Expand request
- Treated as if they exist in the store during evaluation
- Override database tuples if the same relationship exists
- Limited to 100 tuples per request

### How Contextual Tuples Work

#### Normal Authorization Check (Stored Tuples)

```javascript
// Relationship stored in database
await fgaClient.write({
  writes: [
    {
      user: "user:anne",
      relation: "member",
      object: "organization:acme"
    }
  ]
});

// Authorization check
const { allowed } = await fgaClient.check({
  user: "user:anne",
  relation: "can_manage_bots",
  object: "bot:battle-bot-1"
});
```

#### With Contextual Tuples (Dynamic Data)

```javascript
// No pre-stored relationship needed
// Context provided at query time
const { allowed } = await fgaClient.check({
  user: "user:anne",
  relation: "can_manage_bots",
  object: "bot:battle-bot-1",
  contextualTuples: [
    {
      user: "user:anne",
      relation: "member",
      object: "organization:acme"  // Dynamic context
    }
  ]
});
```

### Use Cases

#### 1. Multi-Tenancy / Organization Context

User works in multiple organizations; current context comes from session/JWT:

```javascript
// JWT payload contains current organization
const token = {
  sub: "user:anne",
  current_org: "org:acme",
  orgs: ["org:acme", "org:techcorp"]
};

// Authorization check uses JWT org context
const { allowed } = await fgaClient.check({
  user: token.sub,
  relation: "can_manage_bots",
  object: "bot:battle-bot-1",
  contextualTuples: [
    {
      user: token.sub,
      relation: "active_member",
      object: token.current_org  // From JWT
    }
  ]
});
```

**Authorization Model:**
```
type organization
  relations
    define member: [user]
    define active_member: [user]

type bot
  relations
    define parent_org: [organization]
    define can_manage: active_member from parent_org
```

**Benefits:**
- User can switch organizations without changing stored relationships
- Same user session, different authorization context
- No need to write/delete relationships on org switch

#### 2. Avoiding Data Synchronization

Some authorization data lives in external systems (identity provider, HR system):

```javascript
// User roles come from external IAM system
const externalRoles = await iamSystem.getUserRoles(userId);

// Use contextual tuples instead of syncing to OpenFGA
const contextualTuples = externalRoles.map(role => ({
  user: `user:${userId}`,
  relation: "member",
  object: `role:${role.id}`
}));

const { allowed } = await fgaClient.check({
  user: `user:${userId}`,
  relation: "can_deploy",
  object: "bot:battle-bot-1",
  contextualTuples
});
```

**Benefits:**
- No dual-write coordination
- Single source of truth in external system
- Always current data without sync delays

#### 3. Time-Based Access

Authorization depends on current time (though conditional relationships are preferred):

```javascript
const now = new Date();
const isBusinessHours = now.getHours() >= 9 && now.getHours() < 17;

const contextualTuples = [];
if (isBusinessHours) {
  contextualTuples.push({
    user: "user:anne",
    relation: "business_hours_user",
    object: "system:battlebots"
  });
}

const { allowed } = await fgaClient.check({
  user: "user:anne",
  relation: "can_deploy_to_production",
  object: "bot:battle-bot-1",
  contextualTuples
});
```

**Authorization Model:**
```
type system
  relations
    define business_hours_user: [user]

type bot
  relations
    define can_deploy_to_production: business_hours_user from parent_system
```

#### 4. IP-Based or Location Restrictions

```javascript
const userIp = request.ip;
const isInternalNetwork = ipRangeCheck(userIp, "10.0.0.0/8");

const contextualTuples = [];
if (isInternalNetwork) {
  contextualTuples.push({
    user: `user:${userId}`,
    relation: "internal_network_user",
    object: "network:corporate"
  });
}

const { allowed } = await fgaClient.check({
  user: `user:${userId}`,
  relation: "can_access_admin_panel",
  object: "system:battlebots",
  contextualTuples
});
```

#### 5. Disambiguating Multiple Relationships

User has multiple relationships with same object; specify which applies:

```javascript
// User is member of multiple teams
// Specify which team context for this request

const { allowed } = await fgaClient.check({
  user: "user:anne",
  relation: "can_view_metrics",
  object: "dashboard:team-performance",
  contextualTuples: [
    {
      user: "user:anne",
      relation: "viewing_as_member",
      object: "team:red-team"  // Which team context
    }
  ]
});
```

### Contextual Tuples API Examples

#### Node.js SDK

```javascript
const { OpenFgaClient } = require('@openfga/sdk');

const fgaClient = new OpenFgaClient({
  apiUrl: process.env.FGA_API_URL,
  storeId: process.env.FGA_STORE_ID,
  authorizationModelId: process.env.FGA_MODEL_ID,
});

// Check with contextual tuples
const response = await fgaClient.check({
  user: 'user:anne',
  relation: 'can_edit',
  object: 'document:1',
  contextualTuples: [
    {
      user: 'user:anne',
      relation: 'member',
      object: 'organization:acme'
    },
    {
      user: 'organization:acme',
      relation: 'owner',
      object: 'document:1'
    }
  ]
});

console.log(response.allowed); // true or false
```

#### Go SDK

```go
import (
    "context"
    "github.com/openfga/go-sdk/client"
)

func checkWithContext(userId, botId, orgId string) (bool, error) {
    resp, err := fgaClient.Check(context.Background()).
        Body(client.ClientCheckRequest{
            User:     fmt.Sprintf("user:%s", userId),
            Relation: "can_deploy",
            Object:   fmt.Sprintf("bot:%s", botId),
            ContextualTuples: []client.ClientTupleKey{
                {
                    User:     fmt.Sprintf("user:%s", userId),
                    Relation: "active_member",
                    Object:   fmt.Sprintf("organization:%s", orgId),
                },
            },
        }).
        Execute()

    if err != nil {
        return false, err
    }

    return resp.GetAllowed(), nil
}
```

#### REST API

```bash
curl -X POST http://localhost:8080/stores/${STORE_ID}/check \
  -H "Content-Type: application/json" \
  -d '{
    "authorization_model_id": "${MODEL_ID}",
    "tuple_key": {
      "user": "user:anne",
      "relation": "can_edit",
      "object": "document:1"
    },
    "contextual_tuples": {
      "tuple_keys": [
        {
          "user": "user:anne",
          "relation": "member",
          "object": "organization:acme"
        }
      ]
    }
  }'
```

### Limitations and Considerations

#### Hard Limits

| Limitation | Value | Reason |
|------------|-------|--------|
| **Maximum contextual tuples** | 100 per request | Performance and complexity bounds |
| **Tuple lifetime** | Single request | Not persisted to database |
| **Validation** | Same as stored tuples | Must conform to authorization model |

#### Precedence Rules

When both a contextual tuple and stored tuple exist with same user, relation, and object:

**Contextual tuple takes precedence** and stored tuple is ignored for that check.

```javascript
// Database has:
// user:anne, viewer, document:1

// Check with contextual tuple:
const { allowed } = await fgaClient.check({
  user: "user:anne",
  relation: "can_edit",
  object: "document:1",
  contextualTuples: [
    {
      user: "user:anne",
      relation: "editor",  // Different relation
      object: "document:1"
    }
  ]
});
// Evaluation uses contextual "editor" relation, not stored "viewer"
```

#### Performance Considerations

**Pros:**
- Avoids database writes for temporary/dynamic data
- Reduces data synchronization complexity
- Enables runtime-dependent authorization

**Cons:**
- Increases request payload size
- Cannot be indexed or optimized like stored tuples
- May increase check latency for complex models
- Requires client to provide context data

**Recommendation:** Use contextual tuples for truly dynamic data. For stable relationships, store tuples in database for better performance.

#### Security Considerations

**Trust Boundary:** Application must validate contextual tuple data before sending to OpenFGA.

```javascript
// ❌ DANGEROUS: Using user-provided data directly
const { allowed } = await fgaClient.check({
  user: currentUser,
  relation: "can_delete",
  object: targetBot,
  contextualTuples: req.body.context  // User controls this!
});

// ✅ SAFE: Validate and construct context from trusted sources
const orgId = await getValidatedOrgFromSession(req);
const { allowed } = await fgaClient.check({
  user: currentUser,
  relation: "can_delete",
  object: targetBot,
  contextualTuples: [
    {
      user: currentUser,
      relation: "member",
      object: `organization:${orgId}`  // Validated
    }
  ]
});
```

**Token Expiration Issue:** If using JWT claims as contextual tuples, access continues until token expires even if underlying claims change in identity provider.

**Mitigation:** Use short-lived tokens or implement token revocation checking.

### Best Practices

1. **Use for Truly Dynamic Data:** Session context, current time, IP address, temporary grants
2. **Store Stable Relationships:** User-organization membership, resource ownership, team membership
3. **Validate Context Data:** Never trust user-provided contextual tuples
4. **Limit Contextual Tuple Count:** Stay well below 100-tuple limit for performance
5. **Document Context Requirements:** Make clear what context is needed for each authorization check
6. **Cache When Appropriate:** Some dynamic data can be cached short-term to reduce computation

## Integration Patterns for Battle Bots

### Recommended Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Battle Bots Platform                   │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐      ┌──────────────┐                 │
│  │   Frontend   │      │  API Gateway  │                 │
│  │  (React/etc) │─────▶│   (Go)        │                 │
│  └──────────────┘      └───────┬───────┘                 │
│                                │                          │
│                                ▼                          │
│                    ┌───────────────────────┐             │
│                    │ Auth Middleware       │             │
│                    │ 1. Verify JWT         │             │
│                    │ 2. Extract user ID    │             │
│                    └───────────┬───────────┘             │
│                                │                          │
│                                ▼                          │
│                    ┌───────────────────────┐             │
│                    │ AuthZ Middleware      │             │
│                    │ 1. Map operation      │             │
│                    │ 2. Call OpenFGA       │             │
│                    │ 3. Allow/Deny         │             │
│                    └───────────┬───────────┘             │
│                                │                          │
│                    ┌───────────▼───────────┐             │
│                    │   Business Logic      │             │
│                    │   - Bot Management    │             │
│                    │   - Game Services     │             │
│                    │   - Team Management   │             │
│                    └───────────────────────┘             │
│                                                           │
└─────────────────────────────────────────────────────────┘
         │                              │
         │ GitHub OAuth                 │ gRPC/HTTP
         ▼                              ▼
┌─────────────────┐          ┌─────────────────┐
│  GitHub OAuth   │          │    OpenFGA      │
│                 │          │  - PostgreSQL   │
│                 │          │  - K8s/Docker   │
└─────────────────┘          └─────────────────┘
```

### Example: Bot Deployment Authorization

```go
// middleware/authz.go
package middleware

import (
    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    fgaClient "github.com/openfga/go-sdk/client"
)

type AuthZMiddleware struct {
    fga *fgaClient.OpenFgaClient
}

func (m *AuthZMiddleware) RequirePermission(
    resourceType string,
    relation string,
) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract user from JWT (set by auth middleware)
        userID := c.Locals("userID").(string)

        // Get resource ID from route params
        resourceID := c.Params("id")

        // Construct OpenFGA object
        object := fmt.Sprintf("%s:%s", resourceType, resourceID)
        user := fmt.Sprintf("user:%s", userID)

        // Check authorization
        resp, err := m.fga.Check(context.Background()).
            Body(fgaClient.ClientCheckRequest{
                User:     user,
                Relation: relation,
                Object:   object,
            }).
            Execute()

        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "Authorization check failed",
            })
        }

        if !resp.GetAllowed() {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden",
            })
        }

        return c.Next()
    }
}

// Usage in routes
app.Post("/bots/:id/deploy",
    authMiddleware.Authenticate,
    authzMiddleware.RequirePermission("bot", "can_deploy"),
    handlers.DeployBot,
)
```

### Relationship Management

```go
// services/organization.go
func (s *OrgService) AddTeamMember(
    ctx context.Context,
    teamID string,
    userID string,
) error {
    // Write relationship to OpenFGA
    _, err := s.fga.Write(ctx).
        Body(fgaClient.ClientWriteRequest{
            Writes: []fgaClient.ClientTupleKey{
                {
                    User:     fmt.Sprintf("user:%s", userID),
                    Relation: "member",
                    Object:   fmt.Sprintf("team:%s", teamID),
                },
            },
        }).
        Execute()

    return err
}

func (s *OrgService) RemoveTeamMember(
    ctx context.Context,
    teamID string,
    userID string,
) error {
    // Delete relationship from OpenFGA
    _, err := s.fga.Write(ctx).
        Body(fgaClient.ClientWriteRequest{
            Deletes: []fgaClient.ClientTupleKey{
                {
                    User:     fmt.Sprintf("user:%s", userID),
                    Relation: "member",
                    Object:   fmt.Sprintf("team:%s", teamID),
                },
            },
        }).
        Execute()

    return err
}
```

## Pros and Cons for Battle Bots

### Pros of Using OpenFGA

1. **Fine-Grained Authorization**
   - Object-level permissions instead of coarse role-based access
   - Can express complex relationships (org → team → bot → game)
   - Supports hierarchical permission inheritance
   - Scales to millions of authorization decisions per second

2. **Separation of Concerns**
   - Authentication (GitHub OAuth + JWT) completely separate from authorization
   - OpenFGA focuses solely on "who can do what"
   - Clean architecture with well-defined boundaries

3. **Flexible and Expressive**
   - Model any relationship pattern (ReBAC, RBAC, ABAC hybrid)
   - Union, intersection, exclusion, transitivity operators
   - Contextual tuples for dynamic authorization
   - Easy to evolve authorization model over time

4. **Production-Ready**
   - Inspired by Google Zanzibar (battle-tested at massive scale)
   - High performance with proper database backend
   - Horizontal scaling with stateless API servers
   - Multiple deployment options (Docker, K8s, binary)

5. **Developer Experience**
   - Intuitive DSL for modeling
   - SDKs for Go, JavaScript, Python, .NET, Java
   - Playground for testing models
   - CLI for model management and testing
   - IDE extensions (VS Code, IntelliJ)

6. **Open Source**
   - Apache 2.0 license
   - No vendor lock-in
   - Self-hosted (no per-user costs)
   - Active community and development
   - Backed by Okta/Auth0 but truly open

7. **Multi-Tenancy Native**
   - Perfect for organization/team/bot hierarchy
   - Contextual tuples for session-based org context
   - Clean permission inheritance across tenants

8. **Auditability**
   - All relationship changes can be logged
   - Query history for "why can user X access Y?"
   - Compliance-friendly authorization tracking

### Cons of Using OpenFGA

1. **Additional Infrastructure**
   - Requires separate service deployment
   - Database dependency (PostgreSQL/MySQL for production)
   - Monitoring and maintenance overhead
   - Network hop for authorization checks

2. **Learning Curve**
   - ReBAC is different from traditional RBAC
   - Authorization modeling requires thoughtful design
   - Team must understand relationship-based concepts
   - DSL syntax learning required

3. **Performance Considerations**
   - Every protected operation requires OpenFGA call
   - Network latency added to request path
   - Complex models can have slower evaluation times
   - Requires caching strategy for hot paths

4. **Operational Complexity**
   - Another service to deploy, monitor, and scale
   - Database backup and migration procedures
   - Version management for authorization models
   - Need observability into authorization decisions

5. **May Be Overkill for MVP**
   - Simple RBAC might suffice initially
   - Early-stage apps may not need fine-grained permissions
   - Added complexity before product-market fit
   - Could implement later when requirements are clearer

6. **Data Consistency Challenges**
   - Must keep OpenFGA relationships in sync with application data
   - Deleting resources requires cleaning up relationships
   - Risk of orphaned tuples if cleanup fails
   - Eventual consistency considerations

7. **Limited ABAC Support**
   - Primarily relationship-based, not attribute-based
   - Contextual tuples help but have limitations
   - Cannot directly query external attributes during evaluation
   - May need hybrid approach for complex ABAC requirements

## Recommendations for Battle Bots

### Recommendation 1: Adopt OpenFGA for Authorization

**Rationale:**
- Battle Bots has inherent hierarchy (org → team → bot → game)
- Fine-grained permissions are required (not all team members should deploy bots)
- Self-hosted OpenFGA aligns with control requirements
- Go platform matches well with OpenFGA's Go SDK
- No per-user costs (important for platform growth)

**When to Implement:**
- Design authorization model during MVP planning phase
- Implement alongside GitHub OAuth authentication
- Deploy before adding multi-user features

**Implementation Approach:**
1. Define authorization model for MVP (orgs, teams, bots)
2. Deploy OpenFGA with PostgreSQL on Kubernetes
3. Implement authorization middleware for all protected resources
4. Write relationships when resources are created/modified
5. Monitor and optimize based on usage patterns

### Recommendation 2: Start Simple, Iterate

**Phase 1: Basic Permissions (MVP)**
```
- Organization owners can create teams
- Team admins can manage bots
- Bot creators can edit/deploy their bots
- Team members can view team bots
```

**Phase 2: Advanced Permissions**
```
- Delegation (team admins grant deploy rights)
- Shared bots (multiple teams)
- Temporary access grants
- Audit logs
```

**Phase 3: Complex Scenarios**
```
- Tournament organizers
- Spectator access controls
- Sponsor/partner access
- API rate limiting by permission tier
```

### Recommendation 3: Design Authorization Model Carefully

**Avoid Over-Modeling:**
- Don't model every possible permission from day one
- Start with clear use cases
- Add relations as needed
- Test model thoroughly before production

**Example Starter Model for Battle Bots:**
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

type bot
  relations
    define parent_team: [team]
    define creator: [user]
    define can_view: member from parent_team
    define can_edit: creator or can_manage_bots from parent_team
    define can_deploy: creator or can_manage_bots from parent_team
    define can_delete: creator or admin from parent_team
```

### Recommendation 4: Implement Caching

**Authorization Check Caching:**
```go
// Cache authorization decisions for short periods
type AuthCache struct {
    cache *ttlcache.Cache
}

func (c *AuthCache) Check(
    user, relation, object string,
) (bool, error) {
    // Check cache first
    key := fmt.Sprintf("%s:%s:%s", user, relation, object)
    if cached, found := c.cache.Get(key); found {
        return cached.(bool), nil
    }

    // Call OpenFGA
    allowed, err := c.fgaCheck(user, relation, object)
    if err != nil {
        return false, err
    }

    // Cache for 30 seconds
    c.cache.Set(key, allowed, 30*time.Second)

    return allowed, nil
}
```

**When to Invalidate:**
- User permissions change
- Resource relationships modified
- Organization membership updated

### Recommendation 5: Monitor OpenFGA Performance

**Metrics to Track:**
- Authorization check latency (p50, p95, p99)
- Check throughput (requests/second)
- Model evaluation complexity
- Cache hit rates
- Database query performance

**Alerting Thresholds:**
- Check latency > 100ms (p95)
- Error rate > 1%
- Database connection pool exhaustion
- Tuple write failures

### Recommendation 6: Plan for Model Evolution

**Version Control:**
- Store authorization models in git
- Test model changes thoroughly
- Use separate stores for dev/staging/prod
- Document model changes in ADRs

**Migration Strategy:**
```bash
# Test new model in playground/staging
fga model write --file new-model.fga --store-id staging

# Run validation tests
fga model test --file model_test.yaml

# Deploy to production
fga model write --file new-model.fga --store-id production
```

## Alternative Considerations

### Alternative 1: Ory Keto

**Pros:**
- Similar to OpenFGA (also Zanzibar-inspired)
- Part of Ory ecosystem (with Kratos, Hydra)
- Open source (Apache 2.0)

**Cons:**
- Less mature than OpenFGA
- Smaller community
- Fewer SDKs and tools

### Alternative 2: SpiceDB

**Pros:**
- Also Zanzibar-inspired
- Strong consistency guarantees
- Built-in schema validation

**Cons:**
- More complex deployment
- Requires CockroachDB or PostgreSQL
- Different query language

### Alternative 3: Custom Authorization Service

**Pros:**
- Full control over implementation
- Optimized for specific use case
- No external dependencies

**Cons:**
- Significant development time
- Security burden on team
- Difficult to evolve
- Reinventing proven solutions

**Verdict:** OpenFGA provides the best balance of features, maturity, and operational simplicity for Battle Bots.

## Conclusion

OpenFGA is a production-ready, scalable authorization system that fits Battle Bots' requirements exceptionally well. Its separation from authentication allows seamless integration with the planned GitHub OAuth strategy, while its ReBAC model naturally represents the platform's organizational hierarchy (organizations → teams → bots → games).

**Key Takeaways:**

1. **JWT Integration:** OpenFGA works alongside JWT authentication, using extracted user identity for authorization checks
2. **Deployment Flexibility:** Multiple options (Docker, Kubernetes Helm charts, operators, binary) for various environments
3. **Intuitive Policy Definition:** DSL provides human-readable authorization models that compile to JSON for API use
4. **Dynamic Authorization:** Contextual tuples enable runtime data (session context, time, location) to influence decisions

**Recommended Path Forward:**
1. Design authorization model during planning phase
2. Deploy OpenFGA on Kubernetes with PostgreSQL
3. Start with simple permissions, iterate based on user needs
4. Implement caching for performance
5. Monitor and optimize as platform scales

OpenFGA's self-hosted nature eliminates per-user costs while providing enterprise-grade authorization capabilities, making it an excellent choice for Battle Bots' long-term growth.

## References

- [OpenFGA Official Documentation](https://openfga.dev/docs)
- [OpenFGA GitHub Repository](https://github.com/openfga/openfga)
- [OpenFGA Configuration Language](https://openfga.dev/docs/configuration-language)
- [OpenFGA Kubernetes Setup](https://openfga.dev/docs/getting-started/setup-openfga/kubernetes)
- [OpenFGA Docker Setup](https://openfga.dev/docs/getting-started/setup-openfga/docker)
- [OpenFGA Contextual Tuples](https://openfga.dev/docs/interacting/contextual-tuples)
- [OpenFGA Framework Integration](https://openfga.dev/docs/getting-started/framework)
- [OpenFGA Modeling Guide](https://openfga.dev/docs/modeling/getting-started)
- [OpenFGA Helm Chart](https://artifacthub.io/packages/helm/openfga/openfga)
- [ZEISS OpenFGA Operator](https://github.com/ZEISS/openfga-operator)
- [Canonical OpenFGA Operator](https://github.com/canonical/openfga-operator)
- [Google Zanzibar Paper](https://research.google/pubs/pub48190/)
- [Auth0 FGA Blog Posts](https://auth0.com/blog/tags/fga/)
- [OpenFGA CLI](https://github.com/openfga/cli)
- [OpenFGA Language Grammar](https://github.com/openfga/language)
