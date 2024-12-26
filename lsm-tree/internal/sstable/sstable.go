package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	bloom "lsm_tree/internal/bloom_filter"
	"lsm_tree/internal/common"
)

type (
	/*
		SSTable file layout:
		[Block1][Block2]...[BlockN]

		Block layout:
		[size1 : int64][entry1 : struct{key, value}][size2][entry2]...[sizeN][entryN]
	*/
	SSTable[T common.Serializable] struct {
		filename  string
		blockSize int64

		file  *os.File
		index []indElem

		blockBuf bytes.Buffer

		modifications
	}

	indElem struct {
		key    string
		offset int64
	}

	modifications struct {
		bloom *bloom.BloomFilter
	}

	Serializable interface {
		Serialize() ([]byte, error)
		Deserialize([]byte) error
	}
)

// Creates new file or building index if file already exists
func NewSSTable[T common.Serializable](filename string, blockSize int64, opts ...fnOpt) (*SSTable[T], error) {
	sstable := SSTable[T]{
		filename:  filename,
		blockSize: blockSize,
		index:     make([]indElem, 0),
	}

	for _, opt := range opts {
		opt(&sstable.modifications)
	}

	var err error
	if _, err = os.Stat(filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if sstable.file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
			return nil, err
		}
		return &sstable, err
	}

	if sstable.file, err = os.OpenFile(filename, os.O_RDWR, 0644); err != nil {
		return nil, err
	}

	if err = sstable.loadIndex(); err != nil {
		return nil, fmt.Errorf("error loading index: %v", err)
	}
	return &sstable, nil
}

func (sst *SSTable[T]) Add(key string, value T) error {
	entryData, err := serializeEntry(key, value)
	if err != nil {
		return err
	}

	entryLen := int64(len(entryData))
	if entryLen > sst.blockSize {
		return errors.New("canont add entry, object is bigger than block")
	}
	if int64(sst.blockBuf.Len())+entryLen > sst.blockSize {
		if err = sst.Flush(); err != nil {
			return err
		}
	}

	curOffset, err := sst.file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("error setting offset to the end of file: %v", curOffset)
	}
	sst.index = append(sst.index, indElem{key, curOffset})

	if err := binary.Write(&sst.blockBuf, binary.LittleEndian, entryLen); err != nil {
		return fmt.Errorf("error writing entryLen to blockBuf: %v", err)
	}
	if _, err := sst.blockBuf.Write(entryData); err != nil {
		return fmt.Errorf("error writing entryData to blockBuf: %v", err)
	}

	if sst.bloom != nil {
		sst.bloom.Add(key)
	}
	return nil
}

func (sst *SSTable[T]) Flush() error {
	if sst.blockBuf.Len() == 0 {
		return nil
	}

	// make sure we appending block to the end of file
	if _, err := sst.file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("error setting seek to end of file: %v", err)
	}
	if _, err := sst.file.Write(sst.blockBuf.Bytes()); err != nil {
		return fmt.Errorf("error flushing to file: %v", err)
	}
	// writing empty bytes to keep block offset
	margin := int(sst.blockSize - int64(sst.blockBuf.Len()))
	if _, err := sst.file.Write(make([]byte, margin)); err != nil {
		return fmt.Errorf("error writing margin to file: %v", err)
	}

	sst.blockBuf.Reset()
	return nil
}

func (sst *SSTable[T]) Close() error {
	if err := sst.Flush(); err != nil {
		return err
	}
	return sst.file.Close()
}

func (sst *SSTable[T]) Get(key string) (T, bool) {
	var defaultT T
	if sst.bloom != nil && !sst.bloom.Contains(key) {
		return defaultT, false
	}

	pos := lowerBound(sst.index, key)
	if key != sst.index[pos].key {
		if pos != 0 {
			pos -= 1
		}
	}

	// loadBlock
	entries, err := sst.loadBlock(sst.index[pos].offset)
	if err != nil {
		panic(err)
	}

	for _, ent := range entries {
		if ent.key == key {
			return ent.value, true
		}
	}

	return defaultT, false

}

func lowerBound(index []indElem, target string) int {
	left, right := 0, len(index)
	for left < right {
		mid := left + (right-left)/2
		if index[mid].key >= target {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

func upperBound(arr []string, target string) int {
	left, right := 0, len(arr)
	for left < right {
		mid := left + (right-left)/2
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}
