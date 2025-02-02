package evaluator

import (
	"fmt"
	"go++/ast"
	"go++/methods"
	"go++/object"
	"strconv"
)

func newString(value string) *object.String {
	return &object.String{Value: value, Members: object.ObjectMembers{Members: methods.GetBuiltinStringMethods(&stringHelperImpl{}), MutableMembers: false}}
}

func newInteger(value int64) *object.Integer {
	return &object.Integer{Value: value, Members: object.ObjectMembers{Members: methods.GetBuiltinNumberMethods(&numberHelperImpl{}), MutableMembers: false}}
}

func newBoolean(value bool) *object.Boolean {
	return &object.Boolean{Value: value, Members: object.ObjectMembers{}}
}

func newNull() *object.Null {
	return &object.Null{}
}

func newReturnValue(v object.Object) *object.ReturnValue {
	return &object.ReturnValue{Value: v}
}

func newFunction(parameters []*ast.Identifier, body *ast.BlockStatement, env *object.Environment) *object.Function {
	return &object.Function{Parameters: parameters, Body: body, Env: env}
}

func newArray(values []object.Object) *object.Array {
	return &object.Array{
		Values:  values,
		Members: object.ObjectMembers{Members: methods.GetBuiltinArrayMethods(&arrayHelperImpl{}), MutableMembers: false},
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func nativeBoolToBooleanObject(isTrue bool) object.Object {
	if isTrue {
		return TRUE
	}

	return FALSE
}

func intToString(integer *object.Integer) *object.String {
	return &object.String{Value: strconv.Itoa(int(integer.Value))}
}
