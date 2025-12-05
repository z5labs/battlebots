---
title: "Loki: OTLP Integration"
description: >
    Detailed guide for integrating Grafana Loki with OpenTelemetry Collector via native OTLP support.
type: docs
---

## Overview

Grafana Loki introduced native OpenTelemetry Protocol (OTLP) support in version 3.0, marking a significant advancement in how logs can be ingested into Loki. This native integration allows applications instrumented with OpenTelemetry to send logs directly to Loki using the standardized OTLP format, eliminating the need for format transformations and simplifying the observability pipeline.

The native OTLP endpoint provides a fully OpenTelemetry-compliant ingestion path where logs sent in OTLP format are stored directly in Loki without requiring conversion to JSON or logfmt blobs. This approach leverages Loki's structured metadata feature, which stores log attributes and other OpenTelemetry LogRecord fields separately from the log body itself. The result is a more intuitive query experience and better performance, as queries no longer need to parse JSON at runtime to access fields.

For the BattleBots observability stack, native OTLP integration offers several advantages: unified telemetry collection across logs, metrics, and traces through the OpenTelemetry Collector; simplified configuration compared to legacy exporters; better correlation between logs and traces through preserved TraceId and SpanId fields; and vendor portability, making it easier to migrate between observability backends without changing instrumentation.

## OTLP Support in Loki

**Answer: YES** - Loki versions 3.0 and later natively support the OpenTelemetry Protocol (OTLP) for log ingestion.

### Native OTLP Endpoint

Loki exposes an OTLP-compliant endpoint at `/otlp` that accepts OpenTelemetry log data. When clients send logs to this endpoint, the collector automatically appends the appropriate path suffix (`/v1/logs`), resulting in requests to `/otlp/v1/logs`.

**Supported Protocols:**

- **HTTP**: POST requests using HTTP/1.1 or HTTP/2
- **gRPC**: Unary RPC calls using the OTLP service definition

**Default Port:**

- Loki typically runs on port `3100` for all HTTP endpoints, including OTLP

### Version Requirements

**Minimum Loki Version:** 3.0 or later

**Schema Requirements:**
- Schema version: `v13` or higher (required for structured metadata)
- Index type: `tsdb` (Time Series Database index required)

Structured metadata is essential for OTLP ingestion because it stores the OpenTelemetry LogRecord fields (resource attributes, instrumentation scope, log attributes) separately from the log body. Without schema v13 and tsdb, Loki cannot properly handle OTLP data.

### Benefits of Native OTLP vs Legacy Loki Exporter

The legacy `lokiexporter` component in the OpenTelemetry Collector encoded logs as JSON or logfmt blobs with Loki-specific label conventions. The native OTLP endpoint provides several improvements:

1. **Simplified Querying**: No JSON parsing required at query time. Instead of `{job="dev/auth"} | json | severity="INFO"`, you can query directly: `{service_name="auth"} | severity_text="INFO"`

2. **Cleaner Log Bodies**: The log message is stored as-is rather than wrapped in a JSON structure. A log "user logged in" is stored exactly as that string, with metadata in structured fields.

3. **Standard Resource Labels**: Uses OpenTelemetry semantic conventions (`service_name`, `service_namespace`) instead of custom labels (`job=service.namespace/service.name`)

4. **Better Performance**: Structured metadata allows efficient filtering without parsing the entire log body

5. **Vendor Portability**: Standard OTLP configuration works across multiple backends without Loki-specific hints

6. **Future-Proof**: The native endpoint represents Grafana's strategic direction for log ingestion

## OTLP Endpoint Configuration

To enable OTLP ingestion in Loki, you must configure structured metadata support and optionally customize which resource attributes become index labels.

### Enabling Structured Metadata

Structured metadata is enabled by default in Loki 3.0+, but you should explicitly configure it in your limits:

```yaml
limits_config:
  allow_structured_metadata: true
  max_structured_metadata_entries_count: 128  # Maximum metadata entries per log record
```

### Schema Configuration

Your Loki schema must use version 13 or higher with the tsdb index:

```yaml
schema_config:
  configs:
    - from: "2024-04-01"
      store: tsdb
      object_store: s3  # or filesystem, gcs, azure, etc.
      schema: v13
      index:
        prefix: loki_index_
        period: 24h
```

### Complete Loki Configuration Example

Here's a complete Loki configuration with OTLP support enabled:

