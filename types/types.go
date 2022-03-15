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
	msg  string
	line int
	col  int
}

func (et *ErrorType) Type() GlimmerType {
	return ERROR
}
func (et *ErrorType) String() string {
	return fmt.Sprintf("[%d,%d]: %s", et.line, et.col, et.msg)
}
