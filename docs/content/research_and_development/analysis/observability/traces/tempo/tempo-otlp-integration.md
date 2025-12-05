---
title: "Tempo: OTLP and OpenTelemetry Collector Integration"
description: >
    Deep dive into Grafana Tempo's native OTLP support and integration with the OpenTelemetry Collector, including configuration examples, best practices, and troubleshooting guidance.
type: docs
---

## Overview

This document provides comprehensive guidance on integrating Grafana Tempo with the OpenTelemetry Collector, addressing two critical questions:

1. **Does Tempo support OTLP?** → **YES** - Native OTLP ingestion since version 1.3.0 (January 2022)
2. **Can the OTel Collector export to Tempo?** → **YES** - Full integration via `otlp` (gRPC) or `otlphttp` (HTTP) exporters

The integration enables a vendor-neutral observability pipeline where the OpenTelemetry Collector collects traces from instrumented applications and forwards them to Tempo for cost-effective, long-term trace storage, querying with TraceQL, and correlation with metrics and logs.

### Why OTLP Matters for Tracing

OpenTelemetry Protocol (OTLP) is the native protocol of the OpenTelemetry project, designed as a vendor-neutral standard for telemetry data transmission. Using OTLP with Tempo provides:

- **Native Protocol**: OTLP is Tempo's primary ingestion protocol—no translation overhead
- **Future-Proof**: OTLP is the CNCF industry standard for distributed tracing
- **Simplified Pipeline**: No protocol conversion required (App → OTel SDK → OTel Collector → OTLP → Tempo)
- **Full Fidelity**: All span attributes, events, links, and context preserved
- **Unified Stack**: Same protocol for logs (Loki), metrics (Mimir), and traces (Tempo)
- **Vendor Independence**: Easy migration between OTLP-compatible backends (Jaeger V2, Tempo, cloud vendors)

### Tempo's Position in the OTLP Ecosystem

Tempo acts as an **OTLP-native distributed tracing backend**, receiving traces via:

1. **Primary Path (Recommended)**: OpenTelemetry Collector → OTLP/gRPC (port 4317) → Tempo distributor
2. **Alternative Path**: OpenTelemetry Collector → OTLP/HTTP (port 4318) → Tempo distributor
3. **Legacy Paths (Also Supported)**: Jaeger protocol, Zipkin protocol, OpenCensus

**Recommendation**: Use **OTLP/gRPC (port 4317)** for best performance and lowest latency. Use OTLP/HTTP (port 4318) when gRPC is not available or firewall-restricted.

## OTLP Support in Tempo

### Native OTLP Support: YES

**Status**: Grafana Tempo has **full native OTLP ingestion support** for both gRPC and HTTP protocols.

**Supported Protocols**:
- ✅ **OTLP over gRPC** (port 4317): High-performance binary protocol, recommended for production
- ✅ **OTLP over HTTP** (port 4318): RESTful endpoint at `/v1/traces`, firewall-friendly

**Version History**:
- **v1.3.0 (January 2022)**: Updated OpenTelemetry libraries to v0.40.0; changed OTLP gRPC default port from legacy 55680 to standard 4317
- **v2.7 (2024)**: Updated OpenTelemetry dependencies to v0.116.0; **changed receiver binding from `0.0.0.0` to `localhost` by default** (security improvement)
- **v2.8 (2024)**: Upgraded OTLP to v1.3.0; removed deprecated InstrumentationLibrary from receivers
- **v2.9 (2024)**: Migrated internal testing from deprecated Jaeger agent/exporter to standard OTLP exporter
- **Current (2025)**: Full production-ready OTLP support with gRPC and HTTP protocols

**Minimum Recommended Version**: v1.3.0 or later (for standard OTLP port 4317)

**Production Recommendation**: v2.8 or v2.9+ (latest stable releases with OTLP 1.3.0)

### OTLP Endpoint Configuration

#### Default Endpoints and Ports

Tempo's distributor component exposes OTLP receivers on these default ports:

```
OTLP/gRPC:
  Default Port: 4317/TCP
  Default Bind: localhost:4317 (Tempo v2.7+)
  Recommended Bind: 0.0.0.0:4317 (for containerized/networked deployments)

OTLP/HTTP:
  Default Port: 4318/TCP
  HTTP Path: /v1/traces
  Default Bind: localhost:4318 (Tempo v2.7+)
  Recommended Bind: 0.0.0.0:4318 (for containerized/networked deployments)
```