```yaml
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096
  log_level: info

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: "2024-04-01"
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: loki_index_
        period: 24h

limits_config:
  allow_structured_metadata: true
  max_structured_metadata_entries_count: 128
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  ingestion_rate_mb: 10
  ingestion_burst_size_mb: 20
  per_stream_rate_limit: 5MB
  per_stream_rate_limit_burst: 15MB

  # OTLP-specific configuration
  otlp_config:
    resource_attributes:
      # Configure which resource attributes become index labels
      attributes_config:
        - action: index_label
          attributes:
            - service.name
            - service.namespace
            - deployment.environment
            - k8s.cluster.name
            - k8s.namespace.name
            - cloud.region
            - cloud.provider
        # Convert high-cardinality attributes to structured metadata
        - action: structured_metadata
          attributes:
            - k8s.pod.name
            - service.instance.id
            - process.pid

storage_config:
  tsdb_shipper:
    active_index_directory: /loki/tsdb-index
    cache_location: /loki/tsdb-cache
  filesystem:
    directory: /loki/chunks

compactor:
  working_directory: /loki/compactor
  compaction_interval: 10m
  retention_enabled: true
  retention_delete_delay: 2h
  retention_delete_worker_count: 150

querier:
  max_concurrent: 4

query_scheduler:
  max_outstanding_requests_per_tenant: 4096

frontend:
  max_outstanding_per_tenant: 4096
```

## OTel Collector Export Configuration

**Answer: YES** - The OpenTelemetry Collector can export logs to Loki using the `otlphttp` exporter.

### Basic OTLP HTTP Exporter

The recommended exporter for Loki is `otlphttp/logs`:

```yaml
exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"
    tls:
      insecure: true  # For local development without TLS
```

**Important:** Do not append `/v1/logs` to the endpoint URL. The OTLP exporter automatically adds the appropriate path suffix.

### Complete Exporter Configuration

Here's a production-ready configuration with retry, timeout, and queue settings:

```yaml
exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"

    # TLS Configuration
    tls:
      insecure: false
      ca_file: /etc/ssl/certs/loki-ca.crt
      cert_file: /etc/ssl/certs/client.crt
      key_file: /etc/ssl/private/client.key
      min_version: "1.2"

    # Timeout for individual requests (default: 30s recommended)
    timeout: 30s

    # Retry configuration
    retry_on_failure:
      enabled: true
      initial_interval: 5s    # Time to wait before first retry
      max_interval: 30s       # Maximum backoff interval
      max_elapsed_time: 300s  # Give up after 5 minutes

    # Queue configuration for reliability
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 5000
      storage: file_storage  # Reference to file_storage extension for persistence

    # Compression (gzip recommended for production)
    compression: gzip

    # Headers (for authentication, multi-tenancy)
    headers:
      X-Scope-OrgID: "battlebots"
```

### File Storage Extension for Persistence

To persist queued logs across collector restarts, configure the file storage extension:

```yaml
extensions:
  file_storage:
    directory: /var/lib/otelcol/file_storage
    timeout: 10s

service:
  extensions: [file_storage]
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp/logs]
```

### Batch Processor Configuration

Always include a batch processor before the exporter to optimize throughput:

```yaml
processors:
  batch:
    timeout: 10s           # Send batch after this duration
    send_batch_size: 8192  # Send when batch reaches this size
    send_batch_max_size: 16384  # Never exceed this size
```

### Complete Service Pipeline

```yaml
service:
  extensions: [file_storage]
  pipelines:
    logs:
      receivers: [otlp, filelog]
      processors: [resource_detection, batch, attributes]
      exporters: [otlphttp/logs]
```

## Resource Attribute Mapping

When logs arrive via OTLP, resource attributes from the OpenTelemetry SDK map to either index labels or structured metadata in Loki.

### Default Resource Attributes as Index Labels

By default, Loki converts these 17 resource attributes to index labels:

