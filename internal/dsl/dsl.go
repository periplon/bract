package dsl

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/periplon/bract/internal/dsl/ast"
	"github.com/periplon/bract/internal/dsl/parser"
	"github.com/periplon/bract/internal/dsl/runtime"
	"github.com/periplon/bract/internal/mcpclient"
)

// Interpreter interprets DSL scripts
type Interpreter struct {
	runtime *runtime.Runtime
}

// NewInterpreter creates a new DSL interpreter
func NewInterpreter() *Interpreter {
	return &Interpreter{
		runtime: runtime.NewRuntime(),
	}
}

// ExecuteString executes a DSL script from a string
func (i *Interpreter) ExecuteString(ctx context.Context, script string) error {
	// Parse the script
	ast, err := ParseString(script)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Execute the AST
	if err := i.runtime.Execute(ctx, ast); err != nil {
		return fmt.Errorf("runtime error: %w", err)
	}

	return nil
}

// ExecuteFile executes a DSL script from a file
func (i *Interpreter) ExecuteFile(ctx context.Context, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return i.ExecuteString(ctx, string(data))
}

// ExecuteReader executes a DSL script from a reader
func (i *Interpreter) ExecuteReader(ctx context.Context, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	return i.ExecuteString(ctx, string(data))
}

// GetOutput returns the accumulated output from print statements
func (i *Interpreter) GetOutput() string {
	return i.runtime.GetOutput()
}

// GetClient returns the MCP client
func (i *Interpreter) GetClient() *mcpclient.Client {
	return i.runtime.GetClient()
}

// ParseString parses a DSL script from a string
func ParseString(script string) (*ast.Script, error) {
	// Tokenize
	lexer := parser.NewLexer(script)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	// Parse
	p := parser.NewParser(tokens)
	return p.Parse()
}

// ParseFile parses a DSL script from a file
func ParseFile(filename string) (*ast.Script, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseString(string(data))
}

// ValidateString validates a DSL script syntax without executing it
func ValidateString(script string) error {
	_, err := ParseString(script)
	return err
}

// ValidateFile validates a DSL script file syntax without executing it
func ValidateFile(filename string) error {
	_, err := ParseFile(filename)
	return err
}

// FormatScript formats a DSL script for readability
func FormatScript(script string) (string, error) {
	// Parse the script first to ensure it's valid
	ast, err := ParseString(script)
	if err != nil {
		return "", err
	}

	// Format the AST back to string
	return FormatAST(ast), nil
}

// FormatAST formats an AST back to DSL syntax
func FormatAST(script *ast.Script) string {
	var sb strings.Builder
	formatter := &astFormatter{indent: 0}
	
	for i, stmt := range script.Statements {
		if i > 0 {
			sb.WriteString("\n")
		}
		formatter.formatStatement(&sb, stmt)
	}
	
	return sb.String()
}

type astFormatter struct {
	indent int
}

func (f *astFormatter) formatStatement(sb *strings.Builder, stmt ast.Statement) {
	f.writeIndent(sb)
	
	switch s := stmt.(type) {
	case *ast.ConnectStatement:
		f.formatConnect(sb, s)
	case *ast.CallStatement:
		f.formatCall(sb, s)
	case *ast.AssertStatement:
		f.formatAssert(sb, s)
	case *ast.WaitStatement:
		f.formatWait(sb, s)
	case *ast.LoopStatement:
		f.formatLoop(sb, s)
	case *ast.IfStatement:
		f.formatIf(sb, s)
	case *ast.SetStatement:
		f.formatSet(sb, s)
	case *ast.PrintStatement:
		f.formatPrint(sb, s)
	case *ast.DefineStatement:
		f.formatDefine(sb, s)
	case *ast.RunStatement:
		f.formatRun(sb, s)
	}
}

