---
title: "Custom TCP/UDP Protocol Analysis"
description: >
    Analysis of custom TCP and UDP protocols for bot-to-server communication
type: docs
weight: 3
date: 2025-12-05
---

## Overview

Custom TCP and UDP protocols offer maximum control and potential performance optimization for game networking. This analysis evaluates whether the complexity and development cost of custom protocols is justified for Battle Bots, drawing on lessons from game industry networking.

## TCP Custom Protocol Analysis

### Protocol Design

**Message Format Example:**
```
[4 bytes: magic number 0x42544C00]
[2 bytes: protocol version]
[2 bytes: message type]
[4 bytes: message length]
[4 bytes: sequence number]
[N bytes: payload (JSON, MessagePack, Protobuf)]
[4 bytes: CRC32 checksum]
```

**Connection-Oriented Characteristics:**
- Reliable, ordered delivery (TCP guarantees)
- Stream-based (requires message framing)
- Connection handshake and teardown
- Flow control and congestion avoidance built-in

### Message Framing Strategies

**Length-Prefixed:**
```
[4-byte length][payload]
```
- Simple to parse
- Requires buffering until full message received

**Delimiter-Based:**
```
[payload]\n
```
- Easy to implement
- Inefficient if payload contains delimiters

**Fixed-Length Headers:**
- Header describes payload structure
- Best practice for custom protocols
- Allows for protocol evolution

### Performance Characteristics

**Latency:**
- **5-15ms** round-trip (local network)
- Similar to gRPC/WebSocket (all use TCP)
- No inherent advantage over HTTP/2 or WebSocket

**Throughput:**
- Limited by TCP, not protocol overhead
- Custom framing slightly more efficient than HTTP
- Marginal gain (5-10%) vs WebSocket

**CPU Overhead:**
- Custom parsing faster than HTTP header parsing
- Binary protocols minimize serialization cost
- Negligible difference in practice for Battle Bots scale

### Advantages

✅ **Maximum control** over message format  
✅ **Minimal overhead** (no HTTP headers)  
✅ **Custom flow control** if needed  
✅ **Protocol optimizations** for specific use case  

### Disadvantages

❌ **High implementation cost** (build, test, debug)  
❌ **Language-specific implementations** (no code generation)  
❌ **Limited tooling** (manual packet inspection)  
❌ **No standard libraries** (roll your own)  
❌ **Bug-prone** (protocol state machines, edge cases)  
❌ **No advantage over WebSocket/gRPC** for this use case  

### Comparison with WebSocket

Both use TCP, so performance is nearly identical:
- WebSocket has ~6 bytes per frame overhead
- Custom protocol saves ~6 bytes per message
- **Conclusion**: Negligible benefit, massive implementation cost

## UDP Custom Protocol Analysis

### Protocol Design

**Unreliable Datagram Characteristics:**
- No connection (connectionless)
- Packets can be lost, duplicated, reordered
- Low latency (no TCP handshake/ACK overhead)
- Manual reliability layer required

**Custom Reliability Layer Example:**
```
[4 bytes: sequence number]
[4 bytes: ack bitfield]
[1 byte: reliability flags]
[N bytes: payload]
```

**Reliability Modes:**
- **Unreliable**: Send and forget (position updates)
- **Reliable**: Retransmit until ACK (critical events)
- **Ordered**: Sequence numbers, discard out-of-order
- **Sequenced**: Latest only, discard older packets

### Performance Characteristics

**Latency:**
- **2-10ms** round-trip (no TCP handshake)
- 30-50% lower latency than TCP in ideal conditions
- Worse under packet loss (retransmissions)

**Throughput:**
- Higher ceiling than TCP (no congestion control)
- Prone to network congestion without custom control
- Better for bursty traffic

**Packet Loss Handling:**
- 1-5% packet loss is common on Internet
- Game engines interpolate/extrapolate positions
- Critical events need reliability layer

### Use Cases in Game Networking

**Good for UDP:**
- **Player position updates** (frequent, latest is best)
- **Input state** (redundant, lose older packets)
- **Non-critical effects** (particle systems, audio)
- **High tick-rate simulations** (60+ updates/sec)

**Bad for UDP:**
- **Player actions** (attack, use item - must be reliable)
- **Inventory changes** (critical state)
- **Chat messages** (must arrive)
- **Battle Bots actions** (move, attack decisions are critical)

### Game Networking Patterns

**Client-Side Prediction:**
```
Client predicts local movement → Send input to server
Server authoritative simulation → Send correction
Client reconciles prediction with server state
```
- Masks network latency
- Requires complex reconciliation logic

