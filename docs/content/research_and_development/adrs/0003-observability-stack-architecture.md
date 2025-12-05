---
type: docs
title: "[0003] Observability Stack Architecture"
description: >
    Deploy an OpenTelemetry Collector with Tempo for traces, Prometheus for metrics, and Loki for logs, unified by Grafana for visualization.
weight: 3
date: 2025-12-05
status: "accepted"
category: "strategic"
deciders: []
consulted: []
informed: []
---

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

Chosen option: **Open-Source Observability Stack with Grafana Ecosystem**

* **Ingestion**: OpenTelemetry Collector (Option I1)
* **Visualization**: Grafana (Option V1)
* **Logs**: Loki (Option L2)
* **Traces**: Tempo (Option T1)
* **Metrics**: Mimir (Option M2)

This combination provides a fully open-source, vendor-neutral observability stack that leverages the Grafana ecosystem for seamless integration across all three signals. The stack prioritizes cost-effectiveness (object storage for all backends), operational simplicity (unified Grafana Labs components), and maintains flexibility to swap individual components as needs evolve.

### Consequences

* Good, because all components are open-source, avoiding vendor lock-in and licensing costs
* Good, because Grafana ecosystem provides unified, seamless integration across all signals
* Good, because Tempo, Loki, and Mimir all use object storage, minimizing infrastructure costs
* Good, because OpenTelemetry Collector decouples applications from backend systems
* Good, because we maintain flexibility to swap individual components independently
* Good, because resource efficiency is optimized (Loki and Tempo avoid expensive indexing)
* Good, because the stack aligns with the project's open-source philosophy
* Good, because PromQL, LogQL, and TraceQL provide consistent query language patterns
* Neutral, because we need to self-host and manage all components
* Neutral, because Grafana requires learning three query languages (PromQL, LogQL, TraceQL)
* Bad, because operational overhead is higher than managed commercial solutions
* Bad, because Mimir is more complex to deploy than single-instance Prometheus
* Bad, because Loki requires careful label design to avoid high cardinality issues
* Bad, because we lack some enterprise features (advanced RBAC, managed services)

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

The selected observability stack architecture:

```
[Battle Bots Services] 
  → [OpenTelemetry SDKs]
    → [OpenTelemetry Collector]
      → [Tempo (Traces)]
      → [Mimir (Metrics)]
      → [Loki (Logs)]
        → [Grafana (Visualization)]
```

All three backends (Tempo, Mimir, Loki) use object storage (S3, GCS, or local disk) for cost-effective, scalable storage. Grafana provides the unified visualization layer with native support for all three data sources.

### Implementation Considerations

* **OpenTelemetry Collector Configuration**: Deploy as a sidecar or gateway to receive OTLP data and export to Tempo, Mimir, and Loki
* **Object Storage**: Configure S3-compatible object storage (or local disk for development) for all three backends
* **Label Design**: Carefully design log labels for Loki to balance query performance and cardinality
* **Mimir vs Prometheus**: Start with Mimir for long-term storage capabilities, though Prometheus could be used initially if simpler deployment is preferred
* **Data Retention**: Configure retention policies per signal type (e.g., traces: 7 days, metrics: 30 days with downsampling, logs: 14 days)
* **Deployment Order**: Deploy in order: object storage → Tempo + Mimir + Loki → OTel Collector → Grafana → configure data sources
* **Grafana Data Sources**: Configure three data sources in Grafana (Tempo, Mimir/Prometheus, Loki) for unified querying
* **Correlation**: Leverage trace IDs in logs and metrics to enable correlation across all three signals in Grafana
* **Resource Planning**: Monitor resource usage in POC environment to right-size deployments for production
* **Migration Path**: The OpenTelemetry Collector allows future migration to alternative backends without changing application instrumentation

### Rationale for Selection

This specific combination was selected based on the following factors:

1. **Unified Ecosystem**: Tempo, Loki, and Mimir are all Grafana Labs projects, ensuring consistent design patterns and seamless integration
2. **Cost Optimization**: All three backends use object storage, significantly reducing storage costs compared to indexed solutions
3. **Operational Simplicity**: Managing components from the same ecosystem reduces operational complexity versus mixing vendors
4. **Vendor Neutrality**: All components are open-source, avoiding commercial lock-in while maintaining professional quality
5. **Scalability**: Each component scales horizontally and is designed for high-volume production use
6. **Query Languages**: PromQL (Mimir), LogQL (Loki), and TraceQL (Tempo) share similar syntax patterns, reducing learning curve
7. **OpenTelemetry Native**: All components have native OpenTelemetry support via OTLP protocol
8. **Resource Efficiency**: Loki and Tempo avoid expensive full-text indexing, keeping resource usage low
9. **Project Alignment**: Open-source approach aligns with Battle Bots project philosophy
10. **Migration Flexibility**: OpenTelemetry Collector allows swapping backends without re-instrumenting applications

### Next Steps

1. **POC Deployment**: Deploy the full stack (OTel Collector + Tempo + Mimir + Loki + Grafana) in development environment
2. **Instrumentation**: Instrument Battle Bots services with OpenTelemetry SDK (per ADR-0002) and verify all three signals are collected
3. **Object Storage Setup**: Configure S3-compatible storage or local disk for backend storage
4. **Grafana Dashboards**: Create initial dashboards for key metrics, traces, and logs
5. **Label Strategy**: Define and document label naming conventions for Loki and Mimir
6. **Retention Policies**: Configure appropriate retention for each signal type based on storage constraints
7. **Performance Testing**: Generate load and measure resource usage, query performance, and storage costs
8. **Documentation**: Document deployment procedures, configuration, and operational runbooks
9. **Production Readiness**: Evaluate high availability, backup, and disaster recovery requirements before production deployment

### Open Questions

* What object storage provider should be used (S3, GCS, MinIO, local disk)?
* Should Mimir be deployed immediately, or start with Prometheus and migrate later?
* What are the specific retention requirements for each signal type?
* Do we need multi-tenancy support for isolating different battle arenas?
* What high availability requirements exist for the observability stack?
