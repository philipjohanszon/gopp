package parser

import (
	"fmt"
	"go++/ast"
	"go++/token"
	"strconv"
)

func (parser *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (parser *Parser) parseBoolean() ast.Expression {
	literal := &ast.Boolean{Token: parser.currentToken, Value: parser.currentTokenIs(token.TRUE)}

	return literal
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
		Left:     left,
	}

	precedence := parser.currentPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	exp := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.currentToken}

	parser.nextToken()

	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()

		if !parser.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: parser.currentToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = parser.parseBlockStatement()

	return literal
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.currentToken, Function: function}
	expression.Arguments = parser.parseCallArguments(token.RPAREN)

	return expression
}

// expr <dot> expr() <dot> expr()
func (parser *Parser) parseMemberAccessExpression(expr ast.Expression) ast.Expression {
	expression := &ast.MemberAccessExpression{Token: parser.currentToken, Expression: expr}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	expression.AccessedMember = ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	return expression
}

func (parser *Parser) parseCallArguments(endingToken token.Type) []ast.Expression {
	args := []ast.Expression{}

	if parser.peekTokenIs(endingToken) {
		parser.nextToken()
		return args
	}

	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(endingToken) {
		return nil
	}

	return args
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseAssignExpression(exp ast.Expression) ast.Expression {
	expression := &ast.AssignExpression{Token: parser.currentToken}

	identifier, ok := exp.(*ast.Identifier)

	if !ok {
		parser.appendError("Was expecting an identifier, got %T", exp)
		return nil
	}

	expression.Assignee = *identifier

	parser.nextToken()

	expression.Value = parser.parseExpression(LOWEST)

	return expression
}

func (parser *Parser) parseArray() ast.Expression {
	arr := &ast.Array{Token: parser.currentToken}
	arr.Values = parser.parseCallArguments(token.RBRACKET)

	return arr
}

func (parser *Parser) parseArrayAccess(array ast.Expression) ast.Expression {
	expr := &ast.ArrayAccessExpression{
		Token:      parser.currentToken,
		Expression: array,
		Index:      nil,
	}

	parser.nextToken()

	expr.Index = parser.parseExpression(LOWEST)

	parser.nextToken()

	return expr
}