**Server Reconciliation:**
- Server is authoritative
- Clients show approximate state
- Server corrects divergence

**Snapshot Interpolation:**
- Server sends state snapshots
- Client interpolates between snapshots
- Smooth animation despite packet loss

**Delta Compression:**
- Send only changed fields
- Reduces bandwidth for large state objects
- Requires baseline tracking

**Are these needed for Battle Bots?**
- Turn-based or low tick-rate: **No**
- Real-time FPS-style: **Maybe**
- Current POC requirements: **No**

### OpenTelemetry Integration Challenges

**Major Issues:**
- **No standard trace propagation** (no HTTP headers)
- **Custom metadata format** required
- **Manual span creation** for every packet
- **Sampling complexity** (trace 1/1000 packets?)

**Implementation Approach:**
```go
// Custom trace context in UDP packet
type PacketHeader struct {
    TraceID  [16]byte
    SpanID   [8]byte
    Flags    byte
    // ... rest of packet
}
```

**Metrics Collection:**
- Packet send/receive counters
- Loss rate calculation
- Latency histogram (manual timing)
- Battle-specific metrics

**Complexity vs gRPC/WebSocket:**
- 10x more complex to instrument
- No automatic integration with OTLP
- High maintenance burden

### Container Networking

**Port Mapping:**
- UDP ports must be explicitly mapped
- Each bot needs unique port (or shared port with routing)

**Firewall Traversal:**
- Many networks block UDP
- Corporate firewalls often UDP-hostile
- NAT64 environments problematic

**Load Balancing:**
- Stateless load balancing possible
- Consistent hashing needed for session affinity
- More complex than TCP

**NAT Traversal for P2P:**

**Hole Punching:**
1. Both bots connect to rendezvous server
2. Server shares IP:port of each bot
3. Bots send UDP packets to each other
4. NAT creates temporary mappings
5. Direct P2P communication established

**STUN/TURN:**
- STUN: Discover public IP:port
- TURN: Relay server for firewall traversal
- TURN required when symmetric NAT prevents hole-punching

**Complexity:**
- Significantly more complex than TCP
- Failure rate: 5-20% of connections may require TURN
- Additional infrastructure cost

### Development & Tooling

**Implementation Cost:**
- **High**: Build protocol from scratch in each language
- **Testing**: Unit tests, integration tests, fuzz testing
- **Edge cases**: Connection state, packet loss, reordering

**Debugging Tools:**
- **tcpdump**: Capture packets
- **Wireshark**: Decode custom protocol (requires dissector plugin)
- **Custom tooling**: Packet replayer, simulator
- **Logging**: Extensive logging required

**Language Support:**
- Manual implementation per language
- No code generation (unlike gRPC/Protobuf)
- Consistency challenges across languages

**Maintenance:**
- Protocol versioning complexity
- Backward compatibility testing
- Documentation burden for bot developers

### Industry Examples & Lessons

#### Unreal Engine Networking

**Architecture:**
- UDP-based with custom reliability layer
- RPCs for critical events (reliable)
- Property replication for state (unreliable with delta compression)
- Client-side prediction and server reconciliation

**Why UDP:**
- Fast-paced FPS games (60-120 tick rate)
- Hundreds of position updates per second
- Packet loss acceptable for intermediate positions

**Applicability to Battle Bots:**
- Unclear if Battle Bots needs 60+ tick rate
- If turn-based, UDP advantage disappears

#### Unity DOTS NetCode

**Architecture:**
- UDP with custom reliable channel
- Delta compression for snapshot synchronization
- Ghost entities (networked objects)
- Client-side prediction

**Why UDP:**
- Real-time multiplayer games
- Low latency critical
- High entity count

#### Valve Source Engine

**Architecture:**
- UDP with separate reliable/unreliable channels
- Lag compensation for hit detection
- Client-side prediction
- Command acknowledgment

**Lessons:**
- UDP justified for < 50ms latency requirements
- Complex to implement correctly
- Years of refinement

#### What Battle Bots Can Learn

**Question**: Does Battle Bots need < 50ms latency?

**Turn-based or slow tick-rate (< 10/sec):**
- TCP is sufficient (WebSocket, gRPC)
- UDP complexity not justified

**Real-time fast-paced (> 30/sec tick rate):**
- UDP may provide benefits
- Requires client-side prediction
- High development cost

**POC Requirement:**
- No specific tick rate mentioned
- Suggests UDP is premature optimization

## Protocol Design Considerations

