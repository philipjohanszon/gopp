package ast

import (
	"bytes"
	"go++/token"
	"strings"
)

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (bl *Boolean) expressionNode()      {}
func (bl *Boolean) TokenLiteral() string { return bl.Token.Literal }
func (bl *Boolean) String() string       { return bl.Token.Literal }

type Array struct {
	Token  token.Token
	Values []Expression
}

func (a *Array) expressionNode()      {}
func (a *Array) TokenLiteral() string { return a.Token.Literal }
func (a *Array) String() string {
	var out bytes.Buffer

	items := []string{}

	for _, value := range a.Values {
		items = append(items, value.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("fn ")

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	out.WriteString(fl.Body.String())

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

type ForLoopLiteral struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (fl *ForLoopLiteral) expressionNode()      {}
func (fl *ForLoopLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *ForLoopLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fl.Condition.String())
	out.WriteString(fl.Body.String())

	return out.String()
}
