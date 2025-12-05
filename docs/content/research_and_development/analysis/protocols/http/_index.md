---
title: "HTTP-based Protocols Analysis"
description: >
    Analysis of HTTP/REST, WebSockets, and SSE for bot-to-server communication
type: docs
weight: 2
date: 2025-12-05
---

## Overview

HTTP-based protocols represent the most widely-adopted communication patterns on the web. This analysis evaluates three approaches for Battle Bots bot-to-server communication:

1. **HTTP/REST** - Request/response pattern with polling or long-polling
2. **WebSockets** - Full-duplex bidirectional communication over TCP
3. **Server-Sent Events (SSE)** - Unidirectional server-to-client streaming

Each approach offers different trade-offs in complexity, performance, and developer familiarity.

## HTTP/REST Analysis

### Architecture Pattern

**Traditional Request/Response:**
```
Bot → POST /battles/{id}/actions → Server
Bot ← 200 OK {result} ← Server

Bot → GET /battles/{id}/state → Server
Bot ← 200 OK {gameState} ← Server
```

**RESTful API Design:**
- `POST /battles/{id}/actions` - Submit bot action
- `GET /battles/{id}/state` - Poll for current game state
- `GET /battles/{id}/events?since={timestamp}` - Fetch events since last poll
- `DELETE /battles/{id}/connections/{botId}` - Disconnect from battle

### State Updates

**Polling:**
- Bot repeatedly requests `/state` endpoint (e.g., every 100ms)
- Simple to implement
- High latency (up to poll interval)
- Wasted requests when no state changes

**Long-Polling:**
- Server holds request open until state changes
- Immediately returns new state when available
- Better latency than polling, worse than streaming
- Reconnection overhead between state changes

### Performance Characteristics

**Latency:**
- **Polling**: 50-500ms depending on interval
- **Long-polling**: 10-100ms with reconnection overhead
- **HTTP/1.1**: Head-of-line blocking, connection limits
- **HTTP/2**: Multiplexing improves concurrent requests

**Throughput:**
- Limited by request/response cycle time
- High overhead for frequent small updates
- Better for discrete actions than continuous state

**Resource Usage:**
- Connection churning with polling/long-polling
- Server load from repeated requests
- Client CPU from polling loops

### Advantages
✅ Universal language support (every language has HTTP client)  
✅ Simple mental model - stateless request/response  
✅ Easy debugging with curl, browser DevTools  
✅ RESTful conventions well-understood  
✅ HTTP caching for read-heavy scenarios  
✅ Firewall-friendly (port 80/443)  

### Disadvantages
❌ Poor real-time performance without long-polling  
❌ Polling wastes bandwidth and CPU  
❌ High latency for state updates  
❌ No server push (without long-polling hacks)  
❌ Inefficient for streaming game state  

### Suitability for Battle Bots
⭐⭐☆☆☆ (Marginal) - Only viable for very slow turn-based games with infrequent updates.

## WebSockets Analysis

### Architecture Pattern

**Full-Duplex Communication:**
```
Bot → WS Connect ws://server/battle/123 → Server
     ← WS Upgrade (101 Switching Protocols) ←

Bot → {"action": "move", "x": 10} →
     ← {"event": "state_update", "positions": [...]} ←
     
Bot → {"action": "attack", "target": "bot2"} →
     ← {"event": "damage", "target": "bot2", "hp": 50} ←
```

**Persistent Connection:**
- Single TCP connection for bidirectional messages
- No HTTP overhead after handshake
- Both sides can send messages anytime
- Connection stays open for battle duration

### Message Framing

**Text Frames:**
- JSON messages (human-readable)
- Easy debugging and development
- Larger payload sizes

**Binary Frames:**
- MessagePack, Protocol Buffers, CBOR
- Efficient encoding/decoding
- Smaller payload sizes (30-50% reduction)

### Performance Characteristics

**Latency:**
- **5-20ms** for message round-trip (local network)
- Near-zero after connection establishment
- No polling overhead
- Comparable to gRPC streaming

**Throughput:**
- Thousands of messages per second per connection
- Limited by TCP bandwidth, not protocol overhead
- Efficient for rapid state updates

**Resource Usage:**
- One TCP connection per bot (persistent)
- Lower CPU than polling (no repeated handshakes)
- Memory for message buffering
- Connection scaling considerations (1000s of concurrent bots)

### Connection Management

**Lifecycle:**
1. HTTP upgrade handshake
2. Persistent WebSocket connection
3. Keepalive pings/pongs
4. Graceful or abrupt close

**Reconnection:**
- Not automatic - client must implement
- State synchronization after reconnect
- Resume semantics (message IDs, sequence numbers)

**Heartbeats:**
- Ping/Pong frames for keepalive
- Detect half-open connections
- Configurable timeout

### Language & Platform Support

**Excellent support across languages:**
- **Go**: `gorilla/websocket`, `nhooyr/websocket`
- **Python**: `websockets`, `aiohttp`
- **JavaScript**: Native WebSocket API (browser + Node.js)
- **Java**: Java WebSocket API (JSR 356), Spring WebSocket
- **Rust**: `tokio-tungstenite`
- **C#**: ASP.NET Core SignalR