1. `service.name` → `service_name`
2. `service.namespace` → `service_namespace`
3. `service.instance.id` → `service_instance_id`
4. `deployment.environment` → `deployment_environment`
5. `cloud.region` → `cloud_region`
6. `cloud.availability_zone` → `cloud_availability_zone`
7. `cloud.platform` → `cloud_platform`
8. `k8s.cluster.name` → `k8s_cluster_name`
9. `k8s.namespace.name` → `k8s_namespace_name`
10. `k8s.pod.name` → `k8s_pod_name`
11. `k8s.container.name` → `k8s_container_name`
12. `container.name` → `container_name`
13. `k8s.replicaset.name` → `k8s_replicaset_name`
14. `k8s.deployment.name` → `k8s_deployment_name`
15. `k8s.statefulset.name` → `k8s_statefulset_name`
16. `k8s.daemonset.name` → `k8s_daemonset_name`
17. `k8s.cronjob.name` → `k8s_cronjob_name`

**Note:** Attribute names with dots (`.`) are converted to underscores (`_`) for Loki label compatibility.

### Attribute Transformation in Collector

To add or modify resource attributes before sending to Loki:

```yaml
processors:
  resource:
    attributes:
      # Add static attributes
      - key: deployment.environment
        value: production
        action: insert

      # Copy attributes with character transformations
      - key: service_name
        from_attribute: service.name
        action: insert

      # Delete attributes you don't want
      - key: telemetry.sdk.version
        action: delete

  attributes:
    actions:
      # Add log-level attributes
      - key: environment
        value: production
        action: insert

      # Extract correlation IDs from log body
      - key: correlation_id
        pattern: "correlation_id=([a-z0-9-]+)"
        action: extract
```

### Example Attribute to Label Mapping

**Input (OpenTelemetry SDK):**
```json
{
  "resource": {
    "attributes": {
      "service.name": "game-server",
      "service.namespace": "battlebots",
      "deployment.environment": "production",
      "k8s.pod.name": "game-server-5d7c8f9b-xq2wr",
      "k8s.namespace.name": "battlebots-prod"
    }
  }
}
```

**Output (Loki Labels):**
```
{
  service_name="game-server",
  service_namespace="battlebots",
  deployment_environment="production",
  k8s_namespace_name="battlebots-prod"
}
```

Note: `k8s.pod.name` should be converted to structured metadata (see Label Strategy section).

## Label Strategy and Best Practices

Loki's performance depends heavily on proper label cardinality management. Every unique combination of label values creates a new stream, and too many streams degrade performance significantly.

### Avoiding High Cardinality Issues

**Problem:** High cardinality causes Loki to build a huge index and flush thousands of tiny chunks to object storage, resulting in poor performance and high costs.

**High-Cardinality Attributes (Avoid as Labels):**
- `k8s.pod.name` - Each pod instance creates a new label value
- `service.instance.id` - Each service instance is unique
- `process.pid` - Changes on every process restart
- User IDs, request IDs, transaction IDs
- Timestamps, UUIDs, hashes

### Which Attributes to Use as Labels

**Good Label Candidates (Low Cardinality):**
- `service.name` - Limited number of services
- `service.namespace` - Few namespaces (dev, staging, prod)
- `deployment.environment` - Usually 3-5 values
- `k8s.cluster.name` - Fixed cluster names
- `k8s.namespace.name` - Limited namespaces per cluster
- `cloud.region` - Fixed set of regions
- `log.severity` or `severity_text` - Limited severity levels

**Cardinality Rule of Thumb:** Keep total stream count under 10,000. With 5 labels averaging 10 values each, you get 10^5 = 100,000 streams (too many). Reduce to 3-4 labels with controlled values.

### Which Attributes to Keep as Structured Metadata

Configure high-cardinality attributes as structured metadata:

```yaml
limits_config:
  otlp_config:
    resource_attributes:
      attributes_config:
        - action: structured_metadata
          attributes:
            - k8s.pod.name
            - service.instance.id
            - process.pid
            - process.command_line
            - host.id
```

Structured metadata remains queryable but doesn't create new streams:

```logql
{service_name="game-server"} | k8s_pod_name="game-server-5d7c8f9b-xq2wr"
```

### Recommended Label Strategy for BattleBots

```yaml
limits_config:
  otlp_config:
    resource_attributes:
      attributes_config:
        # Index labels (low cardinality)
        - action: index_label
          attributes:
            - service.name
            - service.namespace
            - deployment.environment
            - k8s.namespace.name
            - cloud.region

        # Structured metadata (high cardinality or optional)
        - action: structured_metadata
          attributes:
            - k8s.pod.name
            - k8s.container.name
            - service.instance.id
            - host.name
            - process.pid
```

