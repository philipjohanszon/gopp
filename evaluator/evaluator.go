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

	// ------- LITERALS -------

	case *ast.IntegerLiteral:
		return newInteger(node.Value)
	case *ast.StringLiteral:
		return newString(node.Value)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		elements := make([]object.Object, len(node.Values))

		for k, v := range node.Values {
			elements[k] = Evaluate(v, env)
		}

		return newArray(elements)

	case *ast.Identifier:
		return evaluateIdentifier(node, env)

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

	// ------- EXPRESSIONS -------

	case *ast.AssignExpression:
		return evaluateAssignExpression(node, env)

	case *ast.CallExpression:
		return evaluateCallExpression(node, env)

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
		return evaluateMemberAccessExpression(node, env)

	case *ast.ArrayAccessExpression:
		return evaluateArrayAccessExpression(node, env)

	case *ast.IfExpression:
		return evaluateIfExpression(node, env)

	// ------ STATEMENTS ------

	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, env)

		if isError(value) {
			return value
		}

		return newReturnValue(value)

	case *ast.BlockStatement:
		return evaluateBlockStatement(node, env)

	case *ast.LetStatement:
		return evaluateLetStatement(node, env)
	}

	return NULL
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
