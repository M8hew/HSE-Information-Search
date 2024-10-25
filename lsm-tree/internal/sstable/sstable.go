package sstable

const blockSize = 4096

type SSTable interface {
	Get(key string) *string
	Add(key, value string)
	WriteToFile(path string)
}

type sstable struct {
	filePath string
	data     map[string]string
}

func NewSSTable(filePath string) SSTable {
	return &sstable{
		filePath: filePath,
		data:     make(map[string]string),
	}
}

func (s *sstable) Add(key, value string) {
	s.data[key] = value
}

func (s *sstable) Get(key string) *string {
	// TODO
	return nil
}

func (s *sstable) WriteToFile(path string) {
	// TODO
}
