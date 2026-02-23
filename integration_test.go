// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bigqueryexporter

import (
	"strings"
	"testing"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/pdata/testdata"
)

func TestIntegration_ExporterLifecycleAndWrites(t *testing.T) {
	fx := newIntegrationFixture(t)
	defer fx.cleanup(t)

	t.Run("dataset missing fails start", func(t *testing.T) {
		cfg := createDefaultConfig()
		cfg.Dataset.Project = fx.projectID
		cfg.Dataset.ID = temporaryDatasetID()

		exp := newBigQueryExporter(t.Context(), cfg, zap.NewNop())

		err := exp.start(t.Context(), nil)
		if err == nil {
			t.Fatal("start expected error, got nil")
		}
		if !strings.Contains(err.Error(), "dataset") {
			t.Fatalf("start error = %q, want dataset error", err.Error())
		}
	})

	t.Run("creates default tables and writes all signals", func(t *testing.T) {
		cfg := createDefaultConfig()
		cfg.Dataset.Project = fx.projectID
		cfg.Dataset.ID = fx.datasetID

		exp := newBigQueryExporter(t.Context(), cfg, zap.NewNop())
		if err := exp.start(t.Context(), nil); err != nil {
			t.Fatalf("start exporter: %v", err)
		}
		defer func() {
			if err := exp.shutdown(t.Context()); err != nil {
				t.Fatalf("shutdown exporter: %v", err)
			}
		}()

		for _, table := range []string{cfg.Dataset.Table.Trace, cfg.Dataset.Table.Metric, cfg.Dataset.Table.Log} {
			exists, err := fx.tableExists(table)
			if err != nil {
				t.Fatalf("tableExists(%s): %v", table, err)
			}
			if !exists {
				t.Fatalf("expected table %q to exist", table)
			}
		}

		traceData := testdata.GenerateTraces(5)
		metricData := testdata.GenerateMetricsAllTypes()
		logData := testdata.GenerateLogs(5)

		if err := exp.pushTraces(t.Context(), traceData); err != nil {
			t.Fatalf("push traces: %v", err)
		}
		if err := exp.pushMetrics(t.Context(), metricData); err != nil {
			t.Fatalf("push metrics: %v", err)
		}
		if err := exp.pushLogs(t.Context(), logData); err != nil {
			t.Fatalf("push logs: %v", err)
		}

		fx.waitForRows(t, cfg.Dataset.Table.Trace, int64(len(tracesToRows(traceData))))
		fx.waitForRows(t, cfg.Dataset.Table.Metric, int64(len(metricsToRows(metricData))))
		fx.waitForRows(t, cfg.Dataset.Table.Log, int64(len(logsToRows(logData))))
	})

	t.Run("respects custom table names and accumulates multiple writes", func(t *testing.T) {
		cfg := createDefaultConfig()
		cfg.Dataset.Project = fx.projectID
		cfg.Dataset.ID = fx.datasetID
		cfg.Dataset.Table.Trace = "trace_custom"
		cfg.Dataset.Table.Metric = "metric_custom"
		cfg.Dataset.Table.Log = "log_custom"

		exp := newBigQueryExporter(t.Context(), cfg, zap.NewNop())
		if err := exp.start(t.Context(), nil); err != nil {
			t.Fatalf("start exporter: %v", err)
		}
		defer func() {
			if err := exp.shutdown(t.Context()); err != nil {
				t.Fatalf("shutdown exporter: %v", err)
			}
		}()

		for _, table := range []string{cfg.Dataset.Table.Trace, cfg.Dataset.Table.Metric, cfg.Dataset.Table.Log} {
			exists, err := fx.tableExists(table)
			if err != nil {
				t.Fatalf("tableExists(%s): %v", table, err)
			}
			if !exists {
				t.Fatalf("expected custom table %q to exist", table)
			}
		}

		totalTraceRows := 0
		totalMetricRows := 0
		totalLogRows := 0

		for i := range 2 {
			traceData := testdata.GenerateTraces(3)
			metricData := testdata.GenerateMetrics(3)
			logData := testdata.GenerateLogs(3)

			if err := exp.pushTraces(t.Context(), traceData); err != nil {
				t.Fatalf("push traces batch %d: %v", i, err)
			}
			if err := exp.pushMetrics(t.Context(), metricData); err != nil {
				t.Fatalf("push metrics batch %d: %v", i, err)
			}
			if err := exp.pushLogs(t.Context(), logData); err != nil {
				t.Fatalf("push logs batch %d: %v", i, err)
			}

			totalTraceRows += len(tracesToRows(traceData))
			totalMetricRows += len(metricsToRows(metricData))
			totalLogRows += len(logsToRows(logData))
		}

		fx.waitForRows(t, cfg.Dataset.Table.Trace, int64(totalTraceRows))
		fx.waitForRows(t, cfg.Dataset.Table.Metric, int64(totalMetricRows))
		fx.waitForRows(t, cfg.Dataset.Table.Log, int64(totalLogRows))
	})
}