**Important Change in Tempo v2.7+**: If an endpoint is **not explicitly specified**, receivers default to binding on `localhost` only (instead of `0.0.0.0`). For containerized environments (Docker, Kubernetes), you must explicitly configure `0.0.0.0` to listen on all interfaces and accept external connections.

#### Basic Tempo Configuration

**Minimal OTLP receiver configuration**:

```yaml
distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          # Defaults to localhost:4317 (Tempo v2.7+)
        http:
          # Defaults to localhost:4318 (Tempo v2.7+)
```

**Production configuration (all interfaces)**:

```yaml
distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317  # Listen on all interfaces
        http:
          endpoint: 0.0.0.0:4318  # Listen on all interfaces
```

#### TLS/mTLS Configuration

**Enable TLS on OTLP receivers**:

```yaml
distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
          tls:
            cert_file: /certs/server.crt
            key_file: /certs/server.key
        http:
          endpoint: 0.0.0.0:4318
          tls:
            cert_file: /certs/server.crt
            key_file: /certs/server.key
```

**Enable mutual TLS (mTLS)**:

```yaml
distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
          tls:
            cert_file: /certs/server.crt
            key_file: /certs/server.key
            client_ca_file: /certs/ca.crt
            require_client_auth: true
```

### Multi-Tenancy with OTLP

Tempo supports multi-tenancy using the `X-Scope-OrgID` header. When multi-tenancy is enabled, every OTLP request must include this header to identify the tenant.

**Enable multi-tenancy in Tempo**:

```yaml
server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
        http:
          endpoint: 0.0.0.0:4318

# Enable multi-tenancy
multitenancy_enabled: true
```

**OpenTelemetry Collector configuration with tenant header**:

```yaml
exporters:
  otlp:
    endpoint: tempo:4317
    headers:
      X-Scope-OrgID: "tenant-123"  # Tenant identifier
    tls:
      insecure: true
```

**Multi-tenancy benefits**:
- Isolated trace data per tenant
- Per-tenant retention policies
- Per-tenant rate limits
- Per-tenant query isolation

## OpenTelemetry Collector Export Configuration

The OpenTelemetry Collector provides two exporters for sending traces to Tempo:

1. **`otlp` exporter**: Sends traces via OTLP over gRPC (recommended)
2. **`otlphttp` exporter**: Sends traces via OTLP over HTTP

### OTLP Exporter (gRPC) - Recommended

The `otlp` exporter uses gRPC for high-performance trace transmission.

**Basic configuration**:

```yaml
exporters:
  otlp:
    endpoint: tempo:4317
    tls:
      insecure: true  # Use TLS in production
```

**Production configuration with retry and queue**:

```yaml
exporters:
  otlp:
    endpoint: tempo:4317
    tls:
      insecure: false
      cert_file: /certs/client.crt
      key_file: /certs/client.key
      ca_file: /certs/ca.crt

    # Retry configuration
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 10m

    # Sending queue (enables disk buffering)
    sending_queue:
      enabled: true
      queue_size: 1000

    # Compression (enabled by default in Tempo 2.7.1+)
    compression: snappy

    # Timeout
    timeout: 30s
```

**Key parameters**:

- `endpoint`: Tempo distributor address and port (format: `host:port`, no `http://` prefix)
- `tls.insecure`: Set to `false` for production (enable TLS)
- `retry_on_failure`: Handles transient network failures
- `sending_queue`: Buffers traces during Tempo downtime (prevents data loss)
- `compression`: `snappy` (recommended), `gzip`, or `none`
- `timeout`: Max time to wait for Tempo to acknowledge

### OTLPHTTP Exporter (HTTP) - Alternative

The `otlphttp` exporter uses HTTP/1.1 for trace transmission.

**Basic configuration**:

```yaml
exporters:
  otlphttp:
    endpoint: http://tempo:4318
    tls:
      insecure: true
```

**Production configuration**:

