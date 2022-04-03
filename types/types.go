package types

import (
	"bytes"
	"fmt"
)

type GlimmerType string

const (
	INTEGER  = "INTEGER"
	FLOAT    = "FLOAT"
	BOOLEAN  = "BOOLEAN"
	STRING   = "STRING"
	ARRAY    = "ARRAY"
	DICT     = "DICT"
	FUNCTION = "FUNCTION"
	NONE     = "NONE"
	ERROR    = "ERROR"
)

type TypeNode interface {
	Type() GlimmerType
	String() string
}

type IntegerType struct{}

func (it *IntegerType) Type() GlimmerType {
	return INTEGER
}
func (it *IntegerType) String() string {
	return "int"
}

type FloatType struct{}

func (ft *FloatType) Type() GlimmerType {
	return FLOAT
}
func (ft *FloatType) String() string {
	return "float"
}

type BooleanType struct{}

func (bt *BooleanType) Type() GlimmerType {
	return BOOLEAN
}
func (bt *BooleanType) String() string {
	return "bool"
}

type StringType struct{}

func (st *StringType) Type() GlimmerType {
	return STRING
}
func (st *StringType) String() string {
	return "string"
}

type ArrayType struct {
	HeldType TypeNode
}

func (at *ArrayType) Type() GlimmerType {
	return ARRAY
}
func (at *ArrayType) String() string {
	return "array[" + at.HeldType.String() + "]"
}

type DictType struct {
	HeldType TypeNode
}

func (dt *DictType) Type() GlimmerType {
	return DICT
}
func (dt *DictType) String() string {
	return "dict[" + dt.HeldType.String() + "]"
}

type FunctionType struct {
	ParamTypes []TypeNode
	ReturnType TypeNode
	FnCtx      *Context
}

func (ft *FunctionType) Type() GlimmerType {
	return FUNCTION
}
func (ft *FunctionType) String() string {
	if len(ft.ParamTypes) == 0 {
		return "fn() -> " + ft.ReturnType.String()
	}

	var out bytes.Buffer
	out.WriteString("fn(")

	out.WriteString(ft.ParamTypes[0].String())
	for _, typ := range ft.ParamTypes[1:len(ft.ParamTypes)] {
		out.WriteString(", " + typ.String())
	}

	out.WriteString(") -> " + ft.ReturnType.String())

	return out.String()
}

type NoneType struct{}

func (nt *NoneType) Type() GlimmerType {
	return NONE
}
func (nt *NoneType) String() string {
	return "none"
}

type ErrorType struct {
	Msg  string
	Line int
	Col  int
}

func (et *ErrorType) Type() GlimmerType {
	return ERROR
}
func (et *ErrorType) String() string {
	return fmt.Sprintf("Static TypeError at [%d,%d]: %s", et.Line, et.Col, et.Msg)
}

func NewEnclosedContext(outer *Context, retType *TypeNode) *Context {
	ctx := NewContext()
	ctx.outer = outer
	ctx.FnType = retType
	return ctx
}

func NewContext() *Context {
	s := make(map[string]TypeNode)
	return &Context{store: s, outer: nil}
}

type Context struct {
	store  map[string]TypeNode
	outer  *Context
	FnType *TypeNode
}

func (c *Context) Get(name string) (TypeNode, bool) {
	typ, ok := c.store[name]
	if !ok && c.outer != nil {
		typ, ok = c.outer.Get(name)
	}
	return typ, ok
}

func (c *Context) Set(name string, val TypeNode) TypeNode {
	c.store[name] = val
	return val
}

func (c *Context) DeepCopy() *Context {
	newEnv := &Context{}
	if c.outer != nil {
		newEnv.outer = c.outer.DeepCopy()
	}
	newStore := make(map[string]TypeNode)
	for key, val := range c.store {
		newStore[key] = val
	}
	newEnv.store = newStore
	return newEnv
}