**Expected Cardinality:**
- `service.name`: ~10 services (game-server, matchmaker, auth, etc.)
- `service.namespace`: 1 value (battlebots)
- `deployment.environment`: 3 values (dev, staging, production)
- `k8s.namespace.name`: ~5 namespaces
- `cloud.region`: ~3 regions

**Total streams:** 10 × 1 × 3 × 5 × 3 = 450 streams (excellent)

## Log-Trace Correlation

OpenTelemetry's unified data model enables seamless correlation between logs and traces through TraceId and SpanId fields embedded in log records.

### TraceID and SpanID Storage

When applications instrumented with OpenTelemetry SDKs emit logs within a trace context, the SDK automatically includes:

- **TraceId**: Unique identifier for the entire trace
- **SpanId**: Unique identifier for the current span
- **TraceFlags**: Sampling and other flags

Loki stores these fields as structured metadata, making them queryable without parsing the log body.

### Querying Logs by Trace ID

**Find all logs for a specific trace:**
```logql
{service_name="game-server"} | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```

**Find logs with any trace context:**
```logql
{service_name="game-server"} | trace_id != ""
```

**Find logs for a specific span:**
```logql
{service_name="game-server"} | span_id="00f067aa0ba902b7"
```

### Grafana Integration for Correlation

Configure Grafana data source correlations to link Loki and Tempo:

**Loki Data Source Configuration:**
```json
{
  "derivedFields": [
    {
      "datasourceUid": "tempo-uid",
      "matcherRegex": "trace_id=(\\w+)",
      "name": "TraceID",
      "url": "${__value.raw}"
    }
  ]
}
```

This creates clickable trace ID links in the Grafana Explore view, allowing you to:

1. View a log entry in Loki
2. Click the trace ID link
3. Jump directly to the full trace in Tempo
4. See the complete request flow with timing information

### Example Correlation Workflow

**Scenario:** Investigating a slow game server response

1. **Start with metrics:** Notice elevated response times in Prometheus metrics
2. **Query slow traces:** Find traces with duration > 1s in Tempo
3. **Jump to logs:** Click trace ID to see all logs for that request
4. **Identify root cause:** Read detailed error messages and debug logs
5. **Correlate with resources:** Use `k8s_pod_name` metadata to check pod health

**Query pattern:**
```logql
{service_name="game-server"}
| trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
| severity_text="ERROR"
```

## Authentication & Multi-tenancy

Loki supports both basic authentication and multi-tenant deployments for OTLP ingestion.

### Basic Authentication Setup

**Collector Configuration with Basic Auth:**

```yaml
extensions:
  basicauth/otlp:
    client_auth:
      username: "battlebots-collector"
      password: "${LOKI_PASSWORD}"  # Use environment variable

exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"
    auth:
      authenticator: basicauth/otlp

service:
  extensions: [basicauth/otlp]
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp/logs]
```

**Loki Configuration with Basic Auth:**

Configure authentication in your reverse proxy (nginx, Envoy) or API gateway rather than directly in Loki:

```nginx
location /otlp {
    auth_basic "Loki OTLP Endpoint";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://loki:3100;
}
```

### Multi-tenant Headers

For multi-tenant Loki deployments, use the `X-Scope-OrgID` header to specify the tenant:

```yaml
exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"
    headers:
      X-Scope-OrgID: "battlebots-production"
```

**Dynamic Tenant Selection:**

Route different services to different tenants using the resource processor:

```yaml
processors:
  resource:
    attributes:
      - key: loki.tenant
        from_attribute: deployment.environment
        action: insert

exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"
    headers:
      X-Scope-OrgID: "${LOKI_TENANT}"
```

**Loki Multi-tenancy Configuration:**

```yaml
auth_enabled: true

limits_config:
  # Per-tenant rate limits
  ingestion_rate_mb: 10
  ingestion_burst_size_mb: 20

  # Per-tenant OTLP configuration
  per_tenant_override_config: /etc/loki/overrides.yaml
```

**Per-tenant overrides (`/etc/loki/overrides.yaml`):**
```yaml
overrides:
  battlebots-production:
    ingestion_rate_mb: 50
    retention_period: 720h  # 30 days

  battlebots-staging:
    ingestion_rate_mb: 20
    retention_period: 168h  # 7 days
```

### TLS Configuration

**Collector with TLS:**

