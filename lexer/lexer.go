package lexer

import (
	"bytes"
	"go++/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	currentChar  byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readCharacter()

	return lexer
}

func (lexer *Lexer) readCharacter() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.currentChar {
	case '=':
		if lexer.peekChar() == '=' {
			character := lexer.currentChar
			lexer.readCharacter()

			tok.Literal = string(character) + string(lexer.currentChar)
			tok.Type = token.EQUALS
		} else {
			tok = newToken(token.ASSIGN, lexer.currentChar)
		}
	case '+':
		tok = newToken(token.PLUS, lexer.currentChar)
	case '-':
		tok = newToken(token.MINUS, lexer.currentChar)
	case '*':
		tok = newToken(token.ASTERISK, lexer.currentChar)
	case '/':
		tok = newToken(token.SLASH, lexer.currentChar)
	case '<':
		tok = newToken(token.LESSTHAN, lexer.currentChar)
	case '>':
		tok = newToken(token.GREATERTHAN, lexer.currentChar)
	case '!':
		if lexer.peekChar() == '=' {
			character := lexer.currentChar
			lexer.readCharacter()

			tok.Literal = string(character) + string(lexer.currentChar)
			tok.Type = token.NOTEQUALS
		} else {
			tok = newToken(token.NOT, lexer.currentChar)
		}
	case ';':
		tok = newToken(token.SEMICOLON, lexer.currentChar)
	case ',':
		tok = newToken(token.COMMA, lexer.currentChar)
	case '.':
		tok = newToken(token.DOT, lexer.currentChar)
	case '(':
		tok = newToken(token.LPAREN, lexer.currentChar)
	case ')':
		tok = newToken(token.RPAREN, lexer.currentChar)
	case '{':
		tok = newToken(token.LBRACE, lexer.currentChar)
	case '}':
		tok = newToken(token.RBRACE, lexer.currentChar)
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()

	case 0:
		tok = newToken(token.EOF, lexer.currentChar)
		tok.Literal = ""
	default:
		if isLetter(lexer.currentChar) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)

			return tok
		} else if IsDigit(lexer.currentChar) {
			tok.Literal = lexer.readNumber()
			tok.Type = token.INTEGER

			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.currentChar)
		}
	}

	lexer.readCharacter()
	return tok
}

func newToken(tokenType token.Type, character byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	}

	return lexer.input[lexer.readPosition]
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.currentChar) {
		lexer.readCharacter()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readString() string {
	//position := lexer.position + 1
	var out bytes.Buffer

	for {
		lexer.readCharacter()

		if lexer.currentChar == '\\' {
			switch lexer.peekChar() {
			case 'n':
				out.WriteByte('\n')
				lexer.readCharacter() // Skip the 'n'
			case '\\':
				out.WriteByte('\\')
				lexer.readCharacter() // Skip the second '\'
			case '"':
				out.WriteByte('"')
				lexer.readCharacter() // Skip the second '"'
			default:
				out.WriteByte(lexer.currentChar)
			}
			continue
		}

		if lexer.currentChar == '"' || lexer.currentChar == '0' {
			break
		}

		out.WriteByte(lexer.currentChar)
	}

	return out.String()
}

func isLetter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\r' || lexer.currentChar == '\n' {
		lexer.readCharacter()
	}
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position

	for IsDigit(lexer.currentChar) {
		lexer.readCharacter()
	}

	return lexer.input[position:lexer.position]
}

func IsDigit(character byte) bool {
	return '0' <= character && character <= '9'
}
