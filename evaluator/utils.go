package evaluator

import "go++/object"

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn.(*object.Function), args)
		evaluated := Evaluate(fn.(*object.Function).Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.(*object.Builtin).Fn(args...)
	case *object.BuiltinMethod:
		arguments := make([]object.Object, 1)

		arguments[0] = fn.(*object.BuiltinMethod).It
		arguments = append(arguments, args...)

		return fn.(*object.BuiltinMethod).Fn(arguments...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], true)
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func isObjectTruthy(obj object.Object) bool {
	switch obj.(type) {
	case *object.Boolean:
		return obj.(*object.Boolean).Value
	case *object.Integer:
		return obj.(*object.Integer).Value != 0
	case *object.Null:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}
