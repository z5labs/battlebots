---
title: "Grafana Tempo"
description: >
    Research and analysis of Grafana Tempo distributed tracing backend for the BattleBots observability stack.
type: docs
---

## Overview

Grafana Tempo is a high-volume, minimal dependency distributed tracing backend designed for cost-efficiency and operational simplicity. Unlike traditional tracing systems that require complex database infrastructure, Tempo uses object storage as its only dependency, dramatically reducing operational complexity while providing powerful trace querying through TraceQL.

This analysis explores Tempo as the distributed tracing backend for the BattleBots platform, focusing on its architecture, native OTLP support, OpenTelemetry Collector integration, and integration with the broader Grafana observability stack (Loki, Mimir, Prometheus).

## Why Research Tempo?

For the BattleBots observability stack, Tempo offers several compelling advantages:

- **Native OTLP Support**: Full support for OTLP/gRPC (port 4317) and OTLP/HTTP (port 4318) protocols, enabling seamless integration with the OpenTelemetry Collector
- **Cost-Effective Storage**: Object storage-only design (S3, GCS, Azure, MinIO) eliminates expensive database infrastructure and reduces storage costs by 10x or more compared to Jaeger or Zipkin
- **Minimal Dependencies**: No Cassandra, Elasticsearch, or other complex databases required—only object storage
- **Horizontal Scalability**: Microservices architecture supports scaling from development to production workloads handling millions of spans per second
- **TraceQL Query Language**: Powerful SQL-like query language for filtering and analyzing traces by span attributes, duration, and relationships
- **Grafana Integration**: Seamless correlation with metrics (Mimir/Prometheus) and logs (Loki) through exemplars and trace IDs for unified observability
- **Multi-Protocol Support**: Accepts OTLP, Jaeger, Zipkin, and OpenCensus formats, enabling gradual migration from legacy tracing systems

## Document Structure

The Tempo analysis is organized into the following documents:

### [Tempo Overview](tempo-overview.md)

Comprehensive overview covering architecture, deployment modes, and operational best practices.

**Topics covered:**
- What is Tempo and its design philosophy (index-free, object storage-based tracing)
- Core concepts: traces, spans, blocks storage, TraceQL, sampling strategies
- Architecture components: distributor, ingester, querier, query-frontend, compactor, metrics-generator
- Deployment modes: monolithic, scalable, microservices (with comparison table)
- How to run Tempo with Docker Compose for POC environments
- Best practices for sampling, storage optimization, and performance tuning
- When to use Tempo vs. Jaeger vs. Zipkin (decision criteria matrix)
- Resource requirements and capacity planning guidance
- Complete configuration examples with Grafana data source setup

**Audience**: Everyone—provides foundational understanding for evaluating Tempo as the tracing backend.

### [OTLP Integration](tempo-otlp-integration.md)

Deep dive into OTLP compatibility and OpenTelemetry Collector integration (addresses critical user requirements).

**Topics covered:**
- Native OTLP support in Tempo (status: YES since v1.3.0+, endpoint configuration)
- OTLP/gRPC and OTLP/HTTP protocol support and when to use each
- OpenTelemetry Collector otlp and otlphttp exporter configuration for Tempo
- Batch processor, retry policies, and queue management best practices
- Resource attribute mapping and span attribute strategies
- Sampling configuration (head-based, tail-based) in the OTel Collector
- Authentication and multi-tenancy setup with X-Scope-OrgID headers
- Complete working configuration examples (OTel Collector + Tempo + Grafana)
- Troubleshooting common integration issues (connection errors, missing traces)
- BattleBots-specific integration patterns and TraceQL query examples

**Audience**: Developers and operators implementing the OpenTelemetry Collector to Tempo pipeline.

## BattleBots Integration Context

For the BattleBots platform, Tempo would serve as the centralized distributed tracing backend, receiving traces from the OpenTelemetry Collector via OTLP. This enables:

### Battle Workflow Tracing

- **Complete Battle Flow**: Trace entire battle lifecycle from matchmaking → battle initialization → game loop execution → victory condition → results persistence
- **Bot Action Traces**: Track individual bot actions (move, attack, defend) with timing, success/failure, and resource costs
- **State Transitions**: Capture game state changes with span events marking key transitions (battle start, bot death, timer expiration)
- **Latency Analysis**: Identify performance bottlenecks in action processing, state synchronization, and event broadcasting