```yaml
exporters:
  otlphttp:
    endpoint: https://tempo:4318
    tls:
      insecure: false
      cert_file: /certs/client.crt
      key_file: /certs/client.key
      ca_file: /certs/ca.crt

    # Headers (multi-tenancy)
    headers:
      X-Scope-OrgID: "tenant-123"

    # Retry configuration
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 10m

    # Sending queue
    sending_queue:
      enabled: true
      queue_size: 1000

    # Compression
    compression: gzip

    # Timeout
    timeout: 30s
```

**Key differences from gRPC**:

- `endpoint`: Includes `http://` or `https://` prefix
- `headers`: Custom HTTP headers (useful for multi-tenancy)
- `compression`: Typically use `gzip` for HTTP

**When to use OTLPHTTP**:

- Firewall restrictions block gRPC/HTTP2
- Need to inspect traffic with HTTP debugging tools
- Existing infrastructure is HTTP/1.1-only

**When to use OTLP (gRPC)**:

- Performance is critical (gRPC is faster)
- Low latency requirements
- High trace volume (gRPC handles backpressure better)

### Batch Processor Configuration

The **batch processor** is critical for performance—it batches multiple spans before export, reducing network overhead and improving throughput.

**Recommended configuration**:

```yaml
processors:
  batch:
    send_batch_size: 1000      # Send when batch reaches 1000 spans
    send_batch_max_size: 1500  # Max batch size (hard limit)
    timeout: 10s               # Send every 10s regardless of size
```

**Parameter guidance**:

- `send_batch_size`: Grafana Labs internally uses 1,000 spans (recommended baseline)
- `send_batch_max_size`: Safety limit to prevent excessive memory usage
- `timeout`: Balance latency vs. efficiency (10s typical, 5s for low-latency needs)

**Effect of batching**:

- **Without batching**: 10,000 spans = 10,000 network requests
- **With batching (1000/batch)**: 10,000 spans = 10 network requests (100x reduction)
- **CPU/Memory savings**: Larger batches = lower overhead, but higher latency

### Memory Limiter Processor

Prevents the Collector from consuming excessive memory during traffic spikes.

**Configuration**:

```yaml
processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 512              # Hard memory limit (512MB)
    spike_limit_mib: 128        # Allow 128MB spikes above limit
```

**How it works**:

1. Collector monitors memory usage every `check_interval`
2. If memory exceeds `limit_mib`, stop accepting new data
3. Allow temporary spikes up to `limit_mib + spike_limit_mib`
4. Resume accepting data when memory drops below limit

**Recommendation**: Always use `memory_limiter` as the first processor in production.

### Complete Pipeline Configuration

**Full OpenTelemetry Collector configuration for Tempo**:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  # Memory limiter (first processor, prevents OOM)
  memory_limiter:
    check_interval: 1s
    limit_mib: 512
    spike_limit_mib: 128

  # Batch processor (critical for performance)
  batch:
    send_batch_size: 1000
    send_batch_max_size: 1500
    timeout: 10s

  # Resource detection (adds cloud/K8s metadata)
  resourcedetection:
    detectors: [env, system, docker, kubernetes]
    timeout: 5s

  # Attributes processor (add/modify span attributes)
  attributes:
    actions:
      - key: environment
        value: production
        action: upsert

exporters:
  # OTLP exporter to Tempo (gRPC)
  otlp:
    endpoint: tempo:4317
    tls:
      insecure: true
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 10m
    sending_queue:
      enabled: true
      queue_size: 1000
    compression: snappy

  # Logging exporter (debugging)
  logging:
    loglevel: info

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, resourcedetection, attributes]
      exporters: [otlp, logging]
```

**Pipeline flow**:

```
OTLP Receiver (4317/4318)
    ↓
Memory Limiter (prevent OOM)
    ↓
Batch Processor (group spans)
    ↓
Resource Detection (add metadata)
    ↓
Attributes Processor (modify spans)
    ↓
