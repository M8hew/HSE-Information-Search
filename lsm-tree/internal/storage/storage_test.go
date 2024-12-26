package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"lsm_tree/internal/common"
)

type (
	MockMemTable[T common.Serializable] struct {
		data map[string]T
	}

	MockSSTable[T common.Serializable] struct {
		data map[string]T
	}

	String = common.MyString
)

func NewMockMemTable[T common.Serializable]() *MockMemTable[T] {
	return &MockMemTable[T]{data: make(map[string]T)}
}

func (m *MockMemTable[T]) Size() int {
	return len(m.data)
}

func (m *MockMemTable[T]) Set(key string, value T) {
	m.data[key] = value
}

func (m *MockMemTable[T]) Get(key string) (T, bool) {
	val, found := m.data[key]
	return val, found
}

func (m *MockMemTable[T]) Iterate(fn func(key string, value T)) {
	for k, v := range m.data {
		fn(k, v)
	}
}

func NewMockSSTable[T common.Serializable]() *MockSSTable[T] {
	return &MockSSTable[T]{data: make(map[string]T)}
}

func (m *MockSSTable[T]) Add(key string, value T) error {
	m.data[key] = value
	return nil
}

func (m *MockSSTable[T]) Get(key string) (T, bool) {
	val, found := m.data[key]
	return val, found
}

func (m *MockSSTable[T]) Flush() error {
	return nil
}

func TestSetAndGet(t *testing.T) {
	memTable := NewMockMemTable[*String]()
	ssTable := NewMockSSTable[*String]()

	sc := StorageComponent[*String]{
		config:   Config{MemTableSizeThreshold: 2, BlockSize: 1024},
		memTable: memTable,
		ssTables: []SSTable[*String]{ssTable},
	}

	// Test setting and getting from the memtable
	value := common.ToMyString("value1")
	err := sc.Set("key1", value)
	assert.NoError(t, err)
	err = sc.Set("key2", value)
	assert.NoError(t, err)

	result, err := sc.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetFromSSTables(t *testing.T) {
	memTable := NewMockMemTable[*String]()
	ssTable := NewMockSSTable[*String]()
	ssTable.Add("key2", common.ToMyString("value2"))

	sc := StorageComponent[*String]{
		config:   Config{MemTableSizeThreshold: 2, BlockSize: 1024},
		memTable: memTable,
		ssTables: []SSTable[*String]{ssTable},
	}

	// Test retrieving a value from the SSTable
	result, err := sc.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, common.ToMyString("value2"), result)

	_, err = sc.Get("key3")
	assert.Error(t, err)
	assert.Equal(t, "key not found", err.Error())
}