**Browser Native:**
- WebSockets work in all modern browsers
- Enables browser-based bot development
- Web-based battle visualization

### OpenTelemetry Integration

**Challenges:**
- **No standard instrumentation** like gRPC
- **Manual span creation** for each message
- **Context propagation** requires custom headers/metadata

**Implementation Approach:**
```go
// Manual tracing
ctx, span := tracer.Start(ctx, "websocket.message")
defer span.End()

span.SetAttributes(
    attribute.String("message.type", msg.Type),
    attribute.String("bot.id", botID),
)
```

**Metrics:**
- Message count, size, latency
- Connection duration, errors
- Custom battle-specific metrics

**Trace Correlation:**
- Embed trace context in JSON messages
- Reconstruct distributed trace manually
- More complex than gRPC auto-instrumentation

### Container Networking

**Proxy Considerations:**
- **Sticky sessions required**: Connection affinity
- **Nginx**: `ip_hash` or `sticky` directive
- **HAProxy**: `stick-table` for session persistence
- **Envoy**: Consistent hashing

**HTTP/1.1 Limitation:**
- WebSocket upgrade uses HTTP/1.1
- One WebSocket per TCP connection
- No multiplexing like HTTP/2

**Kubernetes:**
- Service type LoadBalancer with session affinity
- Ingress controllers support WebSocket
- Readiness probes need special handling

### Development Experience

**Tooling:**
- **wscat**: Command-line WebSocket client
- **Browser DevTools**: Inspect WebSocket frames
- **Postman**: WebSocket request support
- **websocat**: Netcat-like WebSocket tool

**Debugging:**
- Messages visible in browser DevTools
- Text frames easy to inspect
- Binary frames require decoder

**Testing:**
- Unit test with mock WebSocket
- Integration test with test server
- Load testing with `ws-bench`

### Advantages

✅ **True bidirectional communication** (no polling)  
✅ **Low latency** for real-time updates  
✅ **Efficient** - single persistent connection  
✅ **Universal support** across languages and browsers  
✅ **Familiar** to web developers  
✅ **Flexible** - text or binary messages  
✅ **Firewall-friendly** (works on port 80/443)  

### Disadvantages

❌ **No automatic reconnection** (client must implement)  
❌ **Scaling challenges** (sticky sessions, connection limits)  
❌ **OpenTelemetry integration** requires manual instrumentation  
❌ **No built-in backpressure** mechanism  
❌ **HTTP/1.1 only** (no HTTP/2 multiplexing)  
❌ **Connection state** complicates load balancing  

### Suitability for Battle Bots

**Client/Server Architecture**: ⭐⭐⭐⭐☆ (Very Good)
- Excellent real-time performance
- Manual OTEL integration is manageable
- Sticky sessions solvable with proper load balancing

**P2P Architecture**: ⭐⭐⭐☆☆ (Moderate)
- Bots must run WebSocket servers
- NAT traversal challenges (similar to gRPC)
- Discovery and connection complexity

## Server-Sent Events (SSE) Analysis

### Architecture Pattern

**Unidirectional Server → Client:**
```
Bot → GET /battles/{id}/events → Server
     ← HTTP/1.1 200 OK
     ← Content-Type: text/event-stream
     
     ← data: {"event": "state_update", "positions": [...]}
     
     ← data: {"event": "damage", "target": "bot2", "hp": 50}
     
Bot → POST /battles/{id}/actions → Server (separate connection)
```

**Hybrid Approach:**
- SSE for server → bot (game state, events)
- HTTP POST for bot → server (actions)

### Protocol Details

**Event Stream Format:**
```
event: gameStateUpdate
data: {"positions": [...], "tick": 42}
id: 42

event: botDamaged
data: {"botId": "bot2", "damage": 10}
id: 43
```

- Text-based format
- Named events for filtering
- Auto-incrementing IDs for resume

**Automatic Reconnection:**
- Browser automatically reconnects on disconnect
- `Last-Event-ID` header for resuming from last event
- Built-in retry with exponential backoff

### Performance Characteristics

**Latency:**
- **10-30ms** for server → bot messages
- Similar to WebSocket for downstream
- POST latency for bot → server actions (20-50ms)

**Throughput:**
- Excellent for server → bot (streaming)
- Limited by HTTP POST for bot → server
- Asymmetric performance

**Resource Usage:**
- One long-lived connection per bot (SSE)
- Short-lived connections for actions (POST)
- Lower than polling, higher than WebSocket

### Advantages

✅ **Automatic reconnection** with event ID resume  
✅ **Simple** - just HTTP GET, no upgrade handshake  
✅ **Browser native** - EventSource API  
✅ **Text-based** - easy debugging  
✅ **Firewall-friendly** (HTTP)  
✅ **Built-in event types** and filtering  

### Disadvantages

❌ **Unidirectional only** (server → client)  
❌ **Requires separate channel** for client → server (HTTP POST)  
❌ **Text-only** (no binary without base64 encoding)  
❌ **HTTP/1.1 connection limits** (6 per domain in browsers)  
❌ **Less efficient** than WebSocket for bidirectional  
❌ **Not widely adopted** outside browser contexts  

