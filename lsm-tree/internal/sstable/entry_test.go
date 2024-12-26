package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lsm_tree/internal/common"
)

type SimpleValue struct {
	Value int64
}

func (s SimpleValue) Marshal() ([]byte, error) {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, uint64(s.Value))
	return buffer, nil
}

func (s *SimpleValue) Unmarshal(data []byte) error {
	if s == nil {
		s = &SimpleValue{}
	}
	if len(data) < 8 {
		return errors.New("not enough data to unmarshal SimpleValue")
	}
	s.Value = int64(binary.LittleEndian.Uint64(data))
	return nil
}

func (s *SimpleValue) New() common.Serializable {
	return &SimpleValue{}
}

func TestLoadEntryMyString(t *testing.T) {
	key := "simple"
	value := common.ToMyString("value")

	data, err := serializeEntry(key, value)
	require.NoError(t, err)

	buf := bytes.NewReader(data)
	loadedEntry, size, err := loadEntry[*common.MyString](buf)
	require.NoError(t, err)

	assert.Equal(t, len(data), size)
	assert.Equal(t, key, loadedEntry.key)
	assert.Equal(t, value, loadedEntry.value)
	fmt.Println(value)
	fmt.Println(loadedEntry.value)

}

func TestLoadEntrySimpleValue(t *testing.T) {
	key := "simple"
	value := SimpleValue{Value: 42}

	data, err := serializeEntry(key, &value)
	require.NoError(t, err)

	buf := bytes.NewReader(data)
	loadedEntry, size, err := loadEntry[*SimpleValue](buf)
	require.NoError(t, err)

	assert.Equal(t, len(data), size)
	assert.Equal(t, key, loadedEntry.key)
	assert.Equal(t, &value, loadedEntry.value)
	fmt.Println(value)
	fmt.Println(loadedEntry.value)

}
