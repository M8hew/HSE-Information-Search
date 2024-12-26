package common

import (
	"bytes"
	"encoding/binary"
)

type (
	MyInt    int
	MyString string
)

// static checks
var (
	_ Serializable = new(MyInt)
	_ Serializable = new(MyString)
)

func (mi MyInt) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, int64(mi))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mi *MyInt) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, mi)
}

func (mi *MyInt) New() Serializable {
	return ToMyInt(0)
}

func ToMyInt(num int) *MyInt {
	myNum := MyInt(num)
	return &myNum
}

func (ms MyString) Marshal() ([]byte, error) {
	return []byte(ms), nil
}

func (ms *MyString) Unmarshal(data []byte) error {
	*ms = MyString(string(data))
	return nil
}

func (mi *MyString) New() Serializable {
	return ToMyString("")
}

func ToMyString(str string) *MyString {
	myString := MyString(str)
	return &myString
}
