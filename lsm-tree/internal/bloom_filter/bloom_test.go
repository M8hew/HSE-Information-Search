package bloom_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bloom "lsm-tree/internal/bloom_filter"
)

func TestBloomFilter(t *testing.T) {
	filter, err := bloom.NewBloomFilter(1000, 0.01, bloom.NewMurMurHasher())
	require.NoError(t, err)

	assert.False(t, filter.Contains("apple"))

	filter.Add("apple")
	filter.Add("grapefruit")
	filter.Add("kiwi")

	assert.True(t, filter.Contains("grapefruit"))
	assert.True(t, filter.Contains("kiwi"))
	assert.False(t, filter.Contains("lime"))
}
