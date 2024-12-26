package bloom

import (
	"hash"

	"github.com/spaolacci/murmur3"
)

type MurMurHasher struct{}

func NewMurMurHasher() *MurMurHasher {
	return &MurMurHasher{}
}

func (mh *MurMurHasher) GetHashes(n uint) []hash.Hash64 {
	hashes := make([]hash.Hash64, n)
	for i := range hashes {
		hashes[i] = murmur3.New64WithSeed(uint32(i))
	}
	return hashes
}