OTLP Exporter → Tempo (4317)
```

## Resource Attribute Mapping

OpenTelemetry traces include **resource attributes** (metadata about the entity producing spans, such as service name, host, container ID) and **span attributes** (operation-specific metadata).

### How Tempo Handles Attributes

**Resource Attributes**:
- Stored with every span in the trace
- Queryable via TraceQL using `resource.*` selector
- Examples: `resource.service.name`, `resource.host.name`, `resource.k8s.pod.name`

**Span Attributes**:
- Stored with individual spans
- Queryable via TraceQL using `span.*` selector
- Examples: `span.http.method`, `span.db.statement`, `span.action.type`

### Key Resource Attributes for BattleBots

**Service identification**:
```yaml
resource.service.name = "game-server"
resource.service.version = "1.2.3"
resource.service.namespace = "production"
```

**Infrastructure**:
```yaml
resource.host.name = "game-server-01"
resource.k8s.pod.name = "game-server-abc123"
resource.k8s.namespace.name = "battlebots"
resource.container.id = "docker://abc123"
```

**Environment**:
```yaml
resource.deployment.environment = "production"
resource.cloud.region = "us-east-1"
```

### TraceQL Queries with Attributes

**Find traces by service name**:
```traceql
{ resource.service.name = "game-server" }
```

**Find traces by span attribute**:
```traceql
{ span.http.status_code >= 500 }
```

**Combine resource and span attributes**:
```traceql
{ resource.service.name = "game-server"
  && span.action.type = "bot_move"
  && duration > 100ms }
```

### Best Practices for Attributes

**Use low-cardinality resource attributes**:
- ✅ Service name, environment, region (bounded values)
- ❌ User IDs, trace IDs, timestamps (unbounded values)

**Use descriptive span attributes**:
- Include operation-specific context: `http.method`, `db.statement`, `bot.id`
- Add BattleBots-specific attributes: `battle.id`, `action.type`, `player.id`

**Avoid attribute explosion**:
- Too many unique attributes can degrade Tempo performance
- Keep total unique attribute count < 1000 per service

## Sampling Strategies in the OTel Collector

Sampling controls which traces are sent to Tempo. The OpenTelemetry Collector supports multiple sampling strategies.

### Head-Based Sampling

**Decision point**: Made at span creation time (before seeing complete trace).

**Configuration**:

```yaml
processors:
  probabilistic_sampler:
    sampling_percentage: 10  # Sample 10% of traces
```

**Use cases**:
- Baseline traffic reduction
- Simple, predictable sampling rate
- Low overhead, low latency

**Limitations**:
- Cannot make decisions based on trace outcome (errors, latency)
- May miss rare but important traces

### Tail-Based Sampling (Recommended for Production)

**Decision point**: Made after seeing all/most spans in a trace.

**Configuration**:

```yaml
processors:
  tail_sampling:
    policies:
      # Always sample errors
      - name: errors-policy
        type: status_code
        status_code:
          status_codes: [ERROR]

      # Always sample slow traces
      - name: slow-requests
        type: latency
        latency:
          threshold_ms: 1000

      # Sample 10% of successful traces
      - name: sample-10-percent
        type: probabilistic
        probabilistic:
          sampling_percentage: 10

      # Sample specific span attributes
      - name: critical-actions
        type: string_attribute
        string_attribute:
          key: span.action.type
          values: [bot_death, battle_end]
```

**Architecture requirement**: All spans for a given trace ID must be routed to the same Collector instance.

**Solution**: Use a two-tier architecture:

```
Application Instances
    ↓
Agent Collectors (no tail sampling)
    ↓
Load Balancing Exporter (shard by trace ID)
    ↓
Gateway Collectors (tail sampling)
    ↓
Tempo
```

**Benefits**:
- Capture 100% of errors
- Capture 100% of slow requests
- Sample normal traffic to reduce volume

### Sampling Recommendations for BattleBots

**Development**: 100% sampling (no sampling)

**Production**:
```yaml
processors:
  tail_sampling:
    policies:
      # Always sample errors
      - name: errors
        type: status_code
        status_code:
          status_codes: [ERROR]

      # Always sample critical game events
      - name: critical-events
        type: string_attribute
        string_attribute:
          key: span.event.type
          values: [bot_death, battle_end, matchmaking_failed]

      # Always sample slow battles (> 10 seconds)
      - name: slow-battles
        type: latency
        latency:
          threshold_ms: 10000

      # Sample 20% of normal bot actions
      - name: baseline
        type: probabilistic
        probabilistic:
          sampling_percentage: 20
```

## Complete Working Examples

### Example 1: Basic OpenTelemetry Collector → Tempo

**OpenTelemetry Collector config (`otel-collector.yaml`)**:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    send_batch_size: 1000
    timeout: 10s

exporters:
  otlp:
    endpoint: tempo:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
```

**Tempo config (`tempo.yaml`)**:

