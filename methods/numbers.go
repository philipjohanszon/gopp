package methods

import "go++/object"

type NumberHelper interface {
	NewError(format string, a ...interface{}) *object.Error
	NewInteger(value int64) *object.Integer
	GetNull() *object.Null
}

func GetBuiltinNumberMethods(helper NumberHelper) map[string]object.Object {
	return map[string]object.Object{
		"add": &object.BuiltinMethod{Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return helper.NewError("ERROR: Only one argument should be given to 'add'")
			}

			switch args[0].(type) {
			case *object.Integer:
				adder, ok := args[1].(*object.Integer)

				if !ok {
					return helper.NewError("ERROR: Only one argument should be given to 'add'")
				}

				return helper.NewInteger(args[0].(*object.Integer).Value + adder.Value)
			}

			return helper.NewError("ERROR: Type %T not supported", args[0])
		}},
	}
}
