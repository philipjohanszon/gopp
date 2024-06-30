package object

import (
	"bytes"
	"go++/ast"
	"strings"
)

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type { return RETURN }
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
func (rv *ReturnValue) GetMembers() *ObjectMembers { return nil }

type Error struct {
	Message string
}

func (e *Error) Type() Type                 { return ERROR }
func (e *Error) Inspect() string            { return "ERROR: " + e.Message }
func (e *Error) GetMembers() *ObjectMembers { return nil }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) GetMembers() *ObjectMembers { return nil }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type                 { return BUILTIN }
func (b *Builtin) Inspect() string            { return "builtin function" }
func (b *Builtin) GetMembers() *ObjectMembers { return nil }

type MethodFunction func(args ...Object) Object

type BuiltinMethod struct {
	Fn MethodFunction
	It Object
}

func (b *BuiltinMethod) Type() Type                 { return METHOD }
func (b *BuiltinMethod) Inspect() string            { return "method" }
func (b *BuiltinMethod) GetMembers() *ObjectMembers { return nil }
