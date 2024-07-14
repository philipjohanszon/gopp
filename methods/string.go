package methods

import (
	"go++/object"
	"strings"
)

type StringHelper interface {
	NewError(format string, a ...interface{}) *object.Error
	NewInteger(value int64) *object.Integer
	NewString(value string) *object.String
}

func GetBuiltinStringMethods(helper StringHelper) map[string]object.Object {
	return map[string]object.Object{
		"length": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return helper.NewString("ERROR: No arguments should be given to 'length'")
			}

			if _, ok := args[0].(*object.String); !ok {
				return helper.NewError("ERROR: First argument must be a string")
			}

			return helper.NewInteger(int64(len(args[0].(*object.String).Value)))
		}},
		"replace": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return helper.NewError("ERROR: Replace requires 2 arguments")
			}

			it, ok := args[0].(*object.String)

			if !ok {
				return helper.NewError("ERROR: First argument must be a string")
			}

			replace, ok := args[1].(*object.String)

			if !ok {
				return helper.NewError("ERROR: First argument must be a string")
			}

			replacer, ok := args[2].(*object.String)

			if !ok {
				return helper.NewError("ERROR: Second argument must be a string")
			}

			it.Value = strings.Replace(it.Value, replace.Value, replacer.Value, -1)

			return args[0]
		}},
	}
}
