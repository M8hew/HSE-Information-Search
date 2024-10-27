package index

import (
	"errors"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"os"

	lexicalutils "github.com/M8hew/HSE-Information-Search/reversed_index/internal/lexical_utils"
	"github.com/M8hew/HSE-Information-Search/reversed_index/internal/parser"

	"github.com/RoaringBitmap/roaring/v2/roaring64"
)

var (
	ErrNotFound = errors.New("word not found in index")
	ErrStopWord = errors.New("stop word was given, stop words not indexing")
)

type (
	fileName      = string
	fileID        = uint64
	RoaringBitmap = *roaring64.Bitmap

	Indexer interface {
		IndexFiles(dirPath string) error
		Find(word string) (BitVector, error)
		GetFileNames(vec BitVector) []string
	}

	bitmapIndexer struct {
		fileNameDict map[fileID]fileName
		hasher       hash.Hash64
		index        map[string]RoaringBitmap
		parser       parser.WordParser
	}
)

func NewBitmapIndexer() Indexer {
	return &bitmapIndexer{
		fileNameDict: make(map[fileID]fileName),
		hasher:       crc64.New(crc64.MakeTable(crc64.ECMA)),
		index:        make(map[fileName]RoaringBitmap),
		parser:       parser.NewStemmingParser(),
	}
}

func (bi *bitmapIndexer) IndexFiles(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Printf("error geting fileinfo %s: %v", file.Name(), err)
			continue
		}

		fileID := bi.getFileID(fileInfo)
		fileName := dirPath + "/" + file.Name()
		bi.fileNameDict[fileID] = fileName

		words, err := bi.parser.Parse(fileName)
		if err != nil {
			fmt.Printf("error parsing file %s: %v", file.Name(), err)
			continue
		}

		for word := range words {
			if bitmap, ok := bi.index[word]; ok {
				bitmap.Add(fileID)
				continue
			}

			bi.index[word] = roaring64.BitmapOf(fileID)
		}
	}
	return nil
}

func (bi *bitmapIndexer) Find(word string) (BitVector, error) {
	token := lexicalutils.Regularize(word)
	if lexicalutils.IsStopWord(token) {
		return BitVector{}, ErrStopWord
	}

	term, err := lexicalutils.Stem(token)
	if err != nil {
		return BitVector{}, fmt.Errorf("error while stemming given word: %v", err)
	}

	vec, ok := bi.index[term]
	if !ok {
		return BitVector{}, ErrNotFound
	}
	return BitVector{vec}, nil
}

func (bi *bitmapIndexer) GetFileNames(vec BitVector) []string {
	var fileNames []string

	it := vec.RoaringBitmap.Iterator()
	for it.HasNext() {
		curFileID := it.Next()
		fileNames = append(fileNames, bi.fileNameDict[curFileID])
	}
	return fileNames
}

func (bi *bitmapIndexer) getFileID(fileInfo os.FileInfo) uint64 {
	io.WriteString(bi.hasher, fileInfo.Name())
	hash := bi.hasher.Sum64()
	bi.hasher.Reset()
	return hash
}