```yaml
exporters:
  otlphttp/logs:
    endpoint: "https://loki.battlebots.example.com/otlp"
    tls:
      insecure: false
      ca_file: /etc/ssl/certs/ca-bundle.crt
      cert_file: /etc/ssl/certs/collector-client.crt
      key_file: /etc/ssl/private/collector-client.key
      min_version: "1.3"
      server_name_override: loki.battlebots.example.com
```

**Loki TLS Configuration:**

```yaml
server:
  http_listen_port: 3100
  grpc_listen_port: 9096
  http_tls_config:
    cert_file: /etc/loki/tls/server.crt
    key_file: /etc/loki/tls/server.key
    client_auth_type: RequireAndVerifyClientCert
    client_ca_file: /etc/loki/tls/ca.crt
```

## Complete Configuration Example

This section provides a full, working configuration for integrating OpenTelemetry Collector with Loki OTLP.

### Full OpenTelemetry Collector Configuration

```yaml
# /etc/otelcol/config.yaml

extensions:
  health_check:
    endpoint: 0.0.0.0:13133

  pprof:
    endpoint: 0.0.0.0:1777

  zpages:
    endpoint: 0.0.0.0:55679

  file_storage:
    directory: /var/lib/otelcol/file_storage
    timeout: 10s

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

  filelog:
    include:
      - /var/log/battlebots/*.log
    include_file_path: true
    include_file_name: false
    operators:
      - type: json_parser
        timestamp:
          parse_from: attributes.time
          layout: '%Y-%m-%dT%H:%M:%S.%fZ'

processors:
  resourcedetection:
    detectors: [env, system, docker, kubernetes]
    timeout: 5s

  resource:
    attributes:
      - key: deployment.environment
        value: production
        action: insert
      - key: service.namespace
        value: battlebots
        action: insert

  attributes:
    actions:
      - key: loki.attribute.labels
        value: severity_text, service_name
        action: insert

  batch:
    timeout: 10s
    send_batch_size: 8192
    send_batch_max_size: 16384

  memory_limiter:
    check_interval: 1s
    limit_mib: 512
    spike_limit_mib: 128

exporters:
  otlphttp/logs:
    endpoint: "http://loki:3100/otlp"
    tls:
      insecure: true
    timeout: 30s
    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s
      max_elapsed_time: 300s
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 5000
      storage: file_storage
    compression: gzip
    headers:
      X-Scope-OrgID: "battlebots"

  debug:
    verbosity: detailed
    sampling_initial: 5
    sampling_thereafter: 200

service:
  extensions: [health_check, pprof, zpages, file_storage]

  pipelines:
    logs:
      receivers: [otlp, filelog]
      processors: [memory_limiter, resourcedetection, resource, attributes, batch]
      exporters: [otlphttp/logs]

  telemetry:
    logs:
      level: info
    metrics:
      address: 0.0.0.0:8888
```

### Full Loki Configuration with OTLP

```yaml
# /etc/loki/config.yaml

auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096
  log_level: info

common:
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: "2024-04-01"
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: loki_index_
        period: 24h

limits_config:
  allow_structured_metadata: true
  max_structured_metadata_entries_count: 128
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  ingestion_rate_mb: 50
  ingestion_burst_size_mb: 100
  per_stream_rate_limit: 10MB
  per_stream_rate_limit_burst: 20MB
  max_label_names_per_series: 15

  otlp_config:
    resource_attributes:
      attributes_config:
        - action: index_label
          attributes:
            - service.name
            - service.namespace
            - deployment.environment
            - k8s.namespace.name
            - cloud.region
        - action: structured_metadata
          attributes:
            - k8s.pod.name
            - k8s.container.name
            - service.instance.id
            - host.name

storage_config:
  tsdb_shipper:
    active_index_directory: /loki/tsdb-index
    cache_location: /loki/tsdb-cache
  filesystem:
    directory: /loki/chunks

compactor:
  working_directory: /loki/compactor
  compaction_interval: 10m
  retention_enabled: true
  retention_delete_delay: 2h
  retention_delete_worker_count: 150

querier:
  max_concurrent: 4

query_scheduler:
  max_outstanding_requests_per_tenant: 4096

frontend:
  max_outstanding_per_tenant: 4096
```

### Docker Compose Example

