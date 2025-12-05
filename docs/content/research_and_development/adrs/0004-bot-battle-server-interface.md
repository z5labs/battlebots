---
title: "[0004] Bot to Battle Server Communication Protocol"
description: >
    Selection of gRPC with bidirectional streaming as the communication protocol for bot-to-server and bot-to-bot interfaces
type: docs
weight: 4
category: "strategic"
status: "proposed"
date: 2025-12-05
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

The Battle Bots platform requires a communication protocol for bots to interact with the battle server (client/server architecture) or with each other (peer-to-peer architecture). This protocol must support real-time bidirectional communication for game state updates and bot actions while remaining language-agnostic to enable bot development in any programming language.

Related to User Journey [0001] Proof of Concept - 1v1 Battle, which identifies bot-to-server interface as a pending ADR required for POC implementation.

**Key Requirements:**
- Bidirectional communication (bots send actions, receive game state/events)
- Language-agnostic interface (Python, Go, Java, JavaScript, Rust, etc.)
- Real-time or near-real-time performance
- Support for both client/server and P2P architectural modes
- Integration with OpenTelemetry observability stack (ADR-0002, ADR-0003)
- Container-friendly networking (Docker/Podman)

Which communication protocol should Battle Bots adopt for the bot-to-server and bot-to-bot interface?

## Decision Drivers

* **Performance**: Low latency and high throughput for real-time battle interaction
* **Language Support**: Must enable bot development in diverse programming languages
* **Observability**: Seamless integration with OpenTelemetry (ADR-0002) and OTLP stack (ADR-0003)
* **Developer Experience**: Ease of implementation for bot authors
* **Bidirectional Communication**: Efficient server-push and client-send capabilities
* **Type Safety**: Schema validation and versioning for protocol evolution
* **Container Networking**: Compatibility with Docker/Podman environments
* **Dual Architecture Support**: Viability for both client/server and P2P modes
* **Tooling Ecosystem**: Availability of debugging, testing, and development tools
* **Implementation Complexity**: Development and maintenance cost
* **Industry Maturity**: Proven technology with active community support

## Considered Options

* **Option 1**: gRPC with bidirectional streaming
* **Option 2**: WebSockets with JSON/MessagePack
* **Option 3**: HTTP/REST with Server-Sent Events (SSE)
* **Option 4**: Custom TCP protocol
* **Option 5**: Custom UDP protocol

## Decision Outcome

Chosen option: **"gRPC with bidirectional streaming"**, because it provides the optimal balance of performance, developer experience, and observability integration while meeting all functional requirements for both client/server and P2P architectures.

gRPC delivers near-WebSocket performance with superior type safety (Protocol Buffers), native OpenTelemetry instrumentation, and excellent language support through code generation. The bidirectional streaming model naturally fits the battle communication pattern where bots continuously send actions and receive game state updates.

### Consequences

