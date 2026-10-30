package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDay(t *testing.T) {
	got, err := parseDay("2026-11-15")
	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC), got)

	invalid := []string{"15/11/2026", "2026-13-01", "", "tomorrow"}
	for _, s := range invalid {
		_, err := parseDay(s)
		assert.Error(t, err, "should reject %q", s)
	}

}
