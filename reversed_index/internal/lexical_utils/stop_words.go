package lexicalutils

type HashSet = map[string]struct{}

func IsStopWord(word string) bool {
	_, ok := stopWords[word]
	return ok
}

var stopWords = HashSet{
	"a":     {},
	"also":  {},
	"as":    {},
	"an":    {},
	"at":    {},
	"am":    {},
	"and":   {},
	"are":   {},
	"be":    {},
	"but":   {},
	"by":    {},
	"can":   {},
	"could": {},
	"did":   {},
	"do":    {},
	"does":  {},
	"else":  {},
	"for":   {},
	"from":  {},
	"had":   {},
	"has":   {},
	"have":  {},
	"i":     {},
	"if":    {},
	"in":    {},
	"is":    {},
	"it":    {},
	"its":   {},
	"may":   {},
	"maybe": {},
	"me":    {},
	"mine":  {},
	"must":  {},
	"my":    {},
	"nor":   {},
	"not":   {},
	"of":    {},
	"oh":    {},
	"on":    {},
	"or":    {},
	"that":  {},
	"the":   {},
	"to":    {},
	"was":   {},
	"were":  {},
	"with":  {},
	"you":   {},
	"your":  {},
}
