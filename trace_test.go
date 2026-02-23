// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package bigqueryexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/pdata/testdata"
)

func TestTracesToRows(t *testing.T) {
	td := testdata.GenerateTraces(1)
	rows := tracesToRows(td)
	require.Len(t, rows, 1)

	row := rows[0]
	assert.NotEmpty(t, row["name"])
	assert.NotEmpty(t, row["status_code"])
	assert.Contains(t, row["resource_attributes"].(string), "resource-attr")
	assert.NotEmpty(t, row["events"])
	assert.Equal(t, int64(0), row["flags"])
	assert.IsType(t, int64(0), row["dropped_attributes_count"])
	assert.IsType(t, int64(0), row["dropped_events_count"])
	assert.IsType(t, int64(0), row["dropped_links_count"])
	assert.NotNil(t, row["resource_schema_url"])
	assert.NotNil(t, row["scope_schema_url"])
}

func TestTracesToRowsMultipleSpans(t *testing.T) {
	td := testdata.GenerateTraces(2)
	rows := tracesToRows(td)
	require.Len(t, rows, 2)

	assert.NotEmpty(t, rows[0]["name"])
	assert.NotEmpty(t, rows[1]["name"])
	assert.NotNil(t, rows[1]["links"])
}

func TestTracesToRowsMultipleResources(t *testing.T) {
	td := testdata.GenerateTraces(3)
	rows := tracesToRows(td)
	require.Len(t, rows, 3)
}

func TestTracesToRowsEmpty(t *testing.T) {
	assert.Empty(t, tracesToRows(testdata.GenerateTraces(0)))
}