### Message Format Design

**Efficiency:**
- Binary > Text for bandwidth
- Protocol Buffers balance efficiency + maintainability

**Versioning:**
- Protocol version field in header
- Feature negotiation during handshake
- Backward compatibility strategy

**Extensibility:**
- Reserve fields for future use
- TLV (Type-Length-Value) encoding
- Protocol Buffers automatically extensible

### Error Handling

**Connection Errors:**
- Timeout handling
- Reconnection logic
- State recovery

**Protocol Errors:**
- Invalid message format
- Unexpected message type
- Version mismatch

**Application Errors:**
- Invalid bot actions
- Game rule violations
- Rate limiting

### Authentication & Security

**Connection Authentication:**
- Token-based (JWT in handshake)
- Mutual TLS (mTLS)
- API keys

**Message Integrity:**
- HMAC signatures
- CRC/checksum for corruption
- Encryption (TLS for TCP, DTLS for UDP)

**Denial of Service:**
- Rate limiting per bot
- Connection limits
- Packet flood protection

## Development Complexity Comparison

| Aspect | Custom TCP | Custom UDP | gRPC | WebSocket |
|--------|-----------|-----------|------|-----------|
| Implementation | High | Very High | Low | Low |
| Language Support | Manual | Manual | Generated | Libraries |
| Tooling | Minimal | Minimal | Excellent | Good |
| OTEL Integration | Hard | Very Hard | Native | Moderate |
| Maintenance | High | Very High | Low | Low |
| Debugging | Hard | Very Hard | Easy | Moderate |

**Effort Estimate:**
- Custom TCP: 4-6 weeks per language
- Custom UDP: 8-12 weeks per language (+ reliability layer)
- gRPC: 1-2 days (define .proto, generate code)
- WebSocket: 1-3 days (use library)

## Pros and Cons Summary

### Custom TCP

✅ Maximum protocol control  
✅ Minimal wire overhead  
❌ High implementation cost (weeks per language)  
❌ No tooling ecosystem  
❌ Difficult OTEL integration  
❌ **No performance advantage over WebSocket**  

**Verdict**: ❌ Not justified for Battle Bots

### Custom UDP

✅ Lowest latency (2-10ms)  
✅ No TCP overhead  
✅ Good for high tick-rate (60+/sec)  
❌ Very high implementation cost (months)  
❌ Reliability layer complexity  
❌ Poor firewall/NAT traversal  
❌ Very difficult OTEL integration  
❌ **Only beneficial if < 50ms latency critical**  

**Verdict**: ❌ Premature optimization for POC

### When Custom Protocols Make Sense

**Justified:**
- Extreme performance requirements (< 10ms latency)
- Specialized hardware or protocols
- Unique constraints not met by standards

**Not Justified (Battle Bots POC):**
- Unknown tick rate requirements
- Observability is priority (ADR-0002)
- Language-agnostic bot interface
- Development velocity matters

## Suitability for Battle Bots

**Custom TCP**: ⭐☆☆☆☆ (Not Recommended)
- No performance benefit over WebSocket
- Massive implementation cost
- Poor observability integration

**Custom UDP**: ⭐☆☆☆☆ (Not Recommended for POC)
- Unclear if latency requirements justify complexity
- OTEL integration extremely difficult
- Firewall/NAT traversal issues
- Defer until POC proves need

**Recommendation**: Use standard protocols (gRPC/WebSocket) for POC. Re-evaluate custom UDP only if profiling shows < 50ms latency is critical and unachievable with TCP-based protocols.

## References

- [Game Networking Patterns (Gaffer On Games)](https://gafferongames.com/categories/game-networking/)
- [Source Multiplayer Networking](https://developer.valvesoftware.com/wiki/Source_Multiplayer_Networking)
- [Unreal Engine Network Architecture](https://docs.unrealengine.com/5.0/en-US/networking-overview-for-unreal-engine/)
- [Unity DOTS NetCode](https://docs.unity3d.com/Packages/com.unity.netcode@latest)
- [Fast-Paced Multiplayer (Gabriel Gambetta)](https://www.gabrielgambetta.com/client-server-game-architecture.html)
- [NAT Traversal Techniques](https://tools.ietf.org/html/rfc5128)
- [STUN Protocol (RFC 5389)](https://tools.ietf.org/html/rfc5389)
- [TURN Protocol (RFC 5766)](https://tools.ietf.org/html/rfc5766)
- [TCP vs UDP for Games (Glenn Fiedler)](https://gafferongames.com/post/udp_vs_tcp/)
