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

	_, err := runFromFile(os.Args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
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

	obj := evaluator.Evaluate(program, env)

	if errorObj, ok := obj.(*object.Error); ok {
		return obj, errors.New(errorObj.Message)
	}

	return obj, nil
}

/*
func runFromMultipleFiles() (object.Object, error) {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if splits := strings.Split(path, "."); splits[1] == "gopp" {

			}

			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

*/
