---
type: docs
title: "[0002] Observability SDK Selection"
description: >
    Select OpenTelemetry as the observability SDK for instrumenting metrics, traces, and logs across the Battle Bots platform.
weight: 2
date: 2025-12-05
status: "accepted"
category: "strategic"
deciders: []
consulted: []
informed: []
---

## Context and Problem Statement

Battle Bots requires comprehensive observability to monitor game server performance, track bot behavior, debug issues, and visualize battle state in real-time. We need to select an observability SDK that will provide instrumentation for metrics, traces, and logs across our system. The choice will impact vendor lock-in, cost, flexibility, and integration complexity.

Which observability SDK should we adopt for instrumenting the Battle Bots platform?

## Decision Drivers

* Vendor neutrality and ability to switch backends
* Support for metrics, traces, and logs (unified observability)
* Language support (especially for our game server implementation)
* Integration with container environments
* Community adoption and long-term sustainability
* Cost implications (SDK licensing, vendor fees)
* Ease of integration and developer experience
* Performance overhead
* Support for custom attributes and semantic conventions

## Considered Options

* OpenTelemetry
* Datadog SDK
* Sentry SDK

## Decision Outcome

Chosen option: **OpenTelemetry**, because it is an open standard that prevents vendor lock-in and aligns with the open-source nature of the Battle Bots project. As a CNCF-graduated project, it provides long-term sustainability and broad industry adoption while maintaining flexibility to switch observability backends without re-instrumentation.

### Consequences

* Good, because we maintain complete vendor neutrality and can switch backends (Prometheus, Jaeger, Loki, etc.) without changing instrumentation code
* Good, because we adopt an industry-standard approach to observability that is widely supported and documented
* Good, because the open-source SDK has no licensing costs and aligns with project philosophy
* Good, because semantic conventions will ensure consistent telemetry across all Battle Bots components
* Good, because comprehensive language support enables bot developers to use OpenTelemetry in their preferred languages
* Bad, because we need to deploy and manage separate backend infrastructure for metrics, traces, and logs
* Bad, because initial setup requires more configuration compared to all-in-one commercial solutions
* Bad, because debugging telemetry pipelines may require deeper understanding of the OpenTelemetry architecture

### Confirmation

This decision will be considered successful when:
* Game server successfully exports metrics, traces, and logs via OpenTelemetry SDK
* Observability backends (selected via future ADRs) receive and display telemetry data correctly
* Developer experience for adding custom instrumentation is straightforward
* Performance overhead of instrumentation is acceptable (< 5% CPU/memory impact)
* We can successfully switch between different backend providers without code changes

## Pros and Cons of the Options

### OpenTelemetry

[OpenTelemetry](https://opentelemetry.io/) is an open-source observability framework providing vendor-neutral APIs, SDKs, and tools for generating and collecting telemetry data (metrics, logs, and traces).

* Good, because it is vendor-neutral and prevents lock-in to any specific observability backend
* Good, because it supports all three pillars of observability (metrics, traces, logs)
* Good, because it has broad language support including Go, Java, Python, JavaScript, and many others
* Good, because it is CNCF-graduated with strong industry adoption and long-term sustainability
* Good, because it allows flexibility to switch backends (Prometheus, Jaeger, Loki, commercial vendors) without changing instrumentation
* Good, because it has native Kubernetes and container environment support
* Good, because it defines semantic conventions for consistent telemetry across services
* Good, because it is free and open-source with no licensing costs
* Neutral, because it requires separate backend infrastructure (exporters to Prometheus, Jaeger, etc.)
* Bad, because initial setup complexity is higher than vendor-specific SDKs
* Bad, because some advanced features may require vendor-specific extensions

### Datadog SDK

[Datadog](https://www.datadoghq.com/) provides proprietary SDKs for instrumenting applications and sending telemetry to the Datadog platform.

* Good, because it offers a unified, fully-managed observability platform
* Good, because it provides excellent out-of-the-box dashboards and visualizations
* Good, because it has automatic instrumentation for many frameworks and libraries
* Good, because it offers strong APM (Application Performance Monitoring) features
* Good, because it has excellent documentation and developer experience
* Good, because it includes alerting, incident management, and collaboration features
* Neutral, because it requires a Datadog account and subscription
* Bad, because it creates vendor lock-in - switching away requires significant re-instrumentation
* Bad, because it has ongoing per-host/container costs that scale with usage
* Bad, because telemetry data is only compatible with Datadog's backend
* Bad, because it may not align with open-source philosophy of the project

### Sentry SDK

[Sentry](https://sentry.io/) is primarily an error tracking and performance monitoring platform with SDKs for multiple languages.

* Good, because it excels at error tracking and exception reporting
* Good, because it has good performance monitoring capabilities
* Good, because it offers distributed tracing
* Good, because it has strong developer experience and debugging tools
* Good, because it provides issue grouping and notification features
* Good, because it has a generous free tier for open-source projects
* Neutral, because it focuses more on errors/exceptions than comprehensive observability
* Bad, because metrics support is less mature compared to dedicated observability platforms
* Bad, because it creates vendor lock-in similar to Datadog
* Bad, because log aggregation is not a primary feature
* Bad, because it's more specialized for error tracking than full observability

## More Information

### Related Research

* Analysis documents on observability backends are available in `docs/content/research_and_development/analysis/observability/`
* The POC user journey (ADR-0000 reference pending) identifies observability as a critical requirement

### Implementation Considerations

* If OpenTelemetry is chosen, separate ADRs will be needed for:
  * Metrics backend selection (e.g., Prometheus, Mimir)
  * Tracing backend selection (e.g., Jaeger, Tempo)
  * Log aggregation backend selection (e.g., Loki)
  * OpenTelemetry Collector deployment strategy
* If Datadog or Sentry is chosen, consider OpenTelemetry compatibility for future flexibility
* Consider hybrid approaches (e.g., OpenTelemetry for instrumentation + commercial backend)

### Questions to Resolve

* What is the expected scale of telemetry data (events/sec, log volume)?
* What is the budget for observability infrastructure and services?
* Is vendor neutrality a hard requirement or nice-to-have?
* What are the specific visualization and alerting requirements?
* Will bots require observability instrumentation, or only game server?
