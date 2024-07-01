package evaluator

import "go++/object"

type arrayHelperImpl struct{}

func (h *arrayHelperImpl) ApplyFunction(fn object.Object, args []object.Object) object.Object {
	return applyFunction(fn, args)
}

func (h *arrayHelperImpl) NewInteger(value int64) *object.Integer {
	return newInteger(value)
}

func (h *arrayHelperImpl) NewError(format string, args ...interface{}) *object.Error {
	return newError(format, args...)
}

func (h *arrayHelperImpl) NewArray(values []object.Object) *object.Array {
	return newArray(values)
}

func (h *arrayHelperImpl) GetNull() *object.Null {
	return NULL
}
