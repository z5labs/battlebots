---
title: "Grafana Loki"
description: >
    Research and analysis of Grafana Loki log aggregation system for the BattleBots observability stack.
type: docs
---

## Overview

Grafana Loki is a horizontally scalable, highly available, multi-tenant log aggregation system inspired by Prometheus. Unlike other log aggregation systems, Loki is designed to be cost-effective and easy to operate by not indexing the contents of logs, but rather a set of labels for each log stream.

This analysis explores Loki as a potential log storage backend for the BattleBots platform, focusing on its architecture, deployment options, OTLP compatibility, and integration with the OpenTelemetry Collector.

## Why Research Loki?

For the BattleBots observability stack, Loki offers several compelling advantages:

- **Native OTLP Support**: Loki v3+ provides native OTLP ingestion endpoints, enabling seamless integration with the OpenTelemetry Collector
- **Cost-Effective Storage**: Index-free approach dramatically reduces storage costs compared to full-text indexing systems
- **Horizontal Scalability**: Microservices architecture supports scaling from development to production workloads
- **Grafana Integration**: Tight integration with Grafana provides unified visualization of logs, metrics, and traces
- **Cloud-Native Design**: Built for containerized environments with Kubernetes-first deployment patterns

## Document Structure

The Loki analysis is organized into the following documents:

### [Loki Overview](loki-overview.md)

Comprehensive overview covering architecture, deployment, and operational best practices.

**Topics covered:**
- What is Loki and its design philosophy (index-free, label-based querying)
- Core concepts: streams, labels, chunks, index, LogQL
- Architecture components: distributor, ingester, querier, compactor
- Deployment modes: monolithic, simple scalable, microservices
- How to run Loki with Docker/Podman Compose for POC
- Best practices for label strategy and configuration
- When to use Loki vs. alternatives like Elasticsearch

**Audience**: Everyone—provides foundational understanding for evaluating Loki as a log backend.

### [OTLP Integration](loki-otlp-integration.md)

Deep dive into OTLP compatibility and OpenTelemetry Collector integration.

**Topics covered:**
- Native OTLP support in Loki v3+ (endpoints, configuration)
- OTel Collector otlphttp exporter setup
- Resource attribute mapping to Loki labels
- Log-trace correlation via TraceID/SpanID
- Complete working configuration examples
- Authentication and multi-tenancy
- Troubleshooting common integration issues

**Audience**: Developers and operators implementing the OTel Collector to Loki pipeline.

## BattleBots Integration Context

For the BattleBots platform, Loki would serve as the centralized log storage backend, receiving logs from the OpenTelemetry Collector via OTLP. This enables:

### Game Event Logging
- Battle events (bot actions, damage calculations, victory conditions)
- Game state transitions and timing information
- Player actions and command processing
- Error conditions and system anomalies

### Infrastructure Logging
- Container logs from game servers
- System logs from host infrastructure
- Application logs from Go services
- Network and security logs

### Unified Observability
- Log-trace correlation through shared TraceID/SpanID
- Linking log events to metrics and distributed traces
- Grafana dashboards combining logs, metrics, and traces
- Streamlined debugging workflows across all telemetry signals

## Decision Context

This research will inform the upcoming **ADR-NNNN: Observability Stack Selection**, which will determine the log storage backend for BattleBots. Key decision factors include:

- **Functional fit**: Does Loki meet log storage, query, and correlation requirements?
- **Operational complexity**: How difficult is it to deploy, monitor, and maintain Loki?
- **Cost**: What are the infrastructure and operational costs at POC and production scale?
- **Integration**: How well does it integrate with OTel Collector, Grafana, and the broader observability stack?

## External Resources

- [Grafana Loki Official Documentation](https://grafana.com/docs/loki/latest/)
- [Loki GitHub Repository](https://github.com/grafana/loki)
- [Grafana Labs Blog](https://grafana.com/blog/)
- [CNCF Loki Project](https://www.cncf.io/)
