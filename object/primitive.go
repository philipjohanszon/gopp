package object

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Integer struct {
	Value   int64
	Members ObjectMembers
}

func (i *Integer) Type() Type                 { return INTEGER }
func (i *Integer) Inspect() string            { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) GetMembers() *ObjectMembers { return &i.Members }

type Boolean struct {
	Value   bool
	Members ObjectMembers
}

func (b *Boolean) Type() Type                 { return BOOLEAN }
func (b *Boolean) Inspect() string            { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) GetMembers() *ObjectMembers { return &b.Members }

type Null struct{}

func (n *Null) Type() Type                 { return NULL }
func (n *Null) Inspect() string            { return "null" }
func (n *Null) GetMembers() *ObjectMembers { return nil }

type String struct {
	Value   string
	Members ObjectMembers
}

func (s *String) Type() Type                 { return STRING }
func (s *String) Inspect() string            { return s.Value }
func (s *String) GetMembers() *ObjectMembers { return &s.Members }

type Array struct {
	Values  []Object
	Members ObjectMembers
}

func (a *Array) Type() Type { return ARRAY }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	items := []string{}

	for _, value := range a.Values {
		items = append(items, value.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")

	return out.String()
}
func (a *Array) GetMembers() *ObjectMembers { return &a.Members }
func (a *Array) GetIndex(i int) Object {
	if i >= len(a.Values) {
		return &Error{Message: "ERROR: index " + strconv.Itoa(i) + " out of range"}
	}

	return a.Values[i]
}
