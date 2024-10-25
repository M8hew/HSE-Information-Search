package lsmtree

import (
	"lsm-tree/internal/level"
	"lsm-tree/internal/memtable"
)

type LSMTree interface {
	// add key value pair
	Put(key, value string)
	// returns value by key, nil if no such key
	Get(key string) *string
	// deletes value by key -- return false if key is absent
	Delete(key string)
}

type lsmTree struct {
	memTable memtable.Memtable
	levels   []level.Level
}

func NewLSMTree( /* cfg */ ) LSMTree {
	// TODO
	return &lsmTree{
		memTable: memtable.NewMemtable(),
		levels:   make([]level.Level, 0),
	}
}

func (lsm *lsmTree) Put(key, value string) {
	// TODO
}

func (lsm *lsmTree) Get(key string) *string {
	// TODO
	return nil
}

func (lsm *lsmTree) Delete(key string) {
	// TODO
}
