package bloom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBloomFilter(t *testing.T) {
	filter, err := NewBloomFilter(1000, 0.01, NewMurMurHasher())
	require.NoError(t, err)

	assert.False(t, filter.Contains("apple"))

	filter.Add("apple")
	filter.Add("grapefruit")
	filter.Add("kiwi")

	assert.True(t, filter.Contains("grapefruit"))
	assert.True(t, filter.Contains("kiwi"))
	assert.False(t, filter.Contains("lime"))
}
