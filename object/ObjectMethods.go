package object

type ObjectMembers struct {
	Members        map[string]Object
	MutableMembers bool
}

func NewMembers(members map[string]Object, isMutable bool) *ObjectMembers {
	return &ObjectMembers{members, isMutable}
}

//TODO add methods
/*func (members *ObjectMethods[T]) Add(name string, obj Object) bool {
	if !members.isMutable {
		return false
	}

	members.methods[name] = ObjectMethod[T]{
		function: ,
	}

	return true
}*/

func (members *ObjectMembers) Get(name string) (Object, bool) {
	val, ok := members.Members[name]

	if !ok {
		return nil, ok
	}

	return val, ok
}
