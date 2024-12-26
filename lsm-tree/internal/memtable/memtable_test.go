package memtable

import (
	"lsm_tree/internal/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemTable(t *testing.T) {
	mt := NewMemTable[*common.MyInt]()
	assert.NotNil(t, mt, "MemTable should be initialized")
	assert.Equal(t, 0, mt.Size(), "New MemTable should have size 0")
}

func TestMemTableSetAndGet(t *testing.T) {
	mt := NewMemTable[*common.MyInt]()

	// Test inserting a value
	mt.Set("key1", common.ToMyInt(42))
	value, found := mt.Get("key1")
	assert.True(t, found, "key1 should be found in MemTable")
	assert.Equal(t, common.ToMyInt(42), value, "The value for key1 should be 42")

	// Test updating the same key with a new value
	mt.Set("key1", common.ToMyInt(100))
	value, found = mt.Get("key1")
	assert.True(t, found, "key1 should still be found after update")
	assert.Equal(t, common.ToMyInt(100), value, "The value for key1 should be updated to 100")

	// Test retrieving a non-existent key
	_, found = mt.Get("non_existent_key")
	assert.False(t, found, "non_existent_key should not be found in MemTable")
}

func TestMemTableSize(t *testing.T) {
	mt := NewMemTable[*common.MyString]()

	assert.Equal(t, 0, mt.Size(), "MemTable size should initially be 0")

	mt.Set("key1", common.ToMyString("value1"))
	assert.Equal(t, 1, mt.Size(), "MemTable size should be 1 after adding one item")

	mt.Set("key2", common.ToMyString("value2"))
	assert.Equal(t, 2, mt.Size(), "MemTable size should be 2 after adding a second item")

	mt.Set("key1", common.ToMyString("newValue1"))
	assert.Equal(t, 2, mt.Size(), "MemTable size should remain 2 when updating an existing key")
}

func TestMemTableIterate(t *testing.T) {
	mt := NewMemTable[*common.MyInt]()

	mt.Set("key1", common.ToMyInt(1))
	mt.Set("key2", common.ToMyInt(2))
	mt.Set("key3", common.ToMyInt(3))

	keys := []string{}
	values := []*common.MyInt{}

	mt.Iterate(func(key string, value *common.MyInt) {
		keys = append(keys, key)
		values = append(values, value)
	})

	assert.Equal(t, []string{"key1", "key2", "key3"}, keys, "Keys should be iterated in sorted order")
	assert.Equal(t, []*common.MyInt{
		common.ToMyInt(1),
		common.ToMyInt(2),
		common.ToMyInt(3),
	}, values, "Values should match the inserted values in sorted order")
}

func TestValueLess(t *testing.T) {
	v1 := Value[*common.MyInt]{key: "key1", value: common.ToMyInt(1)}
	v2 := Value[*common.MyInt]{key: "key2", value: common.ToMyInt(2)}
	v3 := Value[*common.MyInt]{key: "key1", value: common.ToMyInt(3)} // Same key as v1, different value

	assert.True(t, v1.Less(v2), "v1 should be less than v2 because key1 < key2")
	assert.False(t, v2.Less(v1), "v2 should not be less than v1 because key2 > key1")
	assert.False(t, v1.Less(v3), "v1 should not be less than v3 because keys are equal")
	assert.False(t, v3.Less(v1), "v3 should not be less than v1 because keys are equal")
}
