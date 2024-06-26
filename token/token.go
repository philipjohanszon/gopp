package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

var Keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdentifier(identifier string) Type {
	if tok, ok := Keywords[identifier]; ok {
		return tok
	}

	return IDENTIFIER
}
