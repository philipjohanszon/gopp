package repl

import (
	"bufio"
	"fmt"
	"go++/evaluator"
	"go++/lexer"
	"go++/object"
	"go++/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		pars := parser.New(lexer.New(line))

		program := pars.ParseProgram()

		if len(pars.Errors()) > 0 {
			printParserErrors(out, pars.Errors())
			continue
		}

		evaluated := evaluator.Evaluate(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error(s) occurred!\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
