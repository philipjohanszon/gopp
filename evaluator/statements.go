package evaluator

import (
	"go++/ast"
	"go++/object"
)

func evaluateLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	value := Evaluate(node.Value, env)

	if isError(value) {
		return value
	}

	env.Set(node.Name.Value, value, node.IsMutable)

	return NULL
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