**Positive:**
* ✅ **Excellent OpenTelemetry integration** - Native auto-instrumentation for traces, metrics, and logs with zero manual setup aligns perfectly with ADR-0002 and ADR-0003
* ✅ **Language-agnostic via Protocol Buffers** - Code generation for all major languages (Go, Python, Java, JavaScript, Rust, C++, C#) enables diverse bot ecosystem
* ✅ **Type safety and versioning** - .proto schema enforces contracts and provides forward/backward compatibility
* ✅ **Bidirectional streaming** - Natural fit for real-time battle where bots send actions and receive continuous state updates
* ✅ **Performance** - Binary protocol with 7-10x better throughput than REST/JSON and 30-50% smaller payloads
* ✅ **Container-friendly** - HTTP/2 works natively in Docker/Podman with simple port mapping
* ✅ **Rich tooling** - grpcurl for testing, ghz for benchmarking, reflection for discovery
* ✅ **Industry proven** - Used at scale by Google, Netflix, Square for production systems
* ✅ **Self-documenting** - .proto files serve as both implementation and documentation

**Negative:**
* ❌ **Learning curve** - Bot developers must understand Protocol Buffers syntax (mitigated by examples and generated code)
* ❌ **Build complexity** - Requires protoc compiler and code generation step in build pipeline
* ❌ **Binary debugging** - Not human-readable like JSON (mitigated by grpcurl and reflection)
* ❌ **Browser support** - Requires grpc-web proxy for browser-based bots (not native WebSocket)
* ❌ **P2P NAT traversal** - More complex than UDP for peer-to-peer due to HTTP/2 over TCP

**Neutral:**
* ⚪ **HTTP/2 requirement** - Benefits from multiplexing but limited to HTTP/2 capabilities
* ⚪ **Port considerations** - Non-standard ports (e.g., 50051) may require firewall rules in some environments

### Confirmation

Implementation compliance will be verified through:

1. **Protocol Definition**: All bot-to-server communication defined in `.proto` files in repository
2. **OTEL Instrumentation**: Automated tests verify traces/metrics are emitted to Tempo/Mimir
3. **Language SDK Examples**: Reference bot implementations in at least 3 languages (Go, Python, JavaScript)
4. **Integration Tests**: Client/server and P2P battle scenarios with containerized bots
5. **Performance Benchmarks**: Latency and throughput meet < 50ms round-trip target

## Pros and Cons of the Options

### gRPC with Bidirectional Streaming

**Description:** HTTP/2-based RPC framework with Protocol Buffers for serialization and bidirectional streaming for real-time communication.

**Detailed Analysis:** See [gRPC Protocol Analysis](/battlebots/research_and_development/analysis/protocols/grpc/grpc-analysis/)

**Pros:**
* ✅ **Native OpenTelemetry support** - Automatic trace propagation, span creation, and metrics collection integrate seamlessly with ADR-0002 observability stack
* ✅ **Protocol Buffers** - Strong typing, schema validation, and versioning prevent runtime errors and enable protocol evolution
* ✅ **Bidirectional streaming** - `stream BotAction → stream GameEvent` pattern matches battle communication model perfectly
* ✅ **Code generation** - Eliminates boilerplate, enforces API contracts across all languages
* ✅ **Performance** - Binary serialization achieves 5-20ms latency, thousands of messages/sec throughput
* ✅ **Language coverage** - Official support for Go, Python, Java, JavaScript, C++, C#, Rust, and more
* ✅ **Container networking** - HTTP/2 works natively in Docker/Podman without special configuration
* ✅ **Tooling ecosystem** - grpcurl, ghz, reflection, IDE plugins provide excellent developer experience
* ✅ **Dual architecture support** - Client/server: bots as clients. P2P: bots expose gRPC servers and connect to each other

**Cons:**
* ❌ **Learning curve** - Developers must learn .proto syntax (estimated 1-2 hours for basics)
* ❌ **Build step** - Code generation adds complexity to build pipelines
* ❌ **Binary format** - Debugging requires tools like grpcurl vs simple curl for JSON
* ❌ **Browser limitation** - Requires grpc-web proxy, not native browser support
* ❌ **P2P NAT** - TCP-based NAT traversal requires rendezvous server or relay for some networks

**Suitability:**
* **Client/Server**: ⭐⭐⭐⭐⭐ (Excellent) - Ideal fit for all requirements
* **P2P**: ⭐⭐⭐⭐☆ (Very Good) - NAT traversal manageable with coordination server

### WebSockets with JSON/MessagePack

**Description:** Full-duplex bidirectional communication over persistent TCP connection with JSON or binary serialization.

**Detailed Analysis:** See [HTTP-based Protocols Analysis](/battlebots/research_and_development/analysis/protocols/http/http-analysis/)

**Pros:**
* ✅ **True bidirectional** - Full-duplex communication without polling
* ✅ **Low latency** - 5-20ms, comparable to gRPC
* ✅ **Universal support** - Native in browsers and all major languages
* ✅ **Flexible serialization** - JSON for debugging, MessagePack for efficiency
* ✅ **Familiar to web developers** - Lower barrier than learning gRPC/protobuf
* ✅ **Simple protocol** - No code generation or build complexity

**Cons:**
* ❌ **Manual OpenTelemetry instrumentation** - No automatic trace propagation or span creation
* ❌ **No schema enforcement** - JSON lacks type safety, versioning requires custom logic
* ❌ **Sticky sessions required** - Load balancing complexity for stateful connections
* ❌ **No code generation** - Manual client/server implementation in each language
* ❌ **HTTP/1.1 only** - No HTTP/2 multiplexing benefits
* ❌ **Reconnection logic** - Client must implement reconnect and state recovery

**Suitability:**
* **Client/Server**: ⭐⭐⭐⭐☆ (Very Good) - Viable, but lacks gRPC's OTEL integration
* **P2P**: ⭐⭐⭐☆☆ (Moderate) - Similar NAT challenges as gRPC

### HTTP/REST with Server-Sent Events

**Description:** Hybrid approach using HTTP POST for bot actions and SSE for server-to-bot event streaming.

**Detailed Analysis:** See [HTTP-based Protocols Analysis](/battlebots/research_and_development/analysis/protocols/http/http-analysis/)

**Pros:**
* ✅ **Automatic reconnection** - SSE handles reconnect with Last-Event-ID resume
* ✅ **Simple model** - Standard HTTP requests plus event stream
* ✅ **Excellent OTEL support** - HTTP instrumentation is mature
* ✅ **Text-based** - Easy debugging with curl and browser DevTools
* ✅ **Browser native** - EventSource API built into browsers

**Cons:**
* ❌ **Unidirectional streaming** - Requires separate HTTP POST for bot → server
* ❌ **Asymmetric performance** - Different latency for upstream vs downstream
* ❌ **HTTP/1.1 connection limits** - Browser constraint (6 per domain)
* ❌ **Text-only** - No native binary support (requires base64 encoding)
* ❌ **Less common** - Fewer production examples than WebSocket or gRPC
* ❌ **Two separate channels** - More complex than single bidirectional stream

**Suitability:**
* **Client/Server**: ⭐⭐⭐☆☆ (Moderate) - Workable but awkward bidirectional pattern
* **P2P**: ⭐⭐☆☆☆ (Marginal) - SSE not designed for P2P scenarios

### Custom TCP Protocol

**Description:** Application-specific binary protocol over raw TCP sockets with custom message framing.

**Detailed Analysis:** See [Custom TCP/UDP Protocol Analysis](/battlebots/research_and_development/analysis/protocols/custom/tcp-udp-analysis/)

**Pros:**
* ✅ **Maximum control** - Full control over wire format and protocol behavior
* ✅ **Minimal overhead** - No HTTP headers or protocol baggage

**Cons:**
* ❌ **Very high implementation cost** - 4-6 weeks per language for production-quality code
* ❌ **No performance advantage** - TCP latency same as WebSocket/gRPC (both use TCP)
* ❌ **Manual OTEL integration** - No automatic instrumentation, requires custom trace context
* ❌ **No tooling** - Must build custom debugging and testing tools
* ❌ **Maintenance burden** - Protocol bugs, versioning, cross-platform issues
* ❌ **Language fragmentation** - Manual implementation in each language, consistency challenges

**Suitability:**
* **Client/Server**: ⭐☆☆☆☆ (Not Recommended) - Cost far exceeds marginal benefit
* **P2P**: ⭐☆☆☆☆ (Not Recommended) - Same TCP limitations as gRPC/WebSocket

**Verdict:** ❌ **Rejected** - No measurable performance benefit over gRPC/WebSocket while requiring 10-20x more development effort.

### Custom UDP Protocol

**Description:** Connectionless packet-based protocol with custom reliability layer for critical messages.

**Detailed Analysis:** See [Custom TCP/UDP Protocol Analysis](/battlebots/research_and_development/analysis/protocols/custom/tcp-udp-analysis/)

**Pros:**
* ✅ **Lowest latency** - 2-10ms potential vs 5-20ms for TCP
* ✅ **No head-of-line blocking** - Lost packets don't block subsequent packets
* ✅ **Suitable for high tick-rate** - Games running at 60+ updates/sec benefit

**Cons:**
* ❌ **Very high complexity** - 8-12 weeks per language plus reliability layer
* ❌ **Unclear requirement** - Battle Bots tick rate and latency needs undefined
* ❌ **Firewall hostile** - Many networks block UDP, requires fallback
* ❌ **NAT traversal complexity** - Hole punching, STUN/TURN infrastructure needed for P2P
* ❌ **Extreme OTEL difficulty** - No standard trace propagation, manual packet-level instrumentation
* ❌ **Premature optimization** - POC should prove < 50ms latency is insufficient before custom UDP

**Suitability:**
* **Client/Server**: ⭐☆☆☆☆ (Not Recommended for POC) - Defer until profiling proves TCP inadequate
* **P2P**: ⭐⭐☆☆☆ (Marginal) - NAT traversal challenges significant

**Verdict:** ❌ **Rejected for POC** - Re-evaluate only if client/server POC demonstrates that TCP-based protocols cannot meet latency requirements (evidence currently absent).

## Interface Specification

### Protocol Buffer Schema

**Core Message Types:**

```protobuf
syntax = "proto3";

package battlebots.v1;

// Battle service defining bot-to-server interaction
service Battle {
  // Bidirectional streaming for real-time battle
  rpc Stream(stream BotMessage) returns (stream ServerMessage);
  
  // Unary RPC for initial connection handshake
  rpc Connect(ConnectRequest) returns (ConnectResponse);
}

// Bot → Server messages
message BotMessage {
  string bot_id = 1;
  int64 timestamp = 2;
  
  oneof payload {
    MoveAction move = 10;
    AttackAction attack = 11;
    DefendAction defend = 12;
    Heartbeat heartbeat = 13;
  }
}

message MoveAction {
  Position target = 1;
}

message AttackAction {
  string target_bot_id = 1;
  int32 weapon_id = 2;
}

message DefendAction {
  bool active = 1;
}

message Heartbeat {
  int64 client_tick = 1;
}

// Server → Bot messages
message ServerMessage {
  int64 server_tick = 1;
  
  oneof payload {
    GameState state = 10;
    ActionResult result = 11;
    BattleEvent event = 12;
  }
}

message GameState {
  repeated BotState bots = 1;
  repeated Projectile projectiles = 2;
  int32 tick = 3;
}

message BotState {
  string bot_id = 1;
  Position position = 2;
  int32 health = 3;
  int32 energy = 4;
}

message ActionResult {
  bool success = 1;
  string error = 2;
}

message BattleEvent {
  enum Type {
    BOT_DAMAGED = 0;
    BOT_DESTROYED = 1;
    BATTLE_ENDED = 2;
  }
  Type type = 1;
  map<string, string> metadata = 2;
}

// Common types
message Position {
  float x = 1;
  float y = 2;
}

message Projectile {
  string id = 1;
  Position position = 2;
  Position velocity = 3;
}

message ConnectRequest {
  string bot_id = 1;
  string api_token = 2;
}

message ConnectResponse {
  bool success = 1;
  string session_id = 2;
}
```

### Connection Lifecycle

**Client/Server Mode:**

1. **Connection**: Bot initiates gRPC connection to battle server
   ```go
   conn, err := grpc.Dial("battle-server:50051", grpc.WithInsecure())
   client := battlebots.NewBattleClient(conn)
   ```

2. **Authentication**: Unary `Connect()` RPC with API token
   ```go
   resp, err := client.Connect(ctx, &ConnectRequest{
       BotId: "bot-123",
       ApiToken: "secret-token",
   })
   ```

3. **Streaming**: Bidirectional `Stream()` RPC for battle duration
   ```go
   stream, err := client.Stream(ctx)
   
   // Send actions
   go func() {
       stream.Send(&BotMessage{Action: moveAction})
   }()
   
   // Receive state
   for {
       msg, err := stream.Recv()
       // Process game state
   }
   ```

4. **Disconnection**: Close stream gracefully or handle abrupt disconnect

**P2P Mode:**

1. **Rendezvous**: Bots connect to coordination server to exchange endpoints
2. **Server Setup**: Each bot runs gRPC server on known port
   ```go
   listener, _ := net.Listen("tcp", ":50051")
   grpcServer := grpc.NewServer()
   battlebots.RegisterBattleServer(grpcServer, &botServer{})
   grpcServer.Serve(listener)
   ```

3. **Peer Connection**: Bots connect to each other as gRPC clients
4. **Bidirectional Streaming**: Use `Stream()` RPC for bot-to-bot communication
5. **Consensus**: Application-level consensus for game state (future ADR)

### Error Handling

**gRPC Status Codes:**
- `OK` - Action succeeded
- `INVALID_ARGUMENT` - Malformed bot action
- `UNAUTHENTICATED` - Invalid API token
- `RESOURCE_EXHAUSTED` - Rate limit exceeded
- `UNAVAILABLE` - Server temporarily unavailable
- `DEADLINE_EXCEEDED` - Action timeout

**Client Retry Strategy:**
- Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
- Max retries: 5
- Reconnect on stream errors

### Versioning Strategy

**Protocol Buffers Backward Compatibility:**
- Never change field numbers
- Mark deprecated fields with `[deprecated = true]`
- Add new fields with unique numbers
- Use `reserved` for removed fields

**Service Versioning:**
- Package naming: `battlebots.v1`, `battlebots.v2`
- Server supports multiple versions simultaneously
- Clients specify version in package import

**Example Evolution:**
```protobuf
// v1
message BotMessage {
  string bot_id = 1;
  MoveAction move = 2;
}

// v2 - add new action type
message BotMessage {
  string bot_id = 1;
  MoveAction move = 2;
  TeleportAction teleport = 3; // New field
}
```

Old clients continue working (ignore unknown fields), new clients use new features.

## More Information

### Related ADRs

- **ADR-0002**: OpenTelemetry SDK Selection - Establishes OTEL as observability standard
- **ADR-0003**: Observability Stack Architecture - Defines Tempo, Mimir, Loki as backends
- **Future ADR**: Game Runtime Architecture - Will define tick rate and game loop mechanics
- **Future ADR**: Client/Server vs P2P Architecture - Will choose primary architecture mode

### Related Requirements

From User Journey [0001] Proof of Concept - 1v1 Battle:
- Language-agnostic bot implementation
- Containerized bots (Docker/Podman)
- Observability signals capture
- Real-time battle visualization

### Implementation Guidance

**Getting Started for Bot Developers:**

1. Install protoc compiler and gRPC plugin for your language
2. Clone battle bots repository with `.proto` definitions
3. Generate client code: `protoc --go_out=. --go-grpc_out=. battle.proto`
4. Implement bot logic using generated stubs
5. Connect to battle server and participate in battles

**Example Bot (Python):**
```python
import grpc
from battlebots.v1 import battle_pb2, battle_pb2_grpc

channel = grpc.insecure_channel('battle-server:50051')
stub = battle_pb2_grpc.BattleStub(channel)

# Connect
response = stub.Connect(battle_pb2.ConnectRequest(
    bot_id='my-bot',
    api_token='secret'
))

# Stream battle
stream = stub.Stream(iter([
    battle_pb2.BotMessage(
        bot_id='my-bot',
        move=battle_pb2.MoveAction(target=battle_pb2.Position(x=10, y=20))
    )
]))

for server_msg in stream:
    print(f"Game state: {server_msg.state}")
```

### Open Questions

1. **Tick Rate**: What is the target game tick rate? (Informs whether gRPC performance is sufficient)
2. **P2P Consensus**: How will bots reach consensus in P2P mode? (Future ADR)
3. **Browser Bots**: Should we support browser-based bots via grpc-web?
4. **Rate Limiting**: What are the action rate limits per bot?

### Further Reading

- [gRPC Protocol Analysis](../analysis/protocols/grpc/grpc-analysis.md) - Detailed gRPC evaluation
- [HTTP-based Protocols Analysis](../analysis/protocols/http/http-analysis.md) - WebSocket and SSE analysis
- [Custom Protocols Analysis](../analysis/protocols/custom/tcp-udp-analysis.md) - TCP/UDP evaluation
- [gRPC Official Documentation](https://grpc.io/docs/)
- [Protocol Buffers Language Guide](https://protobuf.dev/programming-guides/proto3/)
- [OpenTelemetry gRPC Instrumentation](https://opentelemetry.io/docs/instrumentation/)