func (f *astFormatter) formatConnect(sb *strings.Builder, stmt *ast.ConnectStatement) {
	sb.WriteString("connect ")
	f.formatExpression(sb, stmt.Server)
	
	for _, arg := range stmt.Args {
		sb.WriteString(" ")
		f.formatExpression(sb, arg)
	}
	
	if len(stmt.Options) > 0 {
		sb.WriteString(" {\n")
		f.indent++
		for name, value := range stmt.Options {
			f.writeIndent(sb)
			sb.WriteString(name)
			sb.WriteString(": ")
			f.formatExpression(sb, value)
			sb.WriteString("\n")
		}
		f.indent--
		f.writeIndent(sb)
		sb.WriteString("}")
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatCall(sb *strings.Builder, stmt *ast.CallStatement) {
	sb.WriteString("call ")
	sb.WriteString(stmt.Tool)
	
	if stmt.Arguments != nil {
		sb.WriteString(" ")
		f.formatExpression(sb, stmt.Arguments)
	}
	
	if stmt.Variable != "" {
		sb.WriteString(" -> ")
		sb.WriteString(stmt.Variable)
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatAssert(sb *strings.Builder, stmt *ast.AssertStatement) {
	sb.WriteString("assert ")
	f.formatExpression(sb, stmt.Expression)
	
	if stmt.Message != "" {
		sb.WriteString(", ")
		sb.WriteString(fmt.Sprintf("%q", stmt.Message))
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatWait(sb *strings.Builder, stmt *ast.WaitStatement) {
	sb.WriteString("wait ")
	f.formatExpression(sb, stmt.Condition)
	
	if stmt.Timeout != nil {
		sb.WriteString(", ")
		f.formatExpression(sb, stmt.Timeout)
		
		if stmt.Interval != nil {
			sb.WriteString(", ")
			f.formatExpression(sb, stmt.Interval)
		}
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatLoop(sb *strings.Builder, stmt *ast.LoopStatement) {
	sb.WriteString("loop ")
	sb.WriteString(stmt.Iterator)
	sb.WriteString(" in ")
	f.formatExpression(sb, stmt.Collection)
	sb.WriteString(" {\n")
	
	f.indent++
	for _, s := range stmt.Body {
		f.formatStatement(sb, s)
	}
	f.indent--
	
	f.writeIndent(sb)
	sb.WriteString("}\n")
}

func (f *astFormatter) formatIf(sb *strings.Builder, stmt *ast.IfStatement) {
	sb.WriteString("if ")
	f.formatExpression(sb, stmt.Condition)
	sb.WriteString(" {\n")
	
	f.indent++
	for _, s := range stmt.Then {
		f.formatStatement(sb, s)
	}
	f.indent--
	
	f.writeIndent(sb)
	sb.WriteString("}")
	
	if len(stmt.Else) > 0 {
		sb.WriteString(" else ")
		if len(stmt.Else) == 1 {
			if _, ok := stmt.Else[0].(*ast.IfStatement); ok {
				// else if - format inline
				f.formatStatement(sb, stmt.Else[0])
				return
			}
		}
		
		sb.WriteString("{\n")
		f.indent++
		for _, s := range stmt.Else {
			f.formatStatement(sb, s)
		}
		f.indent--
		f.writeIndent(sb)
		sb.WriteString("}")
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatSet(sb *strings.Builder, stmt *ast.SetStatement) {
	sb.WriteString("set ")
	sb.WriteString(stmt.Variable)
	sb.WriteString(" = ")
	f.formatExpression(sb, stmt.Value)
	sb.WriteString("\n")
}

func (f *astFormatter) formatPrint(sb *strings.Builder, stmt *ast.PrintStatement) {
	sb.WriteString("print ")
	f.formatExpression(sb, stmt.Expression)
	sb.WriteString("\n")
}

func (f *astFormatter) formatDefine(sb *strings.Builder, stmt *ast.DefineStatement) {
	sb.WriteString("define ")
	sb.WriteString(stmt.Name)
	
	if len(stmt.Parameters) > 0 {
		sb.WriteString("(")
		for i, param := range stmt.Parameters {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(param)
		}
		sb.WriteString(")")
	}
	
	sb.WriteString(" {\n")
	f.indent++
	for _, s := range stmt.Body {
		f.formatStatement(sb, s)
	}
	f.indent--
	f.writeIndent(sb)
	sb.WriteString("}\n")
}

func (f *astFormatter) formatRun(sb *strings.Builder, stmt *ast.RunStatement) {
	sb.WriteString("run ")
	sb.WriteString(stmt.Name)
	
	if len(stmt.Arguments) > 0 {
		sb.WriteString("(")
		for i, arg := range stmt.Arguments {
			if i > 0 {
				sb.WriteString(", ")
			}
			f.formatExpression(sb, arg)
		}
		sb.WriteString(")")
	}
	sb.WriteString("\n")
}

func (f *astFormatter) formatExpression(sb *strings.Builder, expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		sb.WriteString(fmt.Sprintf("%q", e.Value))
	case *ast.NumberLiteral:
		sb.WriteString(fmt.Sprintf("%v", e.Value))
	case *ast.BooleanLiteral:
		sb.WriteString(fmt.Sprintf("%v", e.Value))
	case *ast.Variable:
		sb.WriteString(e.Name)
	case *ast.ObjectLiteral:
		sb.WriteString("{")
		first := true
		for k, v := range e.Fields {
			if !first {
				sb.WriteString(", ")
			}
			first = false
			sb.WriteString(k)
			sb.WriteString(": ")
			f.formatExpression(sb, v)
		}
		sb.WriteString("}")
	case *ast.ArrayLiteral:
		sb.WriteString("[")
		for i, elem := range e.Elements {
			if i > 0 {
				sb.WriteString(", ")
			}
			f.formatExpression(sb, elem)
		}
		sb.WriteString("]")
	case *ast.FieldAccess:
		f.formatExpression(sb, e.Object)
		sb.WriteString(".")
		sb.WriteString(e.Field)
	case *ast.IndexAccess:
		f.formatExpression(sb, e.Object)
		sb.WriteString("[")
		f.formatExpression(sb, e.Index)
		sb.WriteString("]")
	case *ast.BinaryOp:
		sb.WriteString("(")
		f.formatExpression(sb, e.Left)
		sb.WriteString(" ")
		sb.WriteString(e.Operator)
		sb.WriteString(" ")
		f.formatExpression(sb, e.Right)
		sb.WriteString(")")
	case *ast.UnaryOp:
		sb.WriteString(e.Operator)
		f.formatExpression(sb, e.Operand)
	case *ast.FunctionCall:
		sb.WriteString(e.Name)
		sb.WriteString("(")
		for i, arg := range e.Arguments {
			if i > 0 {
				sb.WriteString(", ")
			}
			f.formatExpression(sb, arg)
		}
		sb.WriteString(")")
	}
}

func (f *astFormatter) writeIndent(sb *strings.Builder) {
	for i := 0; i < f.indent; i++ {
		sb.WriteString("  ")
	}
}