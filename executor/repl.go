package executor

import (
	"bufio"
	"fmt"
	"io"

	"glimmer/evaluator"
	"glimmer/lexer"
	"glimmer/object"
	"glimmer/parser"
	"glimmer/token"
	"glimmer/typechecker"
	"glimmer/types"
)

const PROMPT = ">> "

func StartRLPL(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

func StartRPPL(in io.Reader, out io.Writer, dot bool) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		errors := p.Errors()
		if len(errors) != 0 {
			for _, err := range errors {
				io.WriteString(out, "\t"+err+"\n")
			}
			continue
		}
		io.WriteString(out, program.String()+"\n")

		if dot {
			currTime := program.ToDot()
			io.WriteString(out, "Dot file & image "+currTime+" created in /dot/dotfiles & /dot/dotimages\n")
		}
	}
}

func StartREPL(in io.Reader, out io.Writer, dot bool) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	ctx := types.NewContext()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		errors := p.Errors()
		if len(p.Errors()) != 0 {
			for _, err := range errors {
				io.WriteString(out, err+"\n")
			}
			continue
		}

		pType := typechecker.Typeof(program, ctx)
		if pType.Type() == types.ERROR {
			io.WriteString(out, pType.String()+"\n")
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		if dot {
			currTime := program.ToDot()
			io.WriteString(out, "Dot file & image "+currTime+" created in /dot/dotfiles & /dot/dotimages\n")
		}
	}
}