### Request Flow Tracing

- **API Requests**: Trace HTTP/gRPC requests from client → API gateway → game service → persistence layer
- **WebSocket Connections**: Track WebSocket connection lifecycle, authentication, and message flow
- **Service Interactions**: Visualize calls between microservices (matchmaking service → game server manager → bot runtime)
- **External Dependencies**: Monitor calls to external services (authentication, leaderboards, analytics)

### Error Investigation

- **Exception Tracking**: Capture stack traces and error context through span events
- **Failure Propagation**: Trace error propagation across service boundaries to identify root causes
- **Timeout Analysis**: Identify services or operations exceeding latency SLAs
- **Retry Patterns**: Visualize retry behavior and backoff strategies

### Unified Observability

- **Trace-to-Logs Correlation**: Link traces to logs via TraceID/SpanID for detailed debugging context
- **Trace-to-Metrics Correlation**: Use exemplars to jump from metric spikes to example slow traces
- **Grafana Dashboards**: Combine traces, metrics (Mimir), and logs (Loki) in unified visualizations
- **Root Cause Analysis**: Start with metric alert → find exemplar trace → examine correlated logs

### Performance Optimization

- **Latency Profiling**: Identify slow database queries, external API calls, or computational bottlenecks
- **Resource Utilization**: Correlate trace duration with CPU/memory metrics for capacity planning
- **Sampling Strategies**: Implement tail sampling to capture all errors and slow requests while reducing trace volume
- **TraceQL Analysis**: Query for patterns like "all traces with database latency > 100ms" or "errors in bot action processing"

## Decision Context

This research will inform the upcoming **ADR-NNNN: Observability Stack Selection**, which will determine the distributed tracing backend for BattleBots. Key decision factors specific to Tempo include:

- **Cost Efficiency**: Does Tempo's object storage-only design provide sufficient cost savings to justify its adoption over traditional tracing backends?
- **OTLP Integration**: Does native OTLP support simplify the OpenTelemetry Collector integration compared to alternatives requiring protocol translation?
- **Search Trade-offs**: Is Tempo's trace-ID-first query model (with TraceQL for advanced queries) sufficient for BattleBots debugging workflows?
- **Operational Complexity**: Can the team operate Tempo with minimal infrastructure overhead compared to Jaeger (Cassandra/Elasticsearch) or Zipkin?
- **Stack Integration**: Does Tempo's deep integration with Grafana, Loki, and Mimir provide enough value to justify adopting the full LGTM stack?
- **Scalability**: Can Tempo handle growth from POC (thousands of traces) to production (potentially millions of traces per day)?

The ADR will also consider alternative approaches:
- **Jaeger**: More mature with powerful tag-based search, but requires Cassandra or Elasticsearch and higher operational complexity
- **Zipkin**: Simpler but less scalable, requires database backend, limited query capabilities
- **Elastic APM**: Unified observability in Elasticsearch stack, but higher cost and complexity
- **Managed Services**: Grafana Cloud Tempo, Honeycomb, or Lightstep for reduced operational burden

## Related Observability Components

Tempo completes the three-signal observability architecture alongside:
- [**Grafana Loki**](../logs/loki/) - Log aggregation and storage
- [**Grafana Mimir**](../metrics/mimir/) - Metrics storage and querying
- [**OpenTelemetry Collector**](../otel-collector/) - Telemetry data pipeline

Together, these components form the LGTM stack (Loki, Grafana, Tempo, Mimir), providing unified observability with trace-metric-log correlation through exemplars and shared trace IDs.

## External Resources

- [Grafana Tempo Official Documentation](https://grafana.com/docs/tempo/latest/)
- [Tempo GitHub Repository](https://github.com/grafana/tempo)
- [TraceQL Query Language Reference](https://grafana.com/docs/tempo/latest/traceql/)
- [Grafana Tempo Blog](https://grafana.com/blog/tag/tempo/)
- [OpenTelemetry to Tempo Integration Guide](https://grafana.com/docs/tempo/latest/configuration/)
- [Tempo Architecture Deep Dive](https://grafana.com/docs/tempo/latest/operations/architecture/)
- [LGTM Stack Overview](https://grafana.com/go/webinar/intro-to-mltp/)
