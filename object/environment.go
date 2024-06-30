package object

type EnvironmentObject struct {
	IsMutable bool
	Object    Object
}

type Environment struct {
	store map[string]*EnvironmentObject
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*EnvironmentObject)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(environment *Environment) *Environment {
	env := NewEnvironment()
	env.outer = environment
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	envObj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, outerOk := e.outer.Get(name)

		return obj, outerOk
	}

	if !ok {
		return nil, false
	}

	return envObj.Object, ok
}

func (e *Environment) Set(name string, obj Object, isMutable bool) Object {
	e.store[name] = &EnvironmentObject{isMutable, obj}
	return obj
}

func (e *Environment) ReAssign(name string, obj Object) (Object, bool) {
	envObj, ok := e.store[name]

	if ok {
		if !envObj.IsMutable {
			return &Error{Message: "ERROR: Can't reassign immutable object: " + name}, false
		}

		e.store[name] = &EnvironmentObject{true, obj}
		return obj, true
	}

	if e.outer != nil {
		return e.outer.ReAssign(name, obj)
	}

	return nil, false
}
