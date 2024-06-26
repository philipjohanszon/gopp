package evaluator

import (
	"fmt"
	"go++/ast"
	"go++/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Evaluate(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)

	case *ast.LetStatement:
		value := Evaluate(node.Value, env)

		if isError(value) {
			return value
		}

		env.Set(node.Name.Value, value)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

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

	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)

	case *ast.IfExpression:
		return evaluateIfExpression(node, env)

	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, env)

		if isError(value) {
			return value
		}

		return &object.ReturnValue{Value: value}
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
	function, ok := fn.(*object.Function)

	if !ok {
		return newError("%s is not a function", function.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Evaluate(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}
