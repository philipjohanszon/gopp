package parser

import (
	"fmt"
	"go++/token"
)

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) currentTokenIs(t token.Type) bool {
	return parser.currentToken.Type == t
}

func (parser *Parser) peekTokenIs(t token.Type) bool {
	return parser.peekToken.Type == t
}

func (parser *Parser) expectPeek(t token.Type) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	}

	parser.appendError("expected next token to be of type %s, got type %s instead", t, parser.peekToken.Type)
	return false
}

func (parser *Parser) peekPrecedence() int {
	if p, ok := precedences[parser.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (parser *Parser) currentPrecedence() int {
	if p, ok := precedences[parser.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (parser *Parser) appendError(format string, a ...interface{}) {
	parser.errors = append(parser.errors, fmt.Sprintf(format, a...))
}
