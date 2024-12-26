package storage

import (
	"errors"
	"fmt"
	"os"

	"lsm_tree/internal/common"
	"lsm_tree/internal/memtable"
	"lsm_tree/internal/sstable"
)

type BalancedTree[T common.Serializable] interface {
	Size() int
	Set(key string, value T)
	Get(key string) (T, bool)
	Iterate(fn func(key string, value T))
}

type SSTable[T common.Serializable] interface {
	Add(key string, value T) error
	Get(key string) (T, bool)
	Flush() error
}

type StorageComponent[T common.Serializable] struct {
	config   Config
	memTable BalancedTree[T]
	ssTables []SSTable[T]
}

type Config struct {
	MemTableSizeThreshold int
	BlockSize             int
}

func NewStorageComponent[T common.Serializable](cfg Config) StorageComponent[T] {
	if cfg.BlockSize == 0 {
		cfg.BlockSize = 4096
	}
	if cfg.MemTableSizeThreshold == 0 {
		cfg.MemTableSizeThreshold = 100
	}
	return StorageComponent[T]{
		config:   cfg,
		memTable: memtable.NewMemTable[T](),
	}
}

func (sc *StorageComponent[T]) Set(key string, value T) error {
	sc.memTable.Set(key, value)
	if sc.memTable.Size() >= sc.config.MemTableSizeThreshold {
		return sc.flushMemtable()
	}
	return nil
}

func (sc *StorageComponent[T]) Get(key string) (T, error) {
	if value, found := sc.memTable.Get(key); found {
		return value, nil
	}

	for _, ssTable := range sc.ssTables {
		if value, found := ssTable.Get(key); found {
			return value, nil
		}
	}

	var empty T
	return empty, errors.New("key not found")
}

func (sc *StorageComponent[T]) flushMemtable() error {
	filename := fmt.Sprintf("sstable_%d", len(sc.ssTables)+1)

	ssTable, err := sstable.NewSSTable[T](filename, int64(sc.config.BlockSize))
	if err != nil {
		return fmt.Errorf("failed to create SSTable: %w", err)
	}

	sc.memTable.Iterate(func(key string, value T) {
		err := ssTable.Add(key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to add key-value pair (%s, %v) to SSTable: %v\n", key, value, err)
		}
	})

	if err := ssTable.Flush(); err != nil {
		return fmt.Errorf("failed to flush SSTable: %w", err)
	}

	sc.ssTables = append(sc.ssTables, ssTable)

	sc.memTable = memtable.NewMemTable[T]()
	return nil
}
