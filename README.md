# OpenTelemetry BigQuery Exporter

Export OpenTelemetry traces, metrics, and logs to
[Google BigQuery](https://cloud.google.com/bigquery) using the
[BigQuery Storage Write API](https://cloud.google.com/bigquery/docs/write-api).

This repository contains a standalone version of the `bigquery` exporter
originally developed for OpenTelemetry Collector contrib.

## Status

- Stability: `development` for traces, metrics, and logs
- Signal support: traces, metrics, logs
- Maintainer: [@glacion](https://github.com/glacion)

## What It Does

- Writes trace, metric, and log rows into BigQuery tables
- Creates destination tables automatically if they do not exist
- Uses ingestion-time partitioning
- Supports exporter timeout, queueing, and retry controls

## Requirements

- A Google Cloud project with BigQuery enabled
- An existing BigQuery dataset
- Credentials via
  [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials)
- Go 1.25+ (for local development in this repository)

## Collector Configuration

Exporter type: `bigquery`

```yaml
exporters:
  bigquery:
    dataset:
      project: my-gcp-project
      id: otel_dataset
      trace_table: trace
      metric_table: metric
      log_table: log

    timeout: 30s

    retry_on_failure:
      enabled: true
      initial_interval: 5s
      max_interval: 30s

    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 1000
```

Minimal configuration (project inferred from ADC/env):

```yaml
exporters:
  bigquery:
    dataset:
      id: otel_dataset
```

### Configuration Fields

| Field | Type | Default | Required | Notes |
|---|---|---|---|---|
| `dataset.project` | string | empty | No | If omitted, project is resolved from `GOOGLE_CLOUD_PROJECT`, `GCLOUD_PROJECT`, `GCP_PROJECT`, or ADC metadata |
| `dataset.id` | string | none | Yes | BigQuery dataset ID |
| `dataset.trace_table` | string | `trace` | No | Trace table name |
| `dataset.metric_table` | string | `metric` | No | Metric table name |
| `dataset.log_table` | string | `log` | No | Log table name |
| `timeout` | duration | `30s` | No | Per-export call timeout |
| `retry_on_failure` | object | enabled | No | Export retry policy |
| `sending_queue` | object | disabled | No | Queue and batching behavior |

Identifier validation for dataset and table names:

- Pattern: `^[A-Za-z_][A-Za-z0-9_]*$`
- Maximum length: `1024`

## BigQuery Schema (High Level)

The exporter writes one row per OTLP record/data point/log record and includes
signal-specific fields plus resource/scope context.

- **Traces**: IDs, span name/kind, status, timing, attributes, events, links
- **Metrics**: metric metadata, datapoint values, histogram/summary fields,
  exemplars, attributes
- **Logs**: timestamps, severity, body, trace/span correlation, attributes

Table creation and column mapping are implemented in `storage_writer.go`.

## Running Tests

Unit tests:

```sh
go test ./...
```

Integration tests (requires live GCP + BigQuery + ADC):

```sh
RUN_BIGQUERY_INTEGRATION=1 go test -tags integration -run TestIntegration -v -count=1 ./...
```

Optional environment override for integration tests:

- `BIGQUERY_PROJECT`

## Repository Layout

- `factory.go` / `config.go`: component wiring and configuration validation
- `traces.go`, `metrics.go`, `logs.go`: signal-to-row transformations
- `storage_writer.go`: table management and write path
- `integration_test.go`: end-to-end BigQuery verification
- `metadata.yaml`: exporter metadata/stability

## Module Path Note

This standalone repository now uses:

`github.com/glacion/opentelemetry-collector-bigquery-exporter`

If you publish under a different org/repository name, update:

1. `go.mod` module path
2. Package import comments in `*.go` files (for example `config.go`, `factory.go`)
3. Any downstream references that import this module path
