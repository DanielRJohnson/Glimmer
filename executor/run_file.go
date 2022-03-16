package executor

import (
	"fmt"
	"glimmer/evaluator"
	"glimmer/lexer"
	"glimmer/object"
	"glimmer/parser"
	"glimmer/typechecker"
	"glimmer/types"
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
	ctx := types.NewContext()

	program := p.ParseProgram()
	errors := p.Errors()
	var errObjs []error
	if len(p.Errors()) != 0 {
		for _, err := range errors {
			errObjs = append(errObjs, fmt.Errorf(err))
		}
		return nil, errObjs
	}

	pType := typechecker.Typeof(program, ctx)
	if pType.Type() == types.ERROR {
		errObjs = append(errObjs, fmt.Errorf(pType.String()))
		return nil, errObjs
	}

	if dot {
		program.ToDot()
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	return evaluated, nil
}
