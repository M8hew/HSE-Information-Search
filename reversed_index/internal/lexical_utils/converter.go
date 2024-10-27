package lexicalutils

import (
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
)

var re = regexp.MustCompile(`[^a-z0-9\-]+`)

func Regularize(word string) string {
	return re.ReplaceAllString(strings.ToLower(word), "")
}

func Stem(word string) (string, error) {
	return snowball.Stem(word, "english", true)
}