```yaml
version: '3.8'

services:
  loki:
    image: grafana/loki:3.0.0
    container_name: loki
    ports:
      - "3100:3100"
      - "9096:9096"
    volumes:
      - ./loki-config.yaml:/etc/loki/config.yaml
      - loki-data:/loki
    command: -config.file=/etc/loki/config.yaml
    networks:
      - battlebots-observability

  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.96.0
    container_name: otel-collector
    ports:
      - "4317:4317"   # OTLP gRPC receiver
      - "4318:4318"   # OTLP HTTP receiver
      - "8888:8888"   # Metrics endpoint
      - "13133:13133" # Health check
    volumes:
      - ./otel-config.yaml:/etc/otelcol/config.yaml
      - /var/lib/otelcol:/var/lib/otelcol
    command: ["--config=/etc/otelcol/config.yaml"]
    depends_on:
      - loki
    networks:
      - battlebots-observability

  grafana:
    image: grafana/grafana:10.4.0
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana-datasources.yaml:/etc/grafana/provisioning/datasources/datasources.yaml
    depends_on:
      - loki
    networks:
      - battlebots-observability

volumes:
  loki-data:
  grafana-data:

networks:
  battlebots-observability:
    driver: bridge
```

### Grafana Data Source Configuration

```yaml
# grafana-datasources.yaml
apiVersion: 1

datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    jsonData:
      derivedFields:
        - datasourceUid: tempo
          matcherRegex: "trace_id=(\\w+)"
          name: TraceID
          url: "$${__value.raw}"
```

### Step-by-Step Setup

1. **Create configuration files:**
   ```bash
   mkdir -p battlebots-observability
   cd battlebots-observability
   # Copy the configurations above into:
   # - loki-config.yaml
   # - otel-config.yaml
   # - grafana-datasources.yaml
   # - docker-compose.yaml
   ```

2. **Start the stack:**
   ```bash
   docker-compose up -d
   ```

3. **Verify Loki is running:**
   ```bash
   curl http://localhost:3100/ready
   # Expected: ready
   ```

4. **Verify OTel Collector is running:**
   ```bash
   curl http://localhost:13133
   # Expected: {"status":"Server available"}
   ```

5. **Send test logs:**
   ```bash
   curl -X POST http://localhost:4318/v1/logs \
     -H "Content-Type: application/json" \
     -d '{
       "resourceLogs": [{
         "resource": {
           "attributes": [{
             "key": "service.name",
             "value": {"stringValue": "test-service"}
           }]
         },
         "scopeLogs": [{
           "logRecords": [{
             "timeUnixNano": "1640000000000000000",
             "severityText": "INFO",
             "body": {"stringValue": "Test log message"}
           }]
         }]
       }]
     }'
   ```

6. **Query logs in Grafana:**
   - Open http://localhost:3000
   - Login with admin/admin
   - Navigate to Explore
   - Select Loki data source
   - Query: `{service_name="test-service"}`

## Troubleshooting

### Connection Errors

**Problem:** Collector cannot connect to Loki OTLP endpoint

**Symptoms:**
```
error exporting items: failed to push logs: Post "http://loki:3100/otlp/v1/logs": dial tcp: lookup loki: no such host
```

**Solutions:**
1. Verify Loki is running: `curl http://loki:3100/ready`
2. Check DNS resolution: `nslookup loki` (or use IP address)
3. Verify network connectivity: `telnet loki 3100`
4. Check Docker network configuration if using containers
5. Verify endpoint URL doesn't include `/v1/logs` suffix

**Problem:** 404 Not Found on OTLP endpoint

**Symptoms:**
```
error exporting items: failed to push logs: HTTP 404 Not Found
```

**Solutions:**
1. Verify Loki version is 3.0 or later: `curl http://loki:3100/loki/api/v1/status/buildinfo`
2. Check schema version is v13 in Loki config
3. Verify `allow_structured_metadata: true` in limits_config
4. Restart Loki after configuration changes

### Label Cardinality Problems

**Problem:** Too many streams causing performance degradation

**Symptoms:**
```
level=warn msg="stream limit exceeded" limit=10000 streams=15234
```

**Solutions:**
1. Review current label cardinality:
   ```bash
   curl http://localhost:3100/loki/api/v1/label
   ```

2. Check value distribution per label:
   ```bash
   curl http://localhost:3100/loki/api/v1/label/service_name/values
   ```

