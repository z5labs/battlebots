---
title: "ADR-0003: Observability Stack Architecture"
linkTitle: "ADR-0003: Observability Stack"
weight: 3
date: 2025-12-05
status: "proposed"
category: "strategic"
deciders: []
consulted: []
informed: []
---

# Observability Stack Architecture

## Context and Problem Statement

Following the decision to adopt OpenTelemetry as our observability SDK (ADR-0002), we need to select the backend components for our observability stack. The stack must handle three distinct observability signals: traces, metrics, and logs. We need to decide on the ingestion mechanism for telemetry data, the storage backends for each signal type, and the visualization layer for unified observability.

Which combination of ingestion, backend, and visualization components should we deploy for the Battle Bots observability stack?

## Decision Drivers

* Vendor neutrality and avoiding lock-in
* Open-source preference to align with project philosophy
* Resource efficiency (memory, CPU, storage)
* Ease of deployment and operational overhead
* Query performance and scalability
* Long-term storage capabilities and retention policies
* Unified visualization across all observability signals
* Community support and maturity
* Cost of infrastructure and licensing
* Feature completeness for each signal type
* Flexibility to swap components independently

## Considered Options

### Ingestion
* **Option I1**: OpenTelemetry Collector
* **Option I2**: Direct export to signal-specific systems

### Traces
* **Option T1**: Tempo
* **Option T2**: Jaeger
* **Option T3**: Zipkin

### Metrics
* **Option M1**: Prometheus
* **Option M2**: Mimir

### Logs
* **Option L1**: OpenSearch
* **Option L2**: Loki

### Visualization
* **Option V1**: Grafana
* **Option V2**: Datadog
* **Option V3**: Dynatrace
* **Option V4**: New Relic

## Decision Outcome

Chosen option: **[To be determined after evaluation]**

Proposed evaluation approach:
1. Deploy all options in a test environment
2. Generate representative load for Battle Bots scenarios
3. Compare resource usage, query performance, and operational complexity
4. Select optimal combination based on Decision Drivers

### Consequences

* Good, because selecting open-source components avoids vendor lock-in
* Good, because we maintain flexibility to swap individual components
* Good, because OpenTelemetry SDK allows changing backends without re-instrumentation
* Bad, because managing multiple separate backends increases operational complexity
* Bad, because storage requirements will scale with system usage
* Bad, because commercial visualization solutions introduce ongoing costs

### Confirmation

This decision will be considered successful when:
* All three observability signals (traces, metrics, logs) are collected and queryable
* Visualization layer provides unified view across all signals
* Performance overhead of the stack is acceptable (< 10% infrastructure cost)
* Operators can effectively debug issues using the observability data
* Retention policies maintain historical data for at least 30 days
* The stack can scale to handle expected production load
* Costs (infrastructure + licensing) fit within budget constraints

## Pros and Cons of the Options

### Ingestion

#### Option I1: OpenTelemetry Collector

[OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) is a vendor-agnostic service for receiving, processing, and exporting telemetry data.

* Good, because it provides a unified ingestion point for all observability signals
* Good, because it decouples applications from backend systems
* Good, because it supports protocol translation (OTLP, Jaeger, Zipkin, Prometheus, etc.)
* Good, because it enables sending data to multiple backends simultaneously
* Good, because it provides data processing (filtering, sampling, enrichment)
* Good, because it allows backend changes without application redeployment
* Good, because it can aggregate telemetry from multiple sources
* Good, because it reduces load on applications by offloading export logic
* Good, because it supports tail-based sampling for traces
* Neutral, because it adds another component to deploy and manage
* Bad, because it introduces an additional hop in the telemetry pipeline
* Bad, because it becomes a critical point of failure if not properly configured
* Bad, because it requires additional infrastructure resources

#### Option I2: Direct Export to Signal-Specific Systems

Direct export means OpenTelemetry SDK sends telemetry directly to backend systems without an intermediary collector.

