package executor

import (
	"fmt"
	"glimmer/evaluator"
	"glimmer/lexer"
	"glimmer/object"
	"glimmer/parser"
	"io/ioutil"
)

func RunFile(fpath string, dot bool) (object.Object, []error) {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, []error{err}
	}

	contentString := string(content)

	l := lexer.New(contentString)
	p := parser.New(l)

	program := p.ParseProgram()
	errors := p.Errors()
	if len(p.Errors()) != 0 {
		var errObjs []error
		for _, err := range errors {
			errObjs = append(errObjs, fmt.Errorf(err))
		}
		return nil, errObjs
	}
	if dot {
		program.ToDot()
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	return evaluated, nil
}
