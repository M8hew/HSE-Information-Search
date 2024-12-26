package memtable

import (
	"lsm_tree/internal/common"

	"github.com/google/btree"
)

type MemTable[T common.Serializable] struct {
	btree.BTree
}

func NewMemTable[T common.Serializable]() *MemTable[T] {
	return &MemTable[T]{
		*btree.New(32),
	}
}

func (mt *MemTable[T]) Set(key string, value T) {
	mt.BTree.ReplaceOrInsert(Value[T]{
		key:   key,
		value: value,
	})
}

func (mt *MemTable[T]) Get(key string) (T, bool) {
	item := mt.BTree.Get(Value[T]{key: key})
	if item == nil {
		var result T
		return result, false
	}
	return item.(Value[T]).value, true
}

func (mt *MemTable[T]) Size() int {
	return mt.BTree.Len()
}

func (mt *MemTable[T]) Iterate(fn func(key string, value T)) {
	mt.Ascend(func(item btree.Item) bool {
		value := item.(Value[T])
		fn(value.key, value.value)
		return true
	})
}
