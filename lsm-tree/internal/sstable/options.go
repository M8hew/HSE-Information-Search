package sstable

import (
	bloom "lsm_tree/internal/bloom_filter"
)

type fnOpt func(m *modifications)

func withBloom(m *modifications) {
	m.bloom, _ = bloom.NewBloomFilter(1000, 0.01, bloom.NewMurMurHasher())
}
