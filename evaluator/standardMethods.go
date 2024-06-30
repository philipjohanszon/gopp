package evaluator

import (
	"go++/object"
	"strings"
)

var BuiltinStringMethods = map[string]object.Object{
	"length": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("ERROR: No arguments should be given to 'length'")
		}

		if _, ok := args[0].(*object.String); !ok {
			return newError("ERROR: First argument must be a string")
		}

		return &object.Integer{Value: int64(len(args[0].(*object.String).Value))}
	}},
	"replace": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
		if len(args) != 3 {
			return newError("ERROR: Replace requires 2 arguments")
		}

		it, ok := args[0].(*object.String)

		if !ok {
			return newError("ERROR: First argument must be a string")
		}

		replace, ok := args[1].(*object.String)

		if !ok {
			return newError("ERROR: First argument must be a string")
		}

		replacer, ok := args[2].(*object.String)

		if !ok {
			return newError("ERROR: Second argument must be a string")
		}

		it.Value = strings.Replace(it.Value, replace.Value, replacer.Value, -1)

		return args[0]
	}},
}

var BuiltinNumberMethods = map[string]object.Object{
	"add": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("ERROR: Only one argument should be given to 'add'")
		}

		switch args[0].(type) {
		case *object.Integer:
			adder, ok := args[1].(*object.Integer)

			if !ok {
				return newError("ERROR: Only one argument should be given to 'add'")
			}

			args[0].(*object.Integer).Value += adder.Value
		}

		return args[0]
	}},
}
