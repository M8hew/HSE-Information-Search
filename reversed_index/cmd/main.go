package main

import (
	"fmt"

	"github.com/M8hew/HSE-Information-Search/reversed_index/internal/index"
)

func main() {
	indexer := index.NewBitmapIndexer()
	err := indexer.IndexFiles("./dataset")
	if err != nil {
		fmt.Println(err)
		return
	}

	vec, _ := indexer.Find("lot")
	fmt.Printf("%#v\n", indexer.GetFileNames(vec))
}
