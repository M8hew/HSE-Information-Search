package sstable

import (
	"bytes"
	"encoding/binary"
	"io"

	"lsm_tree/internal/common"
)

type entry[T common.Serializable] struct {
	key   string
	value T
}

func serializeEntry[T common.Serializable](key string, value T) ([]byte, error) {
	valueBytes, err := value.Marshal()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	keyLen := uint32(len([]byte(key)))
	if err := binary.Write(&buf, binary.LittleEndian, keyLen); err != nil {
		return nil, err
	}

	if _, err := buf.Write([]byte(key)); err != nil {
		return nil, err
	}

	valueLen := uint32(len(valueBytes))
	if err := binary.Write(&buf, binary.LittleEndian, valueLen); err != nil {
		return nil, err
	}

	if _, err := buf.Write(valueBytes); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func loadEntry[T common.Serializable](buffer io.Reader) (*entry[T], int, error) {
	var totalBytesRead int

	var keyLen uint32
	if err := binary.Read(buffer, binary.LittleEndian, &keyLen); err != nil {
		return nil, totalBytesRead, err
	}
	totalBytesRead += 4

	key := make([]byte, keyLen)
	n, err := io.ReadFull(buffer, key)
	if err != nil {
		return nil, totalBytesRead, err
	}
	totalBytesRead += n

	var valueLen uint32
	if err = binary.Read(buffer, binary.LittleEndian, &valueLen); err != nil {
		return nil, totalBytesRead, err
	}
	totalBytesRead += 4

	valueBytes := make([]byte, valueLen)
	n, err = io.ReadFull(buffer, valueBytes)
	if err != nil {
		return nil, totalBytesRead, err
	}
	totalBytesRead += n

	var tmp T
	value := tmp.New()
	if err = value.Unmarshal(valueBytes); err != nil {
		return nil, totalBytesRead, err
	}

	return &entry[T]{key: string(key), value: value.(T)}, totalBytesRead, nil
}

// func loadEntry[T common.Serializable](buffer io.Reader) (*entry[T], int64, error) {
// 	var (
// 		size     int64
// 		ent      entry[T]
// 		entBytes []byte
// 	)

// 	if err := binary.Read(buffer, binary.LittleEndian, &size); err != nil {
// 		return nil, 0, fmt.Errorf("erorr reading binary size: %v", err)
// 	}
// 	if size == 0 {
// 		return nil, 0, io.EOF
// 	}

// 	if _, err := buffer.Read(entBytes); err != nil {
// 		return nil, 0, fmt.Errorf("error reading binary entry: %v", err)
// 	}
// 	if err := binCodec.Unmarshal(entBytes, &ent); err != nil {
// 		return nil, 0, fmt.Errorf("error unmarshaling binary entry: %v", err)
// 	}
// 	return &ent, int64(size) + 8, nil
// }
