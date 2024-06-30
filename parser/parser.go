package parser

import (
	"fmt"
	"go++/ast"
	lex "go++/lexer"
	"go++/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	MEMBERACCESS
)

var precedences = map[token.Type]int{
	token.ASSIGN:      EQUALS,
	token.EQUALS:      EQUALS,
	token.NOTEQUALS:   EQUALS,
	token.LESSTHAN:    LESSGREATER,
	token.GREATERTHAN: LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.LPAREN:      CALL,
	token.DOT:         MEMBERACCESS,
}

type Parser struct {
	lexer *lex.Lexer

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn

	errors []string
}

func New(lexer *lex.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFns = make(map[token.Type]prefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INTEGER, parser.parseIntegerLiteral)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.NOT, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)

	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.registerPrefix(token.STRING, parser.parseStringLiteral)

	parser.infixParseFns = make(map[token.Type]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.NOTEQUALS, parser.parseInfixExpression)
	parser.registerInfix(token.LESSTHAN, parser.parseInfixExpression)
	parser.registerInfix(token.GREATERTHAN, parser.parseInfixExpression)

	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.DOT, parser.parseMemberAccessExpression)
	parser.registerInfix(token.ASSIGN, parser.parseAssignExpression)

	return parser
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		stmt := parser.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.FOR:
		return parser.parseForLoopLiteral()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.currentToken}

	if parser.peekTokenIs(token.MUT) {
		stmt.IsMutable = true
		parser.nextToken()
	}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	stmt.Value = parser.parseExpression(LOWEST)

	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	stmt.ReturnValue = parser.parseExpression(LOWEST)

	return stmt
}

func (parser *Parser) parseForLoopLiteral() *ast.ExpressionStatement {
	stmt := &ast.ForLoopLiteral{Token: parser.currentToken}

	parser.nextToken()

	stmt.Condition = parser.parseExpression(LOWEST)

	parser.nextToken()

	stmt.Body = parser.parseBlockStatement()

	return &ast.ExpressionStatement{Token: stmt.Token, Expression: stmt}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.currentToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
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

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.currentToken.Type]

	if prefix == nil {
		parser.noPrefixParseFnError(parser.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		parser.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (parser *Parser) noPrefixParseFnError(t token.Type) {
	parser.appendError("no prefix parse function for %s found", t)
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

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.currentToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.currentTokenIs(token.RBRACE) && !parser.currentTokenIs(token.EOF) {
		stmt := parser.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		parser.nextToken()
	}

	return block
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()

	ident := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		ident := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) appendError(format string, a ...interface{}) {
	parser.errors = append(parser.errors, fmt.Sprintf(format, a...))
}
