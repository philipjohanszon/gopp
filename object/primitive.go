package object

import (
	"fmt"
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