3. Move high-cardinality attributes to structured metadata:
   ```yaml
   limits_config:
     otlp_config:
       resource_attributes:
         attributes_config:
           - action: structured_metadata
             attributes:
               - k8s.pod.name
               - service.instance.id
   ```

4. Set cardinality limits:
   ```yaml
   limits_config:
     max_label_names_per_series: 15
     max_label_value_length: 2048
     max_label_name_length: 1024
   ```

### Structured Metadata Issues

**Problem:** Structured metadata not appearing in queries

**Symptoms:** Attributes missing from log entries in Grafana

**Solutions:**
1. Verify schema version 13 or higher
2. Check `allow_structured_metadata: true` in limits
3. Verify attributes aren't being dropped by processors
4. Query with explicit structured metadata filter:
   ```logql
   {service_name="game-server"} | k8s_pod_name="pod-123"
   ```

### Performance Issues

**Problem:** Slow queries or high memory usage

**Symptoms:** Grafana queries timeout or Loki OOM errors

**Solutions:**
1. **Enable query limits:**
   ```yaml
   limits_config:
     max_query_series: 500
     max_query_lookback: 720h
     max_entries_limit_per_query: 5000
   ```

2. **Optimize batch processor:**
   ```yaml
   processors:
     batch:
       timeout: 5s  # Reduce for faster flushing
       send_batch_size: 4096  # Smaller batches
   ```

3. **Add memory limiter:**
   ```yaml
   processors:
     memory_limiter:
       check_interval: 1s
       limit_mib: 512
   ```

4. **Use query acceleration:**
   ```yaml
   limits_config:
     bloom_gateway_enable_filtering: true
   ```

### Debugging Techniques

**Enable debug logging in Collector:**
```yaml
exporters:
  debug:
    verbosity: detailed

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [debug, otlphttp/logs]  # Add debug exporter

  telemetry:
    logs:
      level: debug  # Collector internal logs
```

**Enable debug logging in Loki:**
```yaml
server:
  log_level: debug
```

**Check Loki metrics:**
```bash
curl http://localhost:3100/metrics | grep loki_distributor
```

**Verify OTLP data structure:**
```bash
# Send test log and capture response
curl -v -X POST http://localhost:4318/v1/logs \
  -H "Content-Type: application/json" \
  -d @test-log.json
```

## BattleBots Integration Points

### Observability Stack Architecture

The Loki OTLP integration fits into the BattleBots observability stack as follows:

```
Game Servers (Go)
  └─> OpenTelemetry SDK
       └─> OTLP/gRPC (4317)
            └─> OTel Collector
                 ├─> Loki (Logs via OTLP)
                 ├─> Tempo (Traces via OTLP)
                 └─> Prometheus (Metrics via OTLP)
                      └─> Grafana (Visualization & Correlation)
```

### Collector → Loki Pipeline for Game Servers

**Recommended pipeline configuration:**

1. **Receive logs from game servers** via OTLP
2. **Detect resource attributes** (Kubernetes, cloud provider)
3. **Add BattleBots-specific labels** (environment, service namespace)
4. **Filter out verbose debug logs** in production
5. **Batch and compress** for efficiency
6. **Export to Loki** via OTLP HTTP

### Log Types and Labeling Strategy

**Game Event Logs:**
```yaml
Labels:
  service_name: "game-server"
  service_namespace: "battlebots"
  deployment_environment: "production"
  event_type: "game_event"  # Custom label

Structured Metadata:
  match_id: "uuid"
  player_count: 4
  game_mode: "elimination"
```

**System Logs:**
```yaml
Labels:
  service_name: "game-server"
  service_namespace: "battlebots"
  deployment_environment: "production"
  log_type: "system"

Structured Metadata:
  k8s_pod_name: "game-server-abc123"
  severity_text: "ERROR"
```

**Client Connection Logs:**
```yaml
Labels:
  service_name: "game-server"
  service_namespace: "battlebots"
  deployment_environment: "production"

Structured Metadata:
  client_id: "uuid"
  connection_state: "connected"
  trace_id: "trace-id"  # For correlation
```

### Example Queries for BattleBots

**Find all errors in production:**
```logql
{service_namespace="battlebots", deployment_environment="production"}
| severity_text="ERROR"
```

**Find logs for a specific match:**
```logql
{service_name="game-server"}
| match_id="550e8400-e29b-41d4-a716-446655440000"
```

**Find all logs related to a slow trace:**
```logql
{service_namespace="battlebots"}
| trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```

