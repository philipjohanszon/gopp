package evaluator

import (
	"bytes"
	"fmt"
	"go++/object"
)

var builtins = map[string]*object.Builtin{
	"println": {
		Fn: func(args ...object.Object) object.Object {
			fmt.Println(getStringFromArgs(args...))

			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			fmt.Print(getStringFromArgs(args...))

			return NULL
		},
	},
	"printf": {
		Fn: func(args ...object.Object) object.Object {
			format, ok := args[0].(*object.String)

			if !ok {
				return newError("ERROR: first argument in printf must be a format string")
			}

			fmt.Printf(format.Value, getStringFromArgs(args[1:]...))

			return NULL
		},
	},
}

func getStringFromArgs(args ...object.Object) string {
	var out bytes.Buffer

	for _, arg := range args {
		out.WriteString(fmt.Sprint(arg.Inspect()))
	}

	return out.String()
}
