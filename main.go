package main

import (
	"errors"
	"fmt"
	"go++/evaluator"
	"go++/lexer"
	"go++/object"
	"go++/parser"
	"go++/repl"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	result, err := runFromFile(os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(result.Inspect())
}

func runFromFile(file string) (object.Object, error) {
	data, err := os.ReadFile(file)

	if err != nil {
		return nil, errors.New("no file found named " + file)
	}

	lex := lexer.New(string(data))
	pars := parser.New(lex)

	program := pars.ParseProgram()
	env := object.NewEnvironment()

	return evaluator.Evaluate(program, env), nil
}
