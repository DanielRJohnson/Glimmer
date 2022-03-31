package ast

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/bradleyjkemp/memviz"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) ToDot() string {
	buf := &bytes.Buffer{}
	memviz.Map(buf, &p)

	// make a new file whose name is unique through the unix epoch
	currTime := fmt.Sprint(time.Now().UnixNano())
	workingDirectory, _ := os.Getwd()
	dotDir := workingDirectory + "/dot-outputs/"
	dotFilePath := dotDir + "dotfiles/" + currTime
	dotImagePath := dotDir + "dotimages/" + currTime + ".png"

	// make dot directories if not exists
	err := os.MkdirAll(dotDir+"dotfiles", os.ModeDir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dotDir+"dotimages", os.ModeDir)
	if err != nil {
		panic(err)
	}

	// write the dotfile
	err = os.WriteFile(dotFilePath, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	// run the dot command with the dotfilepath and dotimagepath
	dotCommand := exec.Command("dot", "-Tpng", dotFilePath, "-o", dotImagePath)
	err = dotCommand.Run()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Nanosecond)

	return currTime
}