* Good, because it reduces architectural complexity (fewer moving parts)
* Good, because it eliminates an additional network hop
* Good, because it reduces infrastructure requirements (no collector to deploy)
* Good, because it has lower operational overhead
* Good, because latency is reduced for telemetry delivery
* Neutral, because applications must be configured with backend endpoints
* Bad, because changing backends requires application reconfiguration and redeployment
* Bad, because applications must handle export logic and retries
* Bad, because it couples applications to specific backend protocols
* Bad, because sending to multiple backends requires multiple exporters in each application
* Bad, because advanced processing (sampling, filtering) must be done in applications
* Bad, because it increases load on applications

### Traces

#### Option T1: Tempo

[Grafana Tempo](https://grafana.com/oss/tempo/) is a high-volume, cost-effective distributed tracing backend that requires only object storage.

* Good, because it uses object storage (S3, GCS, local disk), making it extremely cost-effective
* Good, because it has minimal operational overhead (no complex indexing)
* Good, because it integrates seamlessly with Grafana
* Good, because it scales horizontally and handles high trace volumes
* Good, because it supports multiple trace formats (Jaeger, Zipkin, OpenTelemetry)
* Good, because it has native OpenTelemetry Collector support
* Good, because it uses TraceQL for powerful trace queries
* Good, because resource usage is low compared to full-text indexing solutions
* Neutral, because trace discovery relies on trace IDs or service metadata
* Bad, because it lacks full-text search across all trace data
* Bad, because it is relatively newer compared to Jaeger/Zipkin

#### Option T2: Jaeger

[Jaeger](https://www.jaegertracing.io/) is a CNCF-graduated distributed tracing system originally developed by Uber.

* Good, because it is mature and battle-tested in production environments
* Good, because it has comprehensive trace search and filtering capabilities
* Good, because it provides its own UI in addition to Grafana integration
* Good, because it is CNCF-graduated with strong community support
* Good, because it supports multiple storage backends (Cassandra, Elasticsearch, Badger)
* Good, because it has native OpenTelemetry Collector support
* Neutral, because it requires additional storage infrastructure (database)
* Bad, because operational overhead is higher than Tempo
* Bad, because storage costs can be significant with Elasticsearch/Cassandra
* Bad, because resource consumption is higher due to indexing

#### Option T3: Zipkin

[Zipkin](https://zipkin.io/) is a distributed tracing system originally created by Twitter.

* Good, because it is very mature and widely adopted
* Good, because it has a simple architecture and deployment model
* Good, because it provides its own UI for trace visualization
* Good, because it supports multiple storage backends (MySQL, Cassandra, Elasticsearch)
* Good, because it has low resource requirements for small deployments
* Neutral, because OpenTelemetry Collector support is available but less emphasized
* Bad, because it has less active development compared to Tempo/Jaeger
* Bad, because Grafana integration is not as seamless
* Bad, because it lacks some modern features found in Tempo/Jaeger
* Bad, because community momentum has shifted toward newer solutions

### Metrics

#### Option M1: Prometheus

[Prometheus](https://prometheus.io/) is a CNCF-graduated monitoring and alerting toolkit designed for reliability and simplicity.

* Good, because it is the industry standard for metrics collection
* Good, because it has excellent Grafana integration
* Good, because it is CNCF-graduated with massive community support
* Good, because it has a powerful query language (PromQL)
* Good, because it uses efficient time-series storage
* Good, because it has extensive ecosystem of exporters and integrations
* Good, because it has native OpenTelemetry Collector support
* Good, because it is simple to deploy and operate for small-to-medium scale
* Good, because it has built-in alerting capabilities
* Neutral, because it is designed for single-server deployments
* Bad, because it has limited long-term storage capabilities (local disk only)
* Bad, because high availability requires complex federation setups
* Bad, because it doesn't scale horizontally for writes

#### Option M2: Mimir

[Grafana Mimir](https://grafana.com/oss/mimir/) is a horizontally scalable, long-term storage for Prometheus metrics.

* Good, because it provides unlimited scalability for metrics storage
* Good, because it is fully compatible with Prometheus (PromQL, remote write)
* Good, because it uses object storage for cost-effective long-term retention
* Good, because it has seamless Grafana integration
* Good, because it supports multi-tenancy out of the box
* Good, because it provides high availability natively
* Good, because it has native OpenTelemetry Collector support
* Good, because it can act as a Prometheus replacement with better scalability
* Neutral, because it is more complex to deploy than single-instance Prometheus
* Bad, because operational overhead is higher than Prometheus
* Bad, because it requires more infrastructure (object storage, multiple components)
* Bad, because it may be over-engineered for small deployments

### Logs

#### Option L1: OpenSearch

[OpenSearch](https://opensearch.org/) is an open-source fork of Elasticsearch, providing search and analytics capabilities.

* Good, because it offers powerful full-text search across all log data
* Good, because it has rich query capabilities (SQL, DSL)
* Good, because it has strong data visualization tools (OpenSearch Dashboards)
* Good, because it integrates with Grafana
* Good, because it is mature and well-documented
* Good, because it supports structured and unstructured log data
* Good, because it has native OpenTelemetry Collector support
* Neutral, because it requires significant infrastructure (cluster setup)
* Bad, because resource consumption is high (CPU, memory, storage)
* Bad, because operational complexity is significant (cluster management, shard optimization)
* Bad, because storage costs can be expensive at scale
* Bad, because it may be over-engineered if simple grep-style queries suffice

#### Option L2: Loki

[Grafana Loki](https://grafana.com/oss/loki/) is a log aggregation system designed to be cost-effective and easy to operate.

* Good, because it is extremely resource-efficient compared to full-text search solutions
* Good, because it uses object storage, making it cost-effective
* Good, because it has seamless Grafana integration (native support)
* Good, because operational overhead is minimal
* Good, because it scales horizontally
* Good, because it indexes only metadata (labels), not full log content
* Good, because it has native OpenTelemetry Collector support via OTLP
* Good, because it uses LogQL, which is similar to PromQL for consistency
* Neutral, because it requires labels to be well-structured for efficient queries
* Bad, because it doesn't support full-text search across log content
* Bad, because complex queries across unlabeled data are slow
* Bad, because it requires careful label design to avoid high cardinality

### Visualization

#### Option V1: Grafana

[Grafana](https://grafana.com/oss/grafana/) is an open-source observability platform for visualizing metrics, logs, and traces.

* Good, because it is open-source and free, avoiding vendor lock-in
* Good, because it provides unified visualization for all three signals
* Good, because it has native integrations with Prometheus, Tempo, Loki, Jaeger, and many others
* Good, because it supports custom dashboards with powerful query editors
* Good, because it has a large community and extensive plugin ecosystem
* Good, because it supports alerting capabilities
* Good, because it can correlate data across different data sources
* Good, because it aligns with open-source philosophy of the project
* Good, because it has no per-user or per-host licensing costs
* Neutral, because it requires self-hosting and operational management
* Bad, because creating effective dashboards requires learning different query languages (PromQL, LogQL, TraceQL)
* Bad, because advanced features may require additional plugins or configuration
* Bad, because it lacks some enterprise features of commercial platforms (advanced RBAC, audit logs)

#### Option V2: Datadog

[Datadog](https://www.datadoghq.com/) is a commercial observability platform providing unified monitoring, logging, and tracing.

* Good, because it provides a fully-managed, unified observability platform
* Good, because it has excellent out-of-the-box dashboards and visualizations
* Good, because it offers advanced analytics and correlation across signals
* Good, because it provides APM, RUM (Real User Monitoring), and many integrated features
* Good, because it has strong alerting, incident management, and collaboration tools
* Good, because operational overhead is minimal (fully managed)
* Good, because it has comprehensive documentation and support
* Neutral, because it can ingest OpenTelemetry data via OTLP
* Bad, because it creates significant vendor lock-in
* Bad, because it has substantial per-host/container costs that scale with usage
* Bad, because it contradicts open-source philosophy
* Bad, because costs can become prohibitive at scale
* Bad, because data is stored in Datadog's infrastructure only

#### Option V3: Dynatrace

[Dynatrace](https://www.dynatrace.com/) is an enterprise observability and AIOps platform.

* Good, because it provides AI-powered automatic root cause analysis
* Good, because it has comprehensive full-stack monitoring capabilities
* Good, because it offers automatic discovery and dependency mapping
* Good, because it excels at application performance monitoring
* Good, because operational overhead is minimal (fully managed)
* Good, because it has strong enterprise features and support
* Neutral, because it supports OpenTelemetry ingestion
* Bad, because it creates significant vendor lock-in
* Bad, because it has premium pricing model (enterprise-focused)
* Bad, because it contradicts open-source philosophy
* Bad, because it may be over-engineered for Battle Bots scale
* Bad, because costs are typically higher than other commercial solutions

#### Option V4: New Relic

[New Relic](https://newrelic.com/) is a commercial observability platform with full-stack monitoring capabilities.

* Good, because it provides unified observability for metrics, logs, and traces
* Good, because it has a generous free tier for small projects
* Good, because it offers good out-of-the-box visualizations and dashboards
* Good, because operational overhead is minimal (fully managed)
* Good, because it has native OpenTelemetry support
* Good, because it provides powerful query language (NRQL)
* Good, because it has decent community and documentation
* Neutral, because pricing is based on data ingestion rather than hosts
* Bad, because it creates vendor lock-in
* Bad, because costs scale with data volume
* Bad, because it contradicts open-source philosophy
* Bad, because advanced features require paid tiers
* Bad, because data is stored in New Relic's infrastructure only

## More Information

### Related Decisions

* ADR-0002: Observability SDK Selection - establishes OpenTelemetry as the instrumentation layer

### Related Research

* Analysis documents available in `docs/content/research_and_development/analysis/observability/`

### Architecture Overview

The observability stack follows this general architecture:

```
[Battle Bots Services] 
  → [OpenTelemetry SDKs]
    → [Ingestion: OTel Collector OR Direct Export]
      → [Traces Backend: Tempo/Jaeger/Zipkin]
      → [Metrics Backend: Prometheus/Mimir]
      → [Logs Backend: OpenSearch/Loki]
        → [Visualization: Grafana/Datadog/Dynatrace/New Relic]
```

Or with commercial visualization platforms that provide their own backends:

```
[Battle Bots Services]
  → [OpenTelemetry SDKs]
    → [Ingestion: OTel Collector OR Direct Export]
      → [Commercial Platform: Datadog/Dynatrace/New Relic]
        (handles storage + visualization)
```

### Implementation Considerations

* **Vendor Lock-in**: Commercial visualization platforms (Datadog, Dynatrace, New Relic) provide convenience but create lock-in, while open-source options (Grafana + backends) maintain flexibility
* **Cost Model**: Commercial platforms have ongoing per-host/per-GB costs; open-source has infrastructure costs but no licensing fees
* **Operational Trade-off**: Commercial platforms reduce operational overhead; open-source requires managing multiple components
* **OpenTelemetry Collector**: Strongly recommended for flexibility and decoupling, even with direct export capability
* **Starting Point**: Consider starting with lightweight open-source options (OTel Collector → Tempo + Prometheus + Loki → Grafana) and evaluating commercial platforms if needed
* **Hybrid Approach**: Possible to use OTel Collector to send data to both open-source backends and commercial platforms simultaneously
* Plan for data retention policies (e.g., 7 days high-resolution, 30 days aggregated)
* Evaluate resource usage in POC environment before production deployment

### Recommended Evaluation Criteria

For each component combination, measure:
1. **Total Cost of Ownership**: Infrastructure costs + licensing fees + operational overhead
2. **Resource Usage**: CPU, memory, disk I/O during typical load
3. **Query Performance**: P50, P95, P99 latency for common queries
4. **Storage Efficiency**: Compression ratios, storage costs over 30 days
5. **Operational Complexity**: Deployment steps, configuration complexity, upgrade process, required expertise
6. **Feature Completeness**: Support for required query patterns and use cases
7. **Vendor Lock-in Risk**: Ease of migrating to alternative solutions

### Questions to Resolve

* What is the expected volume for each signal type (traces/sec, metrics/sec, log lines/sec)?
* What are the retention requirements for each signal?
* Is long-term storage (> 30 days) required for any signal?
* What is the total budget for observability (infrastructure + licensing)?
* Do we need multi-tenancy support for multiple battle arenas?
* What are the specific query patterns needed for debugging?
* Should we prioritize operational simplicity or maintain vendor neutrality?
* Is there budget for commercial observability platforms?
* What level of operational expertise is available for managing open-source components?
* Are there specific compliance or data sovereignty requirements?