```yaml
server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
        http:
          endpoint: 0.0.0.0:4318

storage:
  trace:
    backend: local
    local:
      path: /var/tempo
```

**Docker Compose**:

```yaml
version: '3'

services:
  otel-collector:
    image: otel/opentelemetry-collector:latest
    command: ["--config=/etc/otel-collector.yaml"]
    volumes:
      - ./otel-collector.yaml:/etc/otel-collector.yaml
    ports:
      - "4317:4317"
      - "4318:4318"

  tempo:
    image: grafana/tempo:latest
    command: ["-config.file=/etc/tempo.yaml"]
    volumes:
      - ./tempo.yaml:/etc/tempo.yaml
      - tempo-data:/var/tempo
    ports:
      - "3200:3200"

volumes:
  tempo-data:
```

### Example 2: Production with Multi-Tenancy and TLS

**OpenTelemetry Collector config**:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

  batch:
    send_batch_size: 1000
    timeout: 10s

  resourcedetection:
    detectors: [env, system, kubernetes]

exporters:
  otlp:
    endpoint: tempo:4317
    headers:
      X-Scope-OrgID: "tenant-battlebots"
    tls:
      insecure: false
      cert_file: /certs/client.crt
      key_file: /certs/client.key
      ca_file: /certs/ca.crt
    retry_on_failure:
      enabled: true
    sending_queue:
      enabled: true
      queue_size: 1000
    compression: snappy

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, resourcedetection]
      exporters: [otlp]
```

**Tempo config**:

```yaml
server:
  http_listen_port: 3200

distributor:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
          tls:
            cert_file: /certs/server.crt
            key_file: /certs/server.key

ingester:
  lifecycler:
    ring:
      replication_factor: 3

storage:
  trace:
    backend: s3
    s3:
      bucket: tempo-traces
      endpoint: s3.amazonaws.com
      region: us-east-1

multitenancy_enabled: true

overrides:
  per_tenant_override_config: /etc/overrides.yaml
```

## Troubleshooting

### Issue: Connection Refused on Port 4317/4318

**Symptoms**:
```
Error: rpc error: code = Unavailable desc = connection error: dial tcp 192.168.1.10:4317: connect: connection refused
```

**Causes**:
1. Tempo's OTLP receivers not configured properly
2. Tempo v2.7+ defaulting to `localhost` (not accessible from other containers)
3. Firewall blocking ports 4317/4318

**Solutions**:

1. Verify Tempo OTLP configuration:
   ```yaml
   distributor:
     receivers:
       otlp:
         protocols:
           grpc:
             endpoint: 0.0.0.0:4317  # Must be 0.0.0.0, not localhost
   ```

2. Check Tempo logs:
   ```bash
   docker logs tempo | grep "OTLP"
   ```

3. Verify ports are exposed:
   ```bash
   docker ps | grep tempo
   # Should show 4317/tcp and 4318/tcp
   ```

4. Test connectivity:
   ```bash
   # From OTel Collector container
   nc -zv tempo 4317
   ```

### Issue: Traces Not Appearing in Tempo

**Symptoms**:
- OTel Collector shows no errors
- Traces not visible in Grafana

**Debugging steps**:

1. **Enable debug logging in OTel Collector**:
   ```yaml
   exporters:
     logging:
       loglevel: debug

   service:
     pipelines:
       traces:
         receivers: [otlp]
         processors: [batch]
         exporters: [otlp, logging]  # Add logging
   ```

2. **Check OTel Collector metrics**:
   ```bash
   curl http://otel-collector:8888/metrics | grep otelcol_exporter_sent_spans
   # Should show > 0 if spans are being sent
   ```

3. **Check Tempo metrics**:
   ```bash
   curl http://tempo:3200/metrics | grep tempo_distributor_spans_received_total
   # Should show > 0 if Tempo is receiving spans
   ```

4. **Wait for ingester flush**:
   - Traces are buffered in ingesters before flushing to storage
   - Default `max_block_duration` = 30-60 minutes
   - For POC, set `max_block_duration: 5m` for faster availability

5. **Query by trace ID directly**:
   ```bash
   curl "http://tempo:3200/api/traces/<trace-id>"
   ```

### Issue: High Memory Usage in OTel Collector

**Symptoms**:
- OTel Collector OOM (out of memory)
- Collector restarts frequently

**Solutions**:

1. **Add memory limiter processor**:
   ```yaml
   processors:
     memory_limiter:
       check_interval: 1s
       limit_mib: 512

   service:
     pipelines:
       traces:
         processors: [memory_limiter, batch]  # memory_limiter first
   ```

2. **Reduce batch size**:
   ```yaml
   processors:
     batch:
       send_batch_size: 500    # Reduce from 1000
       timeout: 5s             # Reduce from 10s
   ```

3. **Increase OTel Collector resources**:
   ```yaml
   # Docker Compose
   otel-collector:
     deploy:
       resources:
         limits:
           memory: 1G
   ```

### Issue: Spans Dropped or Missing

**Symptoms**:
- Incomplete traces (missing spans)
- OTel Collector shows dropped spans

**Causes**:
1. Rate limiting in Tempo
2. Batch processor queue full
3. Network timeouts

**Solutions**:

1. **Check OTel Collector queue metrics**:
   ```bash
   curl http://otel-collector:8888/metrics | grep otelcol_exporter_queue_size
   ```

2. **Increase queue size**:
   ```yaml
   exporters:
     otlp:
       sending_queue:
         enabled: true
         queue_size: 5000  # Increase from 1000
   ```

3. **Check Tempo rate limits**:
   ```yaml
   # tempo.yaml
   overrides:
     defaults:
       ingestion_rate_limit_bytes: 10485760  # 10 MB/s
       ingestion_burst_size_bytes: 20971520  # 20 MB burst
   ```

4. **Enable retry on failure**:
   ```yaml
   exporters:
     otlp:
       retry_on_failure:
         enabled: true
         max_elapsed_time: 10m
   ```

### Issue: "Unimplemented" or Wrong Endpoint Errors

**Symptoms**:
```
Error: rpc error: code = Unimplemented desc = unknown service
```

**Cause**: OpenTelemetry Collector trying to send to wrong endpoint or using wrong protocol.

**Solutions**:

1. **Verify endpoint format**:
   - ✅ Correct: `endpoint: tempo:4317` (no protocol prefix for gRPC)
   - ❌ Wrong: `endpoint: http://tempo:4317` (http prefix for gRPC)

