// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bigqueryexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"go.opentelemetry.io/collector/pdata/testdata"
)

func TestMetricsToRowsAllTypes(t *testing.T) {
	md := testdata.GenerateMetricsAllTypes()
	rows := metricsToRows(md)
	require.NotEmpty(t, rows)

	types := map[string]int{}
	for _, r := range rows {
		types[r["metric_type"].(string)]++
	}
	assert.Greater(t, types["GAUGE"], 0)
	assert.Greater(t, types["SUM"], 0)
	assert.Greater(t, types["HISTOGRAM"], 0)
	assert.Greater(t, types["SUMMARY"], 0)

	for _, r := range rows {
		assert.NotEmpty(t, r["resource_attributes"])
		assert.NotNil(t, r["resource_schema_url"])
		assert.NotNil(t, r["scope_schema_url"])
		assert.IsType(t, int64(0), r["flags"])
	}
}

func TestMetricsToRowsGaugeValues(t *testing.T) {
	md := testdata.GenerateMetrics(1)
	rows := metricsToRows(md)
	require.NotEmpty(t, rows)

	for _, r := range rows {
		assert.NotEmpty(t, r["metric_type"])
		assert.True(t, r["value_double"] != nil || r["value_int"] != nil || r["sum"] != nil || r["count"] != nil)
	}
}

func TestMetricsToRowsEmpty(t *testing.T) {
	assert.Empty(t, metricsToRows(pmetric.NewMetrics()))
}

func TestMetricsJSONDefaults(t *testing.T) {
	assert.Equal(t, "[]", bucketCountsToJSON(nil))
	assert.Equal(t, "[]", explicitBoundsToJSON(nil))
	assert.Equal(t, "[]", quantilesToJSON(pmetric.NewSummaryDataPointValueAtQuantileSlice()))
	assert.Equal(t, "[]", exemplarsToJSON(pmetric.NewExemplarSlice()))
}
