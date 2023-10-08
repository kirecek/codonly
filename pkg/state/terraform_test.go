package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformProvider(t *testing.T) {

	provider, err := NewTerraformProviderFromStateOutput("./testdata/state.output.json")
	assert.Nil(t, err)

	tt := map[string]struct {
		name     string
		r        *Resource
		expected bool
	}{
		"testdata contain sql-instance": {
			r:        &Resource{IDKey: "id", IDValue: "test-instance", Type: "google_sql_database_instance"},
			expected: true,
		},
		"sql instance not present in state": {
			r:        &Resource{IDKey: "id", IDValue: "non-existing-instance", Type: "google_sql_database_instance"},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, provider.Contains(tc.r))
		})
	}
}
