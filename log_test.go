// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bigqueryexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/testdata"
)

func TestLogsToRows(t *testing.T) {
	ld := testdata.GenerateLogs(1)
	rows := logsToRows(ld)
	require.Len(t, rows, 1)

	row := rows[0]
	assert.NotNil(t, row["body"])
	assert.NotEmpty(t, row["severity_text"])
	assert.NotEmpty(t, row["trace_id"])
	assert.NotEmpty(t, row["span_id"])
	assert.NotEmpty(t, row["resource_attributes"])
	assert.IsType(t, int64(0), row["dropped_attributes_count"])
	assert.NotNil(t, row["resource_schema_url"])
	assert.NotNil(t, row["scope_schema_url"])
}

func TestLogsToRowsMultiple(t *testing.T) {
	ld := testdata.GenerateLogs(4)
	rows := logsToRows(ld)
	require.Len(t, rows, 4)

	assert.NotNil(t, rows[0]["body"])
	assert.NotNil(t, rows[1]["body"])
}

func TestLogsToRowsEmpty(t *testing.T) {
	assert.Empty(t, logsToRows(plog.NewLogs()))
}
