package ast

import (
	"bytes"
	"go++/token"
	"strings"
)

type Expression interface {
	Node
	expressionNode()
	String() string
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string

	Left  Expression
	Right Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type AssignExpression struct {
	Token    token.Token
	Assignee Expression
	Value    Expression
}

func (a *AssignExpression) expressionNode()      {}
func (a *AssignExpression) TokenLiteral() string { return a.Token.Literal }
func (a *AssignExpression) String() string {
	return a.Assignee.String() + " = " + a.Value.String()
}

type MemberAccessExpression struct {
	Token          token.Token
	Expression     Expression
	AccessedMember Identifier
}

func (ma *MemberAccessExpression) expressionNode()      {}
func (ma *MemberAccessExpression) TokenLiteral() string { return ma.Token.Literal }
func (ma *MemberAccessExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ma.Expression.String())
	out.WriteString(".")
	out.WriteString(ma.AccessedMember.Value)
	out.WriteString(")")

	return out.String()
}

type ArrayAccessExpression struct {
	Token      token.Token
	Expression Expression
	Index      Expression
}

func (aa *ArrayAccessExpression) expressionNode()      {}
func (aa *ArrayAccessExpression) TokenLiteral() string { return aa.Token.Literal }
func (aa *ArrayAccessExpression) String() string {
	var out bytes.Buffer

	out.WriteString(aa.Expression.String() + "[" + aa.Index.String() + "]")

	return out.String()
}
