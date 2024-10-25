package memtable

import "lsm-tree/internal/sstable"

type Memtable interface {
	Put(key, value string)
	Get(key string) *string
	Delete(key string)
	Flush() (sst sstable.SSTable)
}

type memtable struct {
	data map[string]string
}

func NewMemtable() Memtable {
	return &memtable{
		data: make(map[string]string),
	}
}

func (m *memtable) Put(key, value string) {
	m.data[key] = value
}

func (m *memtable) Get(key string) *string {
	value, ok := m.data[key]
	if !ok {
		return nil
	}
	return &value
}

func (m *memtable) Delete(key string) {
	delete(m.data, key)
}

func (m *memtable) Flush() (sst sstable.SSTable) {
	// TODO
	return nil
}
