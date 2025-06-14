package ast

import (
	"fmt"
)

// Node represents a node in the AST
type Node interface {
	String() string
}

// Script represents the root of a DSL script
type Script struct {
	Statements []Statement
}

func (s *Script) String() string {
	return fmt.Sprintf("Script{Statements: %v}", s.Statements)
}

// Statement represents a statement in the DSL
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the DSL
type Expression interface {
	Node
	expressionNode()
}

// ConnectStatement connects to an MCP server
type ConnectStatement struct {
	Server     Expression
	Args       []Expression
	Options    map[string]Expression
}

func (c *ConnectStatement) statementNode() {}
func (c *ConnectStatement) String() string {
	return fmt.Sprintf("Connect{Server: %v, Args: %v}", c.Server, c.Args)
}

// CallStatement calls an MCP tool
type CallStatement struct {
	Tool       string
	Arguments  Expression
	Variable   string // Optional: store result in variable
}

func (c *CallStatement) statementNode() {}
func (c *CallStatement) String() string {
	return fmt.Sprintf("Call{Tool: %s, Args: %v, Var: %s}", c.Tool, c.Arguments, c.Variable)
}

// AssertStatement makes an assertion
type AssertStatement struct {
	Expression Expression
	Message    string
}

func (a *AssertStatement) statementNode() {}
func (a *AssertStatement) String() string {
	return fmt.Sprintf("Assert{Expr: %v, Msg: %s}", a.Expression, a.Message)
}

// WaitStatement waits for a condition
type WaitStatement struct {
	Condition Expression
	Timeout   Expression
	Interval  Expression
}

func (w *WaitStatement) statementNode() {}
func (w *WaitStatement) String() string {
	return fmt.Sprintf("Wait{Condition: %v, Timeout: %v}", w.Condition, w.Timeout)
}

// LoopStatement represents a loop
type LoopStatement struct {
	Iterator   string
	Collection Expression
	Body       []Statement
}

func (l *LoopStatement) statementNode() {}
func (l *LoopStatement) String() string {
	return fmt.Sprintf("Loop{Iterator: %s, Collection: %v, Body: %v}", l.Iterator, l.Collection, l.Body)
}

// IfStatement represents a conditional
type IfStatement struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
}

func (i *IfStatement) statementNode() {}
func (i *IfStatement) String() string {
	return fmt.Sprintf("If{Condition: %v, Then: %v, Else: %v}", i.Condition, i.Then, i.Else)
}

// SetStatement sets a variable
type SetStatement struct {
	Variable string
	Value    Expression
}

func (s *SetStatement) statementNode() {}
func (s *SetStatement) String() string {
	return fmt.Sprintf("Set{Var: %s, Value: %v}", s.Variable, s.Value)
}

// PrintStatement prints output
type PrintStatement struct {
	Expression Expression
}

func (p *PrintStatement) statementNode() {}
func (p *PrintStatement) String() string {
	return fmt.Sprintf("Print{Expr: %v}", p.Expression)
}

// DefineStatement defines a reusable automation
type DefineStatement struct {
	Name       string
	Parameters []string
	Body       []Statement
}

func (d *DefineStatement) statementNode() {}
func (d *DefineStatement) String() string {
	return fmt.Sprintf("Define{Name: %s, Params: %v, Body: %v}", d.Name, d.Parameters, d.Body)
}

// RunStatement runs a defined automation
type RunStatement struct {
	Name      string
	Arguments []Expression
}

func (r *RunStatement) statementNode() {}
func (r *RunStatement) String() string {
	return fmt.Sprintf("Run{Name: %s, Args: %v}", r.Name, r.Arguments)
}

// StringLiteral represents a string value
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) String() string {
	return fmt.Sprintf("String{%q}", s.Value)
}

// NumberLiteral represents a numeric value
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) expressionNode() {}
func (n *NumberLiteral) String() string {
	return fmt.Sprintf("Number{%v}", n.Value)
}

// BooleanLiteral represents a boolean value
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("Boolean{%v}", b.Value)
}

// ObjectLiteral represents an object/map
type ObjectLiteral struct {
	Fields map[string]Expression
}

func (o *ObjectLiteral) expressionNode() {}
func (o *ObjectLiteral) String() string {
	return fmt.Sprintf("Object{%v}", o.Fields)
}

// ArrayLiteral represents an array
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) String() string {
	return fmt.Sprintf("Array{%v}", a.Elements)
}

// Variable represents a variable reference
type Variable struct {
	Name string
}

func (v *Variable) expressionNode() {}
func (v *Variable) String() string {
	return fmt.Sprintf("Var{%s}", v.Name)
}

// FieldAccess represents accessing a field
type FieldAccess struct {
	Object Expression
	Field  string
}

func (f *FieldAccess) expressionNode() {}
func (f *FieldAccess) String() string {
	return fmt.Sprintf("Field{%v.%s}", f.Object, f.Field)
}

// IndexAccess represents array/map access
type IndexAccess struct {
	Object Expression
	Index  Expression
}

func (i *IndexAccess) expressionNode() {}
func (i *IndexAccess) String() string {
	return fmt.Sprintf("Index{%v[%v]}", i.Object, i.Index)
}

// BinaryOp represents a binary operation
type BinaryOp struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryOp) expressionNode() {}
func (b *BinaryOp) String() string {
	return fmt.Sprintf("BinaryOp{%v %s %v}", b.Left, b.Operator, b.Right)
}

// UnaryOp represents a unary operation
type UnaryOp struct {
	Operator string
	Operand  Expression
}

func (u *UnaryOp) expressionNode() {}
func (u *UnaryOp) String() string {
	return fmt.Sprintf("UnaryOp{%s%v}", u.Operator, u.Operand)
}

// FunctionCall represents a built-in function call
type FunctionCall struct {
	Name      string
	Arguments []Expression
}

func (f *FunctionCall) expressionNode() {}
func (f *FunctionCall) String() string {
	return fmt.Sprintf("Func{%s(%v)}", f.Name, f.Arguments)
}