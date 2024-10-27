package parser

import (
	"bufio"
	"fmt"
	"os"

	lexicalutils "github.com/M8hew/HSE-Information-Search/reversed_index/internal/lexical_utils"
)

type (
	HashSet = map[string]struct{}

	WordParser interface {
		Parse(filename string) (HashSet, error)
	}

	StemmingWordParser struct{}
)

func NewStemmingParser() WordParser {
	return &StemmingWordParser{}
}

func (swp *StemmingWordParser) Parse(filename string) (HashSet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	wordsSet := make(HashSet)
	for scanner.Scan() {
		token := lexicalutils.Regularize(scanner.Text())
		if lexicalutils.IsStopWord(token) {
			continue
		}

		term, err := lexicalutils.Stem(token)
		if err != nil {
			fmt.Printf("error while stemming token %s, %v", token, err)
			continue
		}

		wordsSet[term] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}
	return wordsSet, nil
}
