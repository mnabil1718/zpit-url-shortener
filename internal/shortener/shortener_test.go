package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShorten_Length(t *testing.T) {
	cases := []int{6, 8, 10, 12}
	for _, n := range cases {
		c, err := Shorten(n)
		require.NoError(t, err, "actual length: %d", len(c))
		assert.Len(t, c, n, "expected length: %d", n)
	}
}

func TestShorten_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for range 100 {
		c, err := Shorten(6)
		require.NoError(t, err)
		assert.False(t, seen[c], "code duplicated %s", c)
		seen[c] = true
	}
}
