package sstable

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bloom "lsm_tree/internal/bloom_filter"
	"lsm_tree/internal/common"
)

type (
	String = common.MyString
)

// Mock implementation of Serializable for testing
type MockSerializable struct {
	Value string
}

func TestNewSSTable_Creation(t *testing.T) {
	filename := "test_sstable_creation.sst"
	defer os.Remove(filename)

	sst, err := NewSSTable[*String](filename, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, sst)

	defer sst.Close()

	// Verify that the file is created
	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

func TestAddAndRetrieveEntries(t *testing.T) {
	filename := "test_sstable_add_retrieve.sst"
	defer os.Remove(filename)

	sst, err := NewSSTable[*String](filename, 1024)
	assert.NoError(t, err)
	defer sst.Close()

	err = sst.Add("key1", common.ToMyString("value1"))
	assert.NoError(t, err)
	err = sst.Add("key2", common.ToMyString("value2"))
	assert.NoError(t, err)

	err = sst.Flush()
	assert.NoError(t, err)

	value, found := sst.Get("key1")
	assert.True(t, found)
	assert.Equal(t, common.ToMyString("value1"), value)

	value, found = sst.Get("key2")
	assert.True(t, found)
	assert.Equal(t, common.ToMyString("value2"), value)

	value, found = sst.Get("key3")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestBloomFilterIntegration(t *testing.T) {
	filename := "test_sstable_bloom_filter.sst"
	defer os.Remove(filename)

	bloom, err := bloom.NewBloomFilter(100, 3, bloom.NewMurMurHasher())
	require.NoError(t, err)

	sst, err := NewSSTable[*String](filename, 1024, func(mod *modifications) {
		mod.bloom = bloom
	})
	assert.NoError(t, err)
	defer sst.Close()

	err = sst.Add("key1", common.ToMyString("value1"))
	assert.NoError(t, err)
	err = sst.Add("key2", common.ToMyString("value2"))
	assert.NoError(t, err)

	assert.True(t, bloom.Contains("key1"))
	assert.False(t, bloom.Contains("key3"))

	value, found := sst.Get("key1")
	assert.True(t, found)
	assert.Equal(t, common.ToMyString("value1"), value)

	value, found = sst.Get("key3")
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestSSTableReopenAndIndexLoading(t *testing.T) {
	filename := "test_sstable_reopen.sst"
	defer os.Remove(filename)

	sst, err := NewSSTable[*String](filename, 1024)
	assert.NoError(t, err)

	err = sst.Add("key1", common.ToMyString("value1"))
	assert.NoError(t, err)
	err = sst.Add("key2", common.ToMyString("value2"))
	assert.NoError(t, err)

	err = sst.Close()
	assert.NoError(t, err)

	sst, err = NewSSTable[*String](filename, 1024)
	assert.NoError(t, err)
	defer sst.Close()

	value, found := sst.Get("key1")
	assert.True(t, found)
	assert.Equal(t, common.ToMyString("value1"), value)

	value, found = sst.Get("key2")
	assert.True(t, found)
	assert.Equal(t, common.ToMyString("value2"), value)
}