### Suitability for Battle Bots
⭐⭐⭐☆☆ (Moderate) - Viable for slow-paced battles, but WebSocket is simpler for bidirectional needs.

## Serialization Format Comparison

### JSON (JavaScript Object Notation)

**Characteristics:**
- Text-based, human-readable
- Language-agnostic
- Schema-optional (self-describing)

**Pros:**
- Easy debugging (read with eyes)
- Universal support
- Browser-native parsing
- No build step

**Cons:**
- Larger payload size (2-3x vs binary)
- Slower parsing than binary
- No built-in versioning
- No schema validation

**Example:**
```json
{
  "action": "move",
  "botId": "bot-123",
  "position": {"x": 10, "y": 20},
  "timestamp": 1701820800
}
```
**Size**: ~120 bytes

### MessagePack

**Characteristics:**
- Binary JSON-like format
- Schema-optional
- Faster and smaller than JSON

**Pros:**
- 30-50% smaller than JSON
- 2-5x faster encoding/decoding
- Language support: Go, Python, JS, Java, Rust
- Drop-in JSON replacement

**Cons:**
- Not human-readable
- Less ubiquitous than JSON
- No schema enforcement

**Example (binary):**
```
\x84\xa6action\xa4move\xa5botId\xa7bot-123\xa8position\x82...
```
**Size**: ~70 bytes

### Protocol Buffers over HTTP

**Characteristics:**
- Binary protocol with schema (.proto)
- Same as gRPC but over HTTP/WebSocket
- Strong typing and versioning

**Pros:**
- Smallest payload (50-70% smaller than JSON)
- Type safety and validation
- Forward/backward compatibility
- Code generation

**Cons:**
- Requires .proto files and code generation
- Not human-readable
- Build complexity

**Size**: ~50 bytes

### CBOR (Concise Binary Object Representation)

**Characteristics:**
- Binary JSON alternative (RFC 8949)
- Self-describing like JSON
- Designed for IoT/constrained environments

**Pros:**
- Smaller than JSON
- More data types (binary, dates)
- IETF standard

**Cons:**
- Less adoption than MessagePack
- Not as compact as Protocol Buffers

### Recommendation for Battle Bots

**Development/Debugging**: JSON (easy to inspect)  
**Production**: Protocol Buffers (best performance + type safety)  
**Alternative**: MessagePack (good middle ground)

## OpenTelemetry Integration Summary

| Protocol | OTEL Support | Implementation Effort | Trace Propagation | Auto-Instrumentation |
|----------|--------------|----------------------|-------------------|---------------------|
| HTTP/REST | ⭐⭐⭐⭐⭐ | Easy | W3C Trace Context headers | Yes (most languages) |
| WebSockets | ⭐⭐⭐☆☆ | Moderate | Manual (custom metadata) | No |
| SSE | ⭐⭐⭐☆☆ | Moderate | W3C headers on GET | Partial |

**Best OTEL Integration**: HTTP/REST (but worst real-time performance)  
**Compromise**: WebSocket with manual instrumentation

## Container Networking Summary

| Protocol | Container Support | Load Balancing | Session Affinity | NAT Traversal |
|----------|------------------|----------------|------------------|---------------|
| HTTP/REST | Excellent | Easy | Not required | N/A |
| WebSockets | Good | Moderate | Required | Challenging |
| SSE | Good | Moderate | Required | Moderate |

## Development Experience Summary

| Protocol | Learning Curve | Tooling | Debugging | Language Support |
|----------|---------------|---------|-----------|------------------|
| HTTP/REST | Low | Excellent | Easy | Universal |
| WebSockets | Low-Moderate | Good | Moderate | Excellent |
| SSE | Low | Moderate | Easy | Good |

## Pros and Cons Summary

### Overall Assessment

**HTTP/REST:**
- ✅ Simple, universal, well-tooled
- ❌ Poor real-time performance
- **Use case**: Admin APIs, non-real-time operations

**WebSockets:**
- ✅ Excellent real-time performance, bidirectional
- ❌ Manual OTEL, scaling complexity
- **Use case**: Real-time battle communication (strong candidate)

**SSE:**
- ✅ Auto-reconnect, simple server-push
- ❌ Unidirectional, requires POST for bot actions
- **Use case**: Asymmetric scenarios (server-heavy updates)

## References

- [WebSocket Protocol (RFC 6455)](https://tools.ietf.org/html/rfc6455)
- [Server-Sent Events Specification](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- [HTTP/2 Specification (RFC 7540)](https://tools.ietf.org/html/rfc7540)
- [MessagePack Specification](https://msgpack.org/)
- [CBOR (RFC 8949)](https://tools.ietf.org/html/rfc8949)
- [OpenTelemetry HTTP Instrumentation](https://opentelemetry.io/docs/instrumentation/)
- [WebSocket vs Server-Sent Events](https://ably.com/topic/websockets-vs-sse)
- [gorilla/websocket - Go WebSocket library](https://github.com/gorilla/websocket)
