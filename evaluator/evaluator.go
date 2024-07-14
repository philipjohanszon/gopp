package evaluator

import (
	"go++/ast"
	"go++/object"
)

var (
	NULL  = newNull()
	TRUE  = newBoolean(true)
	FALSE = newBoolean(false)
)

func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)

	case *ast.IntegerLiteral:
		return newInteger(node.Value)
	case *ast.StringLiteral:
		return newString(node.Value)
	case *ast.Array:
		elements := make([]object.Object, len(node.Values))

		for k, v := range node.Values {
			elements[k] = Evaluate(v, env)
		}

		return newArray(elements)
	case *ast.AssignExpression:
		evaluated := Evaluate(node.Value, env)

		if isError(evaluated) {
			return evaluated
		}

		if identifier, ok := node.Assignee.(*ast.Identifier); ok {
			obj, done := assignIdentifier(identifier, evaluated, env)

			if done {
				return obj
			}
		}

		if arrayAccess, ok := node.Assignee.(*ast.ArrayAccessExpression); ok {
			obj, done := assignArray(arrayAccess, evaluated, env)

			if done {
				return obj
			}
		}

		return NULL

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)

	case *ast.LetStatement:
		value := Evaluate(node.Value, env)

		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value, node.IsMutable)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return newFunction(params, body, env)

	case *ast.ForLoopLiteral:
		outerEnv := object.NewEnclosedEnvironment(env)

		for isObjectTruthy(Evaluate(node.Condition, env)) {
			evaluated := evaluateBlockStatement(node.Body, outerEnv)

			if isError(evaluated) {
				return evaluated
			}
		}

	case *ast.CallExpression:
		function := Evaluate(node.Function, env)

		if isError(function) {
			return function
		}

		args := evaluateExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.PrefixExpression:
		right := Evaluate(node.Right, env)

		if isError(right) {
			return right
		}

		return evaluatePrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Evaluate(node.Left, env)

		if isError(left) {
			return left
		}

		right := Evaluate(node.Right, env)

		if isError(right) {
			return right
		}

		return evaluateInfixExpression(node.Operator, left, right)

	case *ast.MemberAccessExpression:
		left := Evaluate(node.Expression, env)

		if isError(left) {
			return left
		}

		val, ok := left.GetMembers().Get(node.AccessedMember.Value)

		if !ok {
			return newError("Error: %s is not member of %s", node.AccessedMember.Value, left.Inspect())
		}

		if isError(val) {
			return val
		}

		if method, ok := val.(*object.BuiltinMethod); ok {
			method.It = left

			return method
		}

		return val

	case *ast.ArrayAccessExpression:
		index := Evaluate(node.Index, env)
		array := Evaluate(node.Expression, env)

		if isError(array) {
			return array
		}

		if isError(index) {
			return index
		}

		if _, ok := array.(*object.Array); !ok {
			return newError("ERROR: %s is not an array", array.Inspect())
		}

		if _, ok := index.(*object.Integer); !ok {
			return newError("ERROR: index: %s is not an integer", index.Inspect())
		}

		return array.(*object.Array).GetIndex(int(index.(*object.Integer).Value))

	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)

	case *ast.IfExpression:
		return evaluateIfExpression(node, env)

	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, env)

		if isError(value) {
			return value
		}

		return newReturnValue(value)
	}

	return nil
}

func evaluateProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evaluateBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evaluateExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Evaluate(expression, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

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

func nativeBoolToBooleanObject(isTrue bool) object.Object {
	if isTrue {
		return TRUE
	}

	return FALSE
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
