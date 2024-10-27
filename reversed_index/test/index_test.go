package test

import (
	"fmt"
	"testing"

	"github.com/M8hew/HSE-Information-Search/reversed_index/internal/index"
	"github.com/stretchr/testify/require"
)

type stringSet = map[string]struct{}

var indexer index.Indexer

func TestMain(m *testing.M) {
	indexer = index.NewBitmapIndexer()
	err := indexer.IndexFiles("../dataset")
	if err != nil {
		fmt.Println(err)
		return
	}

	m.Run()
}

func getFileSet(files []string) stringSet {
	fileSet := make(stringSet)
	for _, file := range files {
		fileSet[file] = struct{}{}
	}
	return fileSet
}

func TestFind(t *testing.T) {
	testCases := []struct {
		name   string
		target string
		err    error
		result stringSet
	}{
		{
			name:   "simple find success",
			target: "voice",
			err:    nil,
			result: stringSet{"../dataset/focus": {}},
		},
		{
			name:   "error stop word",
			target: "an",
			err:    index.ErrStopWord,
			result: nil,
		},
		{
			name:   "error not found",
			target: "abboba",
			err:    index.ErrNotFound,
			result: nil,
		},
		{
			name:   "multiple files search result 1",
			target: "move",
			err:    nil,
			result: stringSet{
				"../dataset/saying_no":        {},
				"../dataset/empire_state":     {},
				"../dataset/shagun_chowdhary": {},
				"../dataset/innovative_ideas": {},
			},
		},
		{
			name:   "multiple files search result 2",
			target: "lot",
			err:    nil,
			result: stringSet{
				"../dataset/focus":            {},
				"../dataset/saying_no":        {},
				"../dataset/creative_brain":   {},
				"../dataset/innovative_ideas": {},
			},
		},
	}

	for _, tc := range testCases {
		vec, err := indexer.Find(tc.target)
		if tc.err != nil {
			require.Equal(t, tc.err, err)
			continue
		}

		fileSet := getFileSet(indexer.GetFileNames(vec))
		require.Equal(t, tc.result, fileSet)
	}
}

func Test_VectorOperations(t *testing.T) {
	vec1, err := indexer.Find("move")
	require.NoError(t, err)

	vec2, err := indexer.Find("lot")
	require.NoError(t, err)

	vec1.AndNot(vec2)
	filenames := getFileSet(indexer.GetFileNames(vec1))

	require.Equal(t,
		stringSet{
			"../dataset/empire_state":     {},
			"../dataset/shagun_chowdhary": {},
		},
		filenames)

}
