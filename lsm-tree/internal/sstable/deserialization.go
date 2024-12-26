package sstable

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func (sst *SSTable[T]) loadIndex() error {
	info, err := sst.file.Stat()
	if err != nil {
		return err
	}

	block := make([]byte, sst.blockSize)
	for offset := int64(0); offset < info.Size(); offset += sst.blockSize {
		if _, err := sst.file.ReadAt(block, offset); err != nil {
			return fmt.Errorf("error reading byteBlock from file: %v", err)
		}

		var ent *entry[T]
		blockReader := bytes.NewReader(block)
		if ent, _, err = loadEntry[T](blockReader); err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("error loading entry from block begining: %v", err)
			}
			return nil
		}
		sst.index = append(sst.index, indElem{ent.key, offset})
	}
	return nil
}

func (sst *SSTable[T]) loadBlock(offset int64) ([]entry[T], error) {
	block := make([]byte, sst.blockSize)
	if _, err := sst.file.ReadAt(block, offset); err != nil {
		return nil, fmt.Errorf("error reading block from file: %v", err)
	}

	var (
		blockOffset int
		entries     []entry[T]
	)
	for int64(blockOffset) < sst.blockSize {
		blockReader := bytes.NewReader(block[offset:])
		ent, entLen, err := loadEntry[T](blockReader)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("error loading entry from block: %v", err)
			}
			break
		}
		entries = append(entries, *ent)
		blockOffset += entLen
	}
	return entries, nil
}
