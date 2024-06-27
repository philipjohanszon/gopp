package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return INTEGER }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BOOLEAN }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() Type      { return NULL }
func (n *Null) Inspect() string { return "null" }

type String struct {
	Value string
}

func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return s.Value }
