package bloom

import (
	"errors"
	"hash"
	"math"
)

type Filter interface {
	Add(key string)
	Test(key string) bool
}

type Hasher interface {
	GetHashes(n uint) []hash.Hash64
}

type BloomFilter struct {
	bitSet []byte
	hashes []hash.Hash64
}

// n -- number of elements, p -- false positive rate
func NewBloomFilter(n uint64, p float64, h Hasher) (*BloomFilter, error) {
	if n == 0 {
		return nil, errors.New("attemp to create filter with 0 capacity")
	}
	if p <= 0 || p >= 1 {
		return nil, errors.New("incorrect value was passed as p")
	}
	if h == nil {
		return nil, errors.New("nil was passed as hasher value")
	}

	bitSetLen, hashesNum := getOptimalParams(n, p)
	return &BloomFilter{
		bitSet: make([]byte, (bitSetLen+7)/8),
		hashes: h.GetHashes(hashesNum),
	}, nil
}

func getOptimalParams(n uint64, p float64) (bitSetLen uint64, hashesNum uint) {
	bitSetLen = uint64(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	if bitSetLen == 0 {
		bitSetLen = 1
	}

	hashesNum = uint(math.Ceil((float64(bitSetLen) / float64(n)) * math.Log(2)))
	if hashesNum == 0 {
		hashesNum = 1
	}
	return
}

func (bf *BloomFilter) Add(key string) {
	for _, hash := range bf.hashes {
		hash.Reset()
		hash.Write([]byte(key))
		hashVal := hash.Sum64() % uint64(len(bf.bitSet))
		bf.setBit(hashVal)
	}
}

func (bf *BloomFilter) Contains(key string) bool {
	for _, hash := range bf.hashes {
		hash.Reset()
		hash.Write([]byte(key))
		hashVal := hash.Sum64() % uint64(len(bf.bitSet))
		if !bf.checkBit(hashVal) {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) setBit(ind uint64) {
	bitPos := ind % 8
	bf.bitSet[ind/8] |= (1 << bitPos)
}

func (bf *BloomFilter) checkBit(ind uint64) bool {
	bitPos := ind % 8
	return (bf.bitSet[ind/8] & (1 << bitPos)) != 0
}