2. **Match exporter to protocol**:
   - Use `otlp` exporter for gRPC (port 4317)
   - Use `otlphttp` exporter for HTTP (port 4318 with `http://` prefix)

3. **Verify Tempo receiver is enabled**:
   ```yaml
   distributor:
     receivers:
       otlp:
         protocols:
           grpc:  # Must be enabled
           http:  # Must be enabled
   ```

## BattleBots Integration Patterns

### Complete LGTM Stack Configuration

**OpenTelemetry Collector configuration for BattleBots**:

```yaml
receivers:
  # OTLP receiver for application traces
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

  batch:
    send_batch_size: 1000
    timeout: 10s

  # Add BattleBots-specific attributes
  attributes:
    actions:
      - key: platform
        value: battlebots
        action: upsert
      - key: environment
        value: production
        action: upsert

  # Detect cloud/K8s metadata
  resourcedetection:
    detectors: [env, system, docker, kubernetes]

  # Tail sampling for intelligent trace selection
  tail_sampling:
    policies:
      # Always sample errors
      - name: errors
        type: status_code
        status_code:
          status_codes: [ERROR]

      # Always sample critical game events
      - name: critical-events
        type: string_attribute
        string_attribute:
          key: span.event.type
          values: [bot_death, battle_end, matchmaking_failed]

      # Always sample slow operations
      - name: slow-operations
        type: latency
        latency:
          threshold_ms: 1000

      # Sample 20% of normal traffic
      - name: baseline
        type: probabilistic
        probabilistic:
          sampling_percentage: 20

exporters:
  # Tempo - traces
  otlp:
    endpoint: tempo:4317
    tls:
      insecure: true
    retry_on_failure:
      enabled: true
    sending_queue:
      enabled: true
    compression: snappy

  # Logging for debugging
  logging:
    loglevel: info

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, attributes, resourcedetection, tail_sampling]
      exporters: [otlp, logging]
```

### BattleBots-Specific TraceQL Queries

**Find all battles with errors**:
```traceql
{ resource.service.name = "game-server" && status = error }
```

