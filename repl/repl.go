package repl

import (
	"Interpreter/evaluator"
	"Interpreter/lexer"
	"Interpreter/object"
	"Interpreter/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">>>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	ctx := object.NewContext()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParserProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, ctx)
		if evaluated != nil {
			_, _ = io.WriteString(out, evaluated.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}

}

func printParserErrors(out io.Writer, err []error) {
	for _, err := range err {
		_, _ = io.WriteString(out, "\t"+err.Error()+"\n")
	}
}