**Count game events by type:**
```logql
sum by (event_type) (
  count_over_time({service_name="game-server"}[5m])
)
```

## Further Reading

### OTLP Specification
- [OpenTelemetry Protocol Specification](https://opentelemetry.io/docs/specs/otlp/) - Complete OTLP specification
- [OTLP Logs Data Model](https://opentelemetry.io/docs/specs/otel/logs/) - OpenTelemetry logs data model
- [OTLP Exporter Configuration](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/) - SDK configuration guide

### Grafana Loki OTLP Documentation
- [Ingesting logs to Loki using OpenTelemetry Collector](https://grafana.com/docs/loki/latest/send-data/otel/) - Official Loki OTLP guide
- [Getting started with OTel Collector and Loki](https://grafana.com/docs/loki/latest/send-data/otel/otel-collector-getting-started/) - Step-by-step tutorial
- [Native OTLP vs Loki Exporter](https://grafana.com/docs/loki/latest/send-data/otel/native_otlp_vs_loki_exporter/) - Migration guide
- [Loki 3.0 Release Notes](https://grafana.com/blog/2024/04/09/grafana-loki-3.0-release-all-the-new-features/) - Features announcement
- [Migration to Native OTLP Format](https://grafana.com/docs/grafana-cloud/send-data/otlp/adopt-new-logs-format/) - Cloud migration guide

### OpenTelemetry Collector Documentation
- [Collector Configuration](https://opentelemetry.io/docs/collector/configuration/) - General configuration guide
- [Collector Resiliency](https://opentelemetry.io/docs/collector/resiliency/) - Retry and queue configuration
- [Configuration Best Practices](https://opentelemetry.io/docs/security/config-best-practices/) - Security and performance
- [Exporter Helper Documentation](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md) - Retry and queue details

### Loki Configuration and Operations
- [Understanding Labels](https://grafana.com/docs/loki/latest/get-started/labels/) - Label strategy guide
- [Structured Metadata](https://grafana.com/docs/loki/latest/get-started/labels/structured-metadata/) - Structured metadata overview
- [Loki Schema Configuration](https://grafana.com/docs/loki/latest/configure/) - Schema v13 setup
- [Upgrade to Loki 3.0](https://grafana.com/docs/loki/latest/setup/upgrade/) - Migration instructions

### Grafana Integration
- [Trace Integration in Explore](https://grafana.com/docs/grafana/latest/explore/trace-integration/) - Log-trace correlation
- [Trace Correlations](https://grafana.com/docs/grafana/latest/datasources/tempo/traces-in-grafana/trace-correlations/) - Setting up correlations
- [Configure Loki Data Source](https://grafana.com/docs/grafana/latest/datasources/loki/) - Data source configuration

### Community Guides and Tutorials
- [Grafana Loki 101: Ingesting with OTel Collector](https://grafana.com/blog/2025/02/24/grafana-loki-101-how-to-ingest-logs-with-alloy-or-the-opentelemetry-collector/) - Official blog tutorial
- [Building a Logging Pipeline with OTel, Loki, and Grafana](https://dev.to/tingwei628/how-to-build-a-logging-pipeline-with-opentelemetry-grafana-loki-and-grafana-in-docker-compose-4kk) - Docker Compose guide
- [Logging with OpenTelemetry and Loki](https://bit.kevinslin.com/p/logging-with-opentelemetry-and-loki) - Practical implementation
- [Efficient Application Log Collection](https://addozhang.medium.com/efficient-application-log-collection-and-analysis-using-opentelemetry-and-loki-baf04bc4a8a2) - Analysis guide

### Troubleshooting and Best Practices
- [OTel Batching Best Practices](https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/best-practices/opentelemetry-best-practices-batching/) - Batch configuration
- [Collector Persistence and Retry](https://axoflow.com/blog/opentelemetry-collector-persistence-and-retry-mechanisms-under-the-hood) - Deep dive
- [Label Cardinality Issues](https://github.com/grafana/loki/issues/2863) - Common problems and solutions

### Related BattleBots Documentation
- [OpenTelemetry Collector: Logs Support](/research_and_development/analysis/observability/otel-collector/otel-collector-logs.md) - Collector log processing
- [OpenTelemetry Collector: Overview](/research_and_development/analysis/observability/otel-collector/opentelemetry-collector-overview.md) - Collector architecture
