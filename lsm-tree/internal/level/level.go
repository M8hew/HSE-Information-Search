package level

import "lsm-tree/internal/sstable"

type Level interface {
	AddSSTable(sst sstable.SSTable) error
	Get(key string) *string
}
