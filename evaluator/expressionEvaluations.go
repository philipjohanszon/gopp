package evaluator

import (
	"go++/ast"
	"go++/object"
	"strconv"
)

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator %s%s", operator, right.Type())
	}
}

func evaluateBangOperatorExpression(right object.Object) object.Object {
	if isObjectTruthy(right) {
		return FALSE
	}

	return TRUE
}

func evaluateMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evaluateIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evaluateStringInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.INTEGER:
		return evaluateStringInfixExpression(operator, left, intToString(right.(*object.Integer)))
	case left.Type() == object.INTEGER && right.Type() == object.STRING:
		return evaluateStringInfixExpression(operator, intToString(right.(*object.Integer)), right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func intToString(integer *object.Integer) *object.String {
	return &object.String{Value: strconv.Itoa(int(integer.Value))}
}

func evaluateIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Evaluate(node.Condition, env)

	if isError(condition) {
		return condition
	}

	outerEnv := object.NewEnclosedEnvironment(env)

	if isObjectTruthy(condition) {
		return Evaluate(node.Consequence, outerEnv)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative, outerEnv)
	} else {
		return NULL
	}
}

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	value, ok := env.Get(node.Value)

	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return value
}
