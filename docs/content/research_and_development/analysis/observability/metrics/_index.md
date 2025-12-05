---
title: "Metrics Storage"
description: >
    Analysis of metrics storage backends for the BattleBots observability stack, focusing on systems that integrate with the OpenTelemetry Collector.
type: docs
---

## Overview

Metrics storage is essential for the BattleBots platform's observability infrastructure, enabling:

- Real-time monitoring of battle events and game state
- Historical analysis of bot performance and system behavior
- Capacity planning and infrastructure optimization
- Alerting on critical system conditions
- Long-term trend analysis and reporting

The metrics storage backend must handle time-series data at scale, support efficient querying, and integrate seamlessly with the OpenTelemetry Collector to provide a unified observability platform alongside logs and traces.

## Why Metrics Storage Matters for BattleBots

The BattleBots platform generates metrics across multiple dimensions:

### Game Metrics
- **Battle Events**: Damage calculations, bot actions, victory conditions
- **Performance Metrics**: Bot response times, action execution latency
- **Game State**: Active battles, queued matches, player counts
- **Resource Utilization**: CPU, memory, and network usage per bot container

### Infrastructure Metrics
- **Container Metrics**: Resource usage, restart counts, health checks
- **Host Metrics**: Node-level CPU, memory, disk, and network utilization
- **Network Metrics**: Request rates, latency distributions, error rates
- **Kubernetes Metrics**: Pod status, deployments, scaling events

### Observability Stack Metrics
- **OpenTelemetry Collector**: Pipeline throughput, batch sizes, export success rates
- **Log Storage**: Ingestion rates, query performance, storage utilization
- **Trace Storage**: Span ingestion, sampling rates, trace completeness

## Components

### [Grafana Mimir](mimir/)

Research on Grafana Mimir, a horizontally scalable, highly available, multi-tenant metrics storage system built for long-term Prometheus data storage.

Mimir transforms Prometheus from a single-server monitoring system into a distributed platform capable of handling over 1 billion active time series, providing:

- **Native OTLP Support**: Direct integration with OpenTelemetry Collector via OTLP over HTTP
- **Horizontal Scalability**: Independent scaling of write path, read path, and backend components
- **Long-Term Storage**: Object storage backend (S3, GCS, MinIO) enables months to years of retention
- **Multi-Tenancy**: Built-in tenant isolation with per-tenant limits and resource controls
- **High Availability**: Replication and distributed architecture eliminate single points of failure
- **PromQL Compatibility**: Full Prometheus query language support for dashboards and alerts

Includes detailed analysis of:
- Architecture components and deployment modes
- Native OTLP ingestion and OpenTelemetry Collector integration
- Object storage backends and retention policies
- Multi-tenancy and cardinality management
- Comparison with Prometheus, Thanos, and Cortex
- Production deployment and operational best practices

## BattleBots Integration Context

For the BattleBots platform, metrics storage serves as the foundation for understanding system behavior and performance:

### Real-Time Monitoring
- Monitor active battles and player engagement in real-time
- Alert on critical conditions (bot container failures, API errors, resource exhaustion)
- Track game server health and availability
- Identify performance degradation before user impact

### Historical Analysis
- Analyze battle outcome patterns and bot performance trends
- Capacity planning based on player growth and peak usage patterns
- Cost optimization through resource utilization analysis
- Root cause analysis for system incidents

### Unified Observability
- Metric-to-trace correlation through shared labels and exemplars
- Jumping from metric anomalies to related distributed traces
- Linking metrics to logs for comprehensive debugging
- Grafana dashboards combining metrics, logs, and traces in single views

### Example Metrics for BattleBots

**Bot Performance**:
```promql
# Bot action latency by bot type
histogram_quantile(0.95,
  sum(rate(bot_action_duration_seconds_bucket{action="attack"}[5m])) by (bot_type, le)
)

# Bot health over time
avg(bot_health_points) by (battle_id, bot_id)
```

**Infrastructure Health**:
```promql
# Container resource utilization
container_memory_usage_bytes{namespace="battlebots", container="game-server"}
  /
container_spec_memory_limit_bytes{namespace="battlebots", container="game-server"}

# API request rate and errors
sum(rate(http_requests_total{service="battle-api"}[5m])) by (status_code)
```

**Game State**:
```promql
# Active battles
sum(battles_active{environment="production"})

# Player queue depth
avg(player_queue_length) by (region)
```

## Decision Context

This research will inform the upcoming **ADR-NNNN: Observability Stack Selection**, which will determine the metrics storage backend for BattleBots. Key decision factors include:

- **Functional Fit**: Does the solution meet metrics storage, query, and correlation requirements?
- **Scalability**: Can it handle expected growth from POC to production scale?
- **Integration**: How well does it integrate with OpenTelemetry Collector, Grafana, and other observability components?
- **Operational Complexity**: What is the operational burden for deployment, monitoring, and maintenance?
- **Cost**: What are the infrastructure and operational costs at POC and production scale?
- **Multi-Tenancy**: Does it support per-player or per-battle isolation if needed?

## Related Documentation

### R&D Documentation

- [Observability Overview](../) - Parent observability analysis section
- [OpenTelemetry Collector Analysis](../otel-collector/) - Metrics collection and processing
- [Loki Analysis](../logs/loki/) - Log storage for correlation with metrics
- [User Journey 0001: POC](../../user-journeys/0001-poc.md) - Observability requirements context
- Future ADRs on observability stack architecture

### External Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [PromQL Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [OpenTelemetry Metrics Specification](https://opentelemetry.io/docs/specs/otel/metrics/)
- [Grafana Metrics Documentation](https://grafana.com/docs/grafana/latest/fundamentals/intro-to-metrics/)

## Contributing

These analysis documents are living documents that should be updated as:
- New metrics storage solutions emerge or mature
- BattleBots observability requirements evolve
- Team members gain operational experience with metrics backends
- Best practices and patterns are discovered
- Comparative analysis reveals new insights

Updates should maintain the high-level overview focus with links to authoritative sources for technical deep-dives.
