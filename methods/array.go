package methods

import (
	"go++/object"
)

type ArrayHelper interface {
	ApplyFunction(fn object.Object, args []object.Object) object.Object
	NewError(format string, a ...interface{}) *object.Error
	NewInteger(value int64) *object.Integer
	NewArray(values []object.Object) *object.Array
	GetNull() *object.Null
}

func GetBuiltinArrayMethods(helper ArrayHelper) map[string]object.Object {
	return map[string]object.Object{
		"length": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return helper.NewError("ERROR: No arguments should be given to 'length'")
			}

			if _, ok := args[0].(*object.Array); !ok {
				return helper.NewError("ERROR: First argument must be an array")
			}

			return &object.Integer{Value: int64(len(args[0].(*object.Array).Values))}
		}},
		"forEach": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return helper.NewError("ERROR: Only a callback should be given to 'forEach'")
			}

			if _, ok := args[0].(*object.Array); !ok {
				return helper.NewError("ERROR: First argument must be an array")
			}

			if _, ok := args[1].(*object.Function); !ok {
				return helper.NewError("ERROR: First argument must be a function")
			}

			for key, value := range args[0].(*object.Array).Values {
				helper.ApplyFunction(args[1].(*object.Function), []object.Object{helper.NewInteger(int64(key)), value})
			}

			return helper.GetNull()
		}},
		"map": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return helper.NewError("ERROR: Only a callback should be given to 'map'")
			}

			if _, ok := args[0].(*object.Array); !ok {
				return helper.NewError("ERROR: First argument must be an array")
			}

			if _, ok := args[1].(*object.Function); !ok {
				return helper.NewError("ERROR: First argument must be a function")
			}

			var newArray []object.Object

			for key, value := range args[0].(*object.Array).Values {
				newArray = append(newArray, helper.ApplyFunction(args[1].(*object.Function), []object.Object{helper.NewInteger(int64(key)), value}))
			}

			return helper.NewArray(newArray)
		}},
	}
}
