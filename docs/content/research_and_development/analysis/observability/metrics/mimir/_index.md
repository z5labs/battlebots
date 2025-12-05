---
title: "Grafana Mimir"
description: >
    Research and analysis of Grafana Mimir metrics storage system for the BattleBots observability stack.
type: docs
---

## Overview

Grafana Mimir is a horizontally scalable, highly available, multi-tenant, long-term storage solution for Prometheus metrics. It transforms Prometheus's single-server architecture into a distributed microservices platform capable of handling over 1 billion active time series with unlimited retention.

This analysis explores Mimir as the metrics storage backend for the BattleBots platform, focusing on its architecture, native OTLP support, OpenTelemetry Collector integration, and operational characteristics.

## Why Research Mimir?

For the BattleBots observability stack, Mimir offers several compelling advantages over standalone Prometheus:

- **Native OTLP Support**: Direct ingestion via `/otlp/v1/metrics` endpoint since version 2.3.0, enabling seamless OpenTelemetry Collector integration without protocol translation
- **Massive Scalability**: Proven at 1 billion active time series (Grafana Labs internal testing) and 500 million series in production customer deployments
- **Long-Term Storage**: Object storage backend (S3, GCS, MinIO) enables months to years of retention at minimal cost compared to local disk storage
- **Multi-Tenancy**: Built-in tenant isolation with per-tenant limits, enabling future use cases like per-player metrics or per-battle analytics
- **High Availability**: 3x replication by default, distributed architecture, and automatic failover eliminate single points of failure
- **Prometheus Compatibility**: Full PromQL support ensures existing Prometheus queries, dashboards, and alerts work unchanged

## Document Structure

The Mimir analysis is organized into the following documents:

### [Mimir Overview](mimir-overview.md)

Comprehensive overview covering architecture, deployment modes, storage backends, and operational best practices.

**Topics covered:**
- What is Mimir and its design philosophy (distributed Prometheus with object storage)
- Core concepts: blocks storage, time series, multi-tenancy, cardinality management
- Architecture components: distributor, ingester, querier, query-frontend, store-gateway, compactor, ruler
- Deployment modes: monolithic, read-write, microservices (with comparison table)
- How to run Mimir with Docker Compose and MinIO for POC environments
- Production deployment patterns with Kubernetes and Helm
- Best practices for label strategy, configuration, storage selection, and performance tuning
- When to use Mimir vs. Prometheus vs. Thanos (decision criteria matrix)
- Resource requirements and capacity planning guidance
- Complete configuration examples for all deployment modes

**Audience**: Everyone—provides foundational understanding for evaluating Mimir as the metrics backend.

### [OTLP Integration](mimir-otlp-integration.md)

Deep dive into OTLP compatibility and OpenTelemetry Collector integration (addresses critical user requirements).

**Topics covered:**
- Native OTLP support in Mimir (status, version requirements, endpoint configuration)
- OTLP vs. Prometheus remote write comparison and recommendations
- OpenTelemetry Collector otlphttp exporter configuration for Mimir
- Alternative: OpenTelemetry Collector prometheusremotewrite exporter setup
- Batch processor, retry policies, and queue management best practices
- Resource attribute mapping from OTel to Mimir labels
- Label strategy and cardinality control for OTel-generated metrics
- Authentication and multi-tenancy setup with X-Scope-OrgID headers
- Complete working configuration examples (OTel Collector + Mimir + Grafana)
- Troubleshooting common integration issues (connection errors, cardinality, performance)
- BattleBots-specific integration patterns and example queries

**Audience**: Developers and operators implementing the OpenTelemetry Collector to Mimir pipeline.

## BattleBots Integration Context

For the BattleBots platform, Mimir would serve as the centralized metrics storage backend, receiving metrics from the OpenTelemetry Collector via native OTLP ingestion. This enables:

### Game Metrics Storage

- **Bot Performance**: Action latency, damage calculations, movement speed, resource utilization per bot
- **Battle Events**: Start/end times, player actions, victory conditions, matchmaking metrics
- **Game State**: Active battles count, queued players, concurrent users, session durations
- **Quality Metrics**: Frame rates, tick rates, network latency, synchronization quality

### Infrastructure Metrics

- **Container Metrics**: CPU/memory usage, restart counts, health checks for bot containers and game servers
- **Kubernetes Metrics**: Pod status, node utilization, deployment health, scaling events
- **Network Metrics**: Request rates, latency distributions, error rates, bandwidth consumption
- **Host Metrics**: System-level CPU, memory, disk I/O, network traffic across infrastructure

### Observability Stack Metrics

- **OpenTelemetry Collector**: Pipeline throughput, batch sizes, queue depths, export success/failure rates
- **Loki (Log Storage)**: Log ingestion rates, query latency, storage utilization
- **Tempo (Trace Storage)**: Span ingestion, trace completeness, sampling rates
- **Mimir Self-Monitoring**: Ingester series counts, query performance, compaction status, object storage health

### Long-Term Analytics

- **Capacity Planning**: Historical resource usage trends to predict scaling needs
- **Cost Optimization**: Identify underutilized resources and optimize allocation
- **Performance Baselines**: Establish normal behavior patterns for anomaly detection
- **Business Intelligence**: Player engagement metrics, battle frequency, peak usage times

## Decision Context

This research will inform the upcoming **ADR-NNNN: Observability Stack Selection**, which will determine the metrics storage backend for BattleBots. Key decision factors specific to Mimir include:

- **Scalability Requirements**: Can Mimir handle expected growth from POC (thousands of series) to production (millions+ of series)?
- **OTLP Integration**: Does native OTLP support simplify the OpenTelemetry Collector integration compared to alternatives?
- **Operational Complexity**: Is the team prepared to operate a distributed metrics system, or should we start with standalone Prometheus?
- **Cost vs. Value**: Do Mimir's features (long-term storage, scalability, multi-tenancy) justify the increased infrastructure cost and complexity?
- **Migration Path**: If starting with Prometheus, how difficult is migration to Mimir when scale demands it?

The ADR will also consider alternative approaches:
- **Standalone Prometheus**: Simpler but limited to ~10M series and short retention
- **Thanos**: Similar capabilities to Mimir with different architectural trade-offs
- **Managed Services**: Grafana Cloud or other hosted Prometheus-compatible solutions

## External Resources

- [Grafana Mimir Official Documentation](https://grafana.com/docs/mimir/latest/)
- [Mimir GitHub Repository](https://github.com/grafana/mimir)
- [Grafana Mimir Blog](https://grafana.com/blog/tag/mimir/)
- [Mimir Capacity Calculator](https://o11y.tools/mimircalc/)
- [OpenTelemetry Metrics to Mimir Guide](https://grafana.com/docs/mimir/latest/configure/configure-otel-collector/)
- [Prometheus Remote Write Specification](https://prometheus.io/docs/concepts/remote_write_spec/)
