package memtable

import (
	"github.com/google/btree"

	"lsm_tree/internal/common"
)

type Value[T common.Serializable] struct {
	key   string
	value T
}

func (v Value[T]) Less(than btree.Item) bool {
	other, ok := than.(Value[T])
	if !ok {
		panic("cannot compare Value with different types")
	}
	return v.key < other.key
}