**Find slow bot actions (> 100ms)**:
```traceql
{ span.action.type =~ "bot_.*" && duration > 100ms }
```

**Find battles where a specific bot died**:
```traceql
{ span.event.type = "bot_death" && span.bot.id = "bot_xyz123" }
```

**Find database queries in battle processing**:
```traceql
{ resource.service.name = "game-server"
  && span.db.statement != nil
  && duration > 50ms }
```

**Aggregate: Count battles by winner**:
```traceql
{ span.battle.result != nil }
| by(span.battle.winner)
| count()
```

**Aggregate: Average battle duration by game mode**:
```traceql
{ resource.service.name = "game-server"
  && span.name = "ExecuteBattle" }
| by(span.battle.mode)
| avg(duration)
```

### Trace-Metric-Log Correlation Example

**Application instrumentation (Go)**:

```go
import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    "github.com/sirupsen/logrus"
)

func ProcessBotAction(ctx context.Context, action BotAction) error {
    // Create span
    tracer := otel.Tracer("game-server")
    ctx, span := tracer.Start(ctx, "ProcessBotAction")
    defer span.End()

    // Add span attributes
    span.SetAttributes(
        attribute.String("action.type", action.Type),
        attribute.String("bot.id", action.BotID),
        attribute.String("battle.id", action.BattleID),
    )

    // Log with trace context
    traceID := span.SpanContext().TraceID().String()
    spanID := span.SpanContext().SpanID().String()

    logger.WithFields(logrus.Fields{
        "trace_id": traceID,
        "span_id":  spanID,
        "bot_id":   action.BotID,
    }).Info("Processing bot action")

    // ... process action ...

    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "Action failed")
        logger.WithField("trace_id", traceID).Error("Action failed", err)
        return err
    }

    return nil
}
```

**Workflow in Grafana**:

1. **Start with metric alert**: "Bot action latency P99 > 500ms"
2. **Click exemplar**: Jump to example slow trace
3. **Identify span**: See `ProcessBotAction` span took 800ms
4. **View logs**: Click "Logs for this span" → see detailed error logs with same `trace_id`
5. **Root cause**: Logs show "database connection pool exhausted at 12:34:56"

## Further Reading

### Official Documentation

- [Configure Tempo](https://grafana.com/docs/tempo/latest/configuration/)
- [Pushing Spans with HTTP](https://grafana.com/docs/tempo/latest/api_docs/pushing-spans-with-http/)
- [Enable Multi-tenancy](https://grafana.com/docs/tempo/latest/operations/manage-advanced-systems/multitenancy/)
- [Tempo Release Notes](https://grafana.com/docs/tempo/latest/release-notes/)

### OpenTelemetry Collector Resources

- [OpenTelemetry Collector Configuration](https://opentelemetry.io/docs/collector/configuration/)
- [OTLP Exporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlpexporter)
- [OTLPHTTP Exporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter)
- [Batch Processor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor)
- [Tail Sampling Processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/tailsamplingprocessor)

### Integration Guides

- [How to Send Traces to Grafana Cloud Tempo with OpenTelemetry Collector](https://grafana.com/blog/2021/04/13/how-to-send-traces-to-grafana-clouds-tempo-service-with-opentelemetry-collector/)
- [End-to-End Distributed Tracing in Kubernetes with Grafana Tempo and OpenTelemetry](https://www.civo.com/learn/distributed-tracing-kubernetes-grafana-tempo-opentelemetry)
- [Send Data to the Grafana Cloud OTLP Endpoint](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)

### Troubleshooting Resources

- [Tempo Troubleshooting Guide](https://grafana.com/docs/tempo/latest/troubleshooting/)
- [OpenTelemetry Collector Troubleshooting](https://opentelemetry.io/docs/collector/troubleshooting/)
- [Grafana Community Forums - Tempo](https://community.grafana.com/c/grafana-tempo/)

### Related BattleBots Documentation

- [Tempo Overview](tempo-overview.md) - Architecture, deployment, and how to run Tempo
- [OpenTelemetry Collector: Traces Support](../otel-collector/otel-collector-traces.md) - Understanding distributed tracing
- [Grafana Loki: OTLP Integration](../logs/loki/loki-otlp-integration.md) - Log correlation
- [Grafana Mimir: OTLP Integration](../metrics/mimir/mimir-otlp-integration.md) - Metrics correlation
