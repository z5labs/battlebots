---
title: "Protocols"
description: >
    Analysis of communication protocols for Battle Bots bot-to-server interface
type: docs
weight: 100
---

This section contains detailed technical analysis of communication protocol options evaluated for the Battle Bots bot-to-server and bot-to-bot interface.

## Protocol Evaluations

- **[gRPC](grpc/grpc-analysis/)** - HTTP/2-based RPC with Protocol Buffers and bidirectional streaming
- **[HTTP-based Protocols](http/http-analysis/)** - REST, WebSockets, and Server-Sent Events analysis
- **[Custom TCP/UDP](custom/tcp-udp-analysis/)** - Low-level custom protocol evaluation

These analyses inform [ADR-0004: Bot to Battle Server Communication Protocol](../../adrs/0004-bot-battle-server-interface/).
