package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/periplon/bract/internal/dsl/ast"
	"github.com/periplon/bract/internal/mcpclient"
)

// Runtime executes DSL scripts
type Runtime struct {
	client      *mcpclient.Client
	variables   map[string]interface{}
	automations map[string]*ast.DefineStatement
	output      strings.Builder
}

// NewRuntime creates a new runtime
func NewRuntime() *Runtime {
	return &Runtime{
		variables:   make(map[string]interface{}),
		automations: make(map[string]*ast.DefineStatement),
	}
}

// Execute runs a DSL script
func (rt *Runtime) Execute(ctx context.Context, script *ast.Script) error {
	for _, stmt := range script.Statements {
		if err := rt.executeStatement(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

// GetOutput returns the accumulated output
func (rt *Runtime) GetOutput() string {
	return rt.output.String()
}

// GetClient returns the MCP client
func (rt *Runtime) GetClient() *mcpclient.Client {
	return rt.client
}

func (rt *Runtime) executeStatement(ctx context.Context, stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ConnectStatement:
		return rt.executeConnect(ctx, s)
	case *ast.CallStatement:
		return rt.executeCall(ctx, s)
	case *ast.AssertStatement:
		return rt.executeAssert(ctx, s)
	case *ast.WaitStatement:
		return rt.executeWait(ctx, s)
	case *ast.LoopStatement:
		return rt.executeLoop(ctx, s)
	case *ast.IfStatement:
		return rt.executeIf(ctx, s)
	case *ast.SetStatement:
		return rt.executeSet(ctx, s)
	case *ast.PrintStatement:
		return rt.executePrint(ctx, s)
	case *ast.DefineStatement:
		return rt.executeDefine(ctx, s)
	case *ast.RunStatement:
		return rt.executeRun(ctx, s)
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (rt *Runtime) executeConnect(ctx context.Context, stmt *ast.ConnectStatement) error {
	// Evaluate server expression
	serverVal, err := rt.evaluateExpression(ctx, stmt.Server)
	if err != nil {
		return fmt.Errorf("failed to evaluate server expression: %w", err)
	}

	serverCmd, ok := serverVal.(string)
	if !ok {
		return fmt.Errorf("server must be a string, got %T", serverVal)
	}

	// Evaluate arguments
	args := []string{}
	for _, argExpr := range stmt.Args {
		argVal, err := rt.evaluateExpression(ctx, argExpr)
		if err != nil {
			return fmt.Errorf("failed to evaluate argument: %w", err)
		}
		args = append(args, fmt.Sprintf("%v", argVal))
	}

	// Create and connect client
	client, err := mcpclient.NewClient(mcpclient.Config{
		ServerCommand: serverCmd,
		ServerArgs:    args,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Connect(ctx, serverCmd, args...); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	rt.client = client
	return nil
}

func (rt *Runtime) executeCall(ctx context.Context, stmt *ast.CallStatement) error {
	if rt.client == nil {
		return fmt.Errorf("not connected to any MCP server")
	}

	// Handle special built-in tools
	if stmt.Tool == "list_tools" {
		tools, err := rt.client.ListTools(ctx)
		if err != nil {
			return fmt.Errorf("failed to list tools: %w", err)
		}

		// Store result if variable specified
		if stmt.Variable != "" {
			// Convert tools to a simple list of names
			toolNames := make([]interface{}, len(tools))
			for i, tool := range tools {
				toolNames[i] = tool.Name
			}
			rt.variables[stmt.Variable] = toolNames
		}
		return nil
	}

	// Evaluate arguments
	var args interface{}
	if stmt.Arguments != nil {
		var err error
		args, err = rt.evaluateExpression(ctx, stmt.Arguments)
		if err != nil {
			return fmt.Errorf("failed to evaluate arguments: %w", err)
		}
	}

	// Call the tool
	result, err := rt.client.CallTool(ctx, stmt.Tool, args)
	if err != nil {
		return fmt.Errorf("tool call failed: %w", err)
	}

	// Store result if variable specified
	if stmt.Variable != "" {
		// Extract content from result
		var resultValue interface{}
		if len(result.Content) == 1 {
			// Single content item - process it
			content := result.Content[0]
			// Check if it's a TextContent with JSON data
			if textContent, ok := content.(mcp.TextContent); ok {
				// Try to parse as JSON
				var jsonData interface{}
				if err := json.Unmarshal([]byte(textContent.Text), &jsonData); err == nil {
					resultValue = jsonData
				} else {
					// Not JSON, store as plain text
					resultValue = textContent.Text
				}
			} else {
				// Not TextContent, store as is
				resultValue = content
			}
		} else {
			// Multiple content items - process each one
			items := make([]interface{}, len(result.Content))
			for i, content := range result.Content {
				if textContent, ok := content.(mcp.TextContent); ok {
					// Try to parse as JSON
					var jsonData interface{}
					if err := json.Unmarshal([]byte(textContent.Text), &jsonData); err == nil {
						items[i] = jsonData
					} else {
						// Not JSON, store as plain text
						items[i] = textContent.Text
					}
				} else {
					// Not TextContent, store as is
					items[i] = content
				}
			}
			resultValue = items
		}
		rt.variables[stmt.Variable] = resultValue
	}

	return nil
}

func (rt *Runtime) executeAssert(ctx context.Context, stmt *ast.AssertStatement) error {
	// Evaluate condition
	result, err := rt.evaluateExpression(ctx, stmt.Expression)
	if err != nil {
		return fmt.Errorf("failed to evaluate assertion: %w", err)
	}

	// Check if condition is true
	isTrue := rt.isTruthy(result)
	if !isTrue {
		msg := stmt.Message
		if msg == "" {
			msg = fmt.Sprintf("assertion failed: %v", stmt.Expression)
		}
		return fmt.Errorf("assertion error: %s", msg)
	}

	return nil
}

func (rt *Runtime) executeWait(ctx context.Context, stmt *ast.WaitStatement) error {
	// Default timeout and interval
	timeout := 30 * time.Second
	interval := 100 * time.Millisecond

	// Evaluate timeout if provided
	if stmt.Timeout != nil {
		timeoutVal, err := rt.evaluateExpression(ctx, stmt.Timeout)
		if err != nil {
			return fmt.Errorf("failed to evaluate timeout: %w", err)
		}
		switch v := timeoutVal.(type) {
		case float64:
			timeout = time.Duration(v) * time.Second
		case int:
			timeout = time.Duration(v) * time.Second
		default:
			return fmt.Errorf("timeout must be a number, got %T", timeoutVal)
		}
	}

	// Evaluate interval if provided
	if stmt.Interval != nil {
		intervalVal, err := rt.evaluateExpression(ctx, stmt.Interval)
		if err != nil {
			return fmt.Errorf("failed to evaluate interval: %w", err)
		}
		switch v := intervalVal.(type) {
		case float64:
			interval = time.Duration(v) * time.Millisecond
		case int:
			interval = time.Duration(v) * time.Millisecond
		default:
			return fmt.Errorf("interval must be a number, got %T", intervalVal)
		}
	}

	// Wait for condition
	deadline := time.Now().Add(timeout)
	for {
		result, err := rt.evaluateExpression(ctx, stmt.Condition)
		if err != nil {
			return fmt.Errorf("failed to evaluate wait condition: %w", err)
		}

		if rt.isTruthy(result) {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("wait timeout: condition did not become true within %v", timeout)
		}

		time.Sleep(interval)
	}
}

func (rt *Runtime) executeLoop(ctx context.Context, stmt *ast.LoopStatement) error {
	// Evaluate collection
	collection, err := rt.evaluateExpression(ctx, stmt.Collection)
	if err != nil {
		return fmt.Errorf("failed to evaluate collection: %w", err)
	}

	// Convert to iterable
	items, err := rt.toIterable(collection)
	if err != nil {
		return fmt.Errorf("collection is not iterable: %w", err)
	}

	// Execute loop body for each item
	for _, item := range items {
		// Set iterator variable
		rt.variables[stmt.Iterator] = item

		// Execute body
		for _, bodyStmt := range stmt.Body {
			if err := rt.executeStatement(ctx, bodyStmt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rt *Runtime) executeIf(ctx context.Context, stmt *ast.IfStatement) error {
	// Evaluate condition
	condition, err := rt.evaluateExpression(ctx, stmt.Condition)
	if err != nil {
		return fmt.Errorf("failed to evaluate if condition: %w", err)
	}

	// Execute appropriate branch
	if rt.isTruthy(condition) {
		for _, thenStmt := range stmt.Then {
			if err := rt.executeStatement(ctx, thenStmt); err != nil {
				return err
			}
		}
	} else if stmt.Else != nil {
		for _, elseStmt := range stmt.Else {
			if err := rt.executeStatement(ctx, elseStmt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rt *Runtime) executeSet(ctx context.Context, stmt *ast.SetStatement) error {
	// Evaluate value
	value, err := rt.evaluateExpression(ctx, stmt.Value)
	if err != nil {
		return fmt.Errorf("failed to evaluate value: %w", err)
	}

	// Set variable
	rt.variables[stmt.Variable] = value
	return nil
}

func (rt *Runtime) executePrint(ctx context.Context, stmt *ast.PrintStatement) error {
	// Evaluate expression
	value, err := rt.evaluateExpression(ctx, stmt.Expression)
	if err != nil {
		return fmt.Errorf("failed to evaluate print expression: %w", err)
	}

	// Format and print value
	output := rt.formatValue(value)
	fmt.Fprintln(&rt.output, output)
	fmt.Println(output)
	return nil
}

func (rt *Runtime) executeDefine(ctx context.Context, stmt *ast.DefineStatement) error {
	// Store automation definition
	rt.automations[stmt.Name] = stmt
	return nil
}

func (rt *Runtime) executeRun(ctx context.Context, stmt *ast.RunStatement) error {
	// Find automation
	automation, ok := rt.automations[stmt.Name]
	if !ok {
		return fmt.Errorf("automation '%s' not defined", stmt.Name)
	}

	// Check argument count
	if len(stmt.Arguments) != len(automation.Parameters) {
		return fmt.Errorf("automation '%s' expects %d arguments, got %d",
			stmt.Name, len(automation.Parameters), len(stmt.Arguments))
	}

	// Create new scope with parameters
	oldVars := rt.variables
	rt.variables = make(map[string]interface{})
	for k, v := range oldVars {
		rt.variables[k] = v
	}

	// Bind arguments to parameters
	for i, param := range automation.Parameters {
		value, err := rt.evaluateExpression(ctx, stmt.Arguments[i])
		if err != nil {
			rt.variables = oldVars
			return fmt.Errorf("failed to evaluate argument %d: %w", i, err)
		}
		rt.variables[param] = value
	}

	// Execute automation body
	for _, bodyStmt := range automation.Body {
		if err := rt.executeStatement(ctx, bodyStmt); err != nil {
			rt.variables = oldVars
			return err
		}
	}

	// Restore scope
	rt.variables = oldVars
	return nil
}

func (rt *Runtime) evaluateExpression(ctx context.Context, expr ast.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.NumberLiteral:
		return e.Value, nil
	case *ast.BooleanLiteral:
		return e.Value, nil
	case *ast.Variable:
		val, ok := rt.variables[e.Name]
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", e.Name)
		}
		return val, nil
	case *ast.ObjectLiteral:
		obj := make(map[string]interface{})
		for k, v := range e.Fields {
			val, err := rt.evaluateExpression(ctx, v)
			if err != nil {
				return nil, err
			}
			obj[k] = val
		}
		return obj, nil
	case *ast.ArrayLiteral:
		arr := make([]interface{}, len(e.Elements))
		for i, elem := range e.Elements {
			val, err := rt.evaluateExpression(ctx, elem)
			if err != nil {
				return nil, err
			}
			arr[i] = val
		}
		return arr, nil
	case *ast.FieldAccess:
		obj, err := rt.evaluateExpression(ctx, e.Object)
		if err != nil {
			return nil, err
		}
		return rt.getField(obj, e.Field)
	case *ast.IndexAccess:
		obj, err := rt.evaluateExpression(ctx, e.Object)
		if err != nil {
			return nil, err
		}
		index, err := rt.evaluateExpression(ctx, e.Index)
		if err != nil {
			return nil, err
		}
		return rt.getIndex(obj, index)
	case *ast.BinaryOp:
		return rt.evaluateBinaryOp(ctx, e)
	case *ast.UnaryOp:
		return rt.evaluateUnaryOp(ctx, e)
	case *ast.FunctionCall:
		return rt.evaluateFunctionCall(ctx, e)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (rt *Runtime) evaluateBinaryOp(ctx context.Context, op *ast.BinaryOp) (interface{}, error) {
	left, err := rt.evaluateExpression(ctx, op.Left)
	if err != nil {
		return nil, err
	}

	// Short-circuit evaluation for logical operators
	if op.Operator == "&&" {
		if !rt.isTruthy(left) {
			return false, nil
		}
	} else if op.Operator == "||" {
		if rt.isTruthy(left) {
			return true, nil
		}
	}

	right, err := rt.evaluateExpression(ctx, op.Right)
	if err != nil {
		return nil, err
	}

	switch op.Operator {
	case "+":
		return rt.add(left, right)
	case "-":
		return rt.subtract(left, right)
	case "*":
		return rt.multiply(left, right)
	case "/":
		return rt.divide(left, right)
	case "==":
		return rt.equals(left, right), nil
	case "!=":
		return !rt.equals(left, right), nil
	case "<":
		return rt.lessThan(left, right)
	case ">":
		return rt.greaterThan(left, right)
	case "<=":
		return rt.lessEqual(left, right)
	case ">=":
		return rt.greaterEqual(left, right)
	case "&&":
		return rt.isTruthy(right), nil
	case "||":
		return true, nil
	default:
		return nil, fmt.Errorf("unknown binary operator: %s", op.Operator)
	}
}

func (rt *Runtime) evaluateUnaryOp(ctx context.Context, op *ast.UnaryOp) (interface{}, error) {
	operand, err := rt.evaluateExpression(ctx, op.Operand)
	if err != nil {
		return nil, err
	}

	switch op.Operator {
	case "!":
		return !rt.isTruthy(operand), nil
	case "-":
		switch v := operand.(type) {
		case float64:
			return -v, nil
		case int:
			return -v, nil
		default:
			return nil, fmt.Errorf("cannot negate %T", operand)
		}
	default:
		return nil, fmt.Errorf("unknown unary operator: %s", op.Operator)
	}
}

func (rt *Runtime) evaluateFunctionCall(ctx context.Context, call *ast.FunctionCall) (interface{}, error) {
	// Built-in functions
	switch call.Name {
	case "len":
		if len(call.Arguments) != 1 {
			return nil, fmt.Errorf("len() expects 1 argument, got %d", len(call.Arguments))
		}
		arg, err := rt.evaluateExpression(ctx, call.Arguments[0])
		if err != nil {
			return nil, err
		}
		return rt.getLength(arg)

	case "str":
		if len(call.Arguments) != 1 {
			return nil, fmt.Errorf("str() expects 1 argument, got %d", len(call.Arguments))
		}
		arg, err := rt.evaluateExpression(ctx, call.Arguments[0])
		if err != nil {
			return nil, err
		}
		return rt.formatValue(arg), nil

	case "int":
		if len(call.Arguments) != 1 {
			return nil, fmt.Errorf("int() expects 1 argument, got %d", len(call.Arguments))
		}
		arg, err := rt.evaluateExpression(ctx, call.Arguments[0])
		if err != nil {
			return nil, err
		}
		return rt.toInt(arg)

	case "float":
		if len(call.Arguments) != 1 {
			return nil, fmt.Errorf("float() expects 1 argument, got %d", len(call.Arguments))
		}
		arg, err := rt.evaluateExpression(ctx, call.Arguments[0])
		if err != nil {
			return nil, err
		}
		return rt.toFloat(arg)

	case "json":
		if len(call.Arguments) != 1 {
			return nil, fmt.Errorf("json() expects 1 argument, got %d", len(call.Arguments))
		}
		arg, err := rt.evaluateExpression(ctx, call.Arguments[0])
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		return string(data), nil

	default:
		return nil, fmt.Errorf("unknown function: %s", call.Name)
	}
}

// Helper methods

func (rt *Runtime) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return v != ""
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	default:
		return true
	}
}

func (rt *Runtime) equals(a, b interface{}) bool {
	// Handle nil values
	if a == nil || b == nil {
		return a == b
	}

	// Use deep equality
	return reflect.DeepEqual(a, b)
}

func (rt *Runtime) toIterable(val interface{}) ([]interface{}, error) {
	switch v := val.(type) {
	case []interface{}:
		return v, nil
	case map[string]interface{}:
		// Convert map to array of key-value pairs
		items := make([]interface{}, 0, len(v))
		for k, val := range v {
			items = append(items, map[string]interface{}{
				"key":   k,
				"value": val,
			})
		}
		return items, nil
	case string:
		// Convert string to array of characters
		items := make([]interface{}, len(v))
		for i, ch := range v {
			items[i] = string(ch)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("value is not iterable: %T", val)
	}
}

func (rt *Runtime) getField(obj interface{}, field string) (interface{}, error) {
	switch v := obj.(type) {
	case map[string]interface{}:
		val, ok := v[field]
		if !ok {
			return nil, nil // Return nil for missing fields
		}
		return val, nil
	default:
		return nil, fmt.Errorf("cannot access field '%s' on %T", field, obj)
	}
}

func (rt *Runtime) getIndex(obj, index interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case []interface{}:
		idx, err := rt.toInt(index)
		if err != nil {
			return nil, fmt.Errorf("array index must be integer: %w", err)
		}
		if idx < 0 || idx >= len(v) {
			return nil, fmt.Errorf("array index out of bounds: %d", idx)
		}
		return v[idx], nil
	case map[string]interface{}:
		key, ok := index.(string)
		if !ok {
			key = fmt.Sprintf("%v", index)
		}
		return v[key], nil
	case string:
		idx, err := rt.toInt(index)
		if err != nil {
			return nil, fmt.Errorf("string index must be integer: %w", err)
		}
		if idx < 0 || idx >= len(v) {
			return nil, fmt.Errorf("string index out of bounds: %d", idx)
		}
		return string(v[idx]), nil
	default:
		return nil, fmt.Errorf("cannot index %T", obj)
	}
}

func (rt *Runtime) formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case nil:
		return "null"
	default:
		data, err := json.MarshalIndent(val, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(data)
	}
}

func (rt *Runtime) getLength(val interface{}) (int, error) {
	switch v := val.(type) {
	case string:
		return len(v), nil
	case []interface{}:
		return len(v), nil
	case map[string]interface{}:
		return len(v), nil
	default:
		return 0, fmt.Errorf("cannot get length of %T", val)
	}
}

func (rt *Runtime) toInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		var i int
		_, err := fmt.Sscanf(v, "%d", &i)
		return i, err
	default:
		return 0, fmt.Errorf("cannot convert %T to int", val)
	}
}

func (rt *Runtime) toFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", &f)
		return f, err
	default:
		return 0, fmt.Errorf("cannot convert %T to float", val)
	}
}

func (rt *Runtime) add(a, b interface{}) (interface{}, error) {
	// String concatenation
	if aStr, ok := a.(string); ok {
		return aStr + rt.formatValue(b), nil
	}
	if bStr, ok := b.(string); ok {
		return rt.formatValue(a) + bStr, nil
	}

	// Numeric addition
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		return aFloat + bFloat, nil
	}

	return nil, fmt.Errorf("cannot add %T and %T", a, b)
}

func (rt *Runtime) subtract(a, b interface{}) (interface{}, error) {
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		return aFloat - bFloat, nil
	}
	return nil, fmt.Errorf("cannot subtract %T and %T", a, b)
}

func (rt *Runtime) multiply(a, b interface{}) (interface{}, error) {
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		return aFloat * bFloat, nil
	}
	return nil, fmt.Errorf("cannot multiply %T and %T", a, b)
}

func (rt *Runtime) divide(a, b interface{}) (interface{}, error) {
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		if bFloat == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return aFloat / bFloat, nil
	}
	return nil, fmt.Errorf("cannot divide %T and %T", a, b)
}

func (rt *Runtime) lessThan(a, b interface{}) (bool, error) {
	// String comparison
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && bIsStr {
		return aStr < bStr, nil
	}

	// Numeric comparison
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		return aFloat < bFloat, nil
	}

	return false, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (rt *Runtime) greaterThan(a, b interface{}) (bool, error) {
	// String comparison
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && bIsStr {
		return aStr > bStr, nil
	}

	// Numeric comparison
	aFloat, err1 := rt.toFloat(a)
	bFloat, err2 := rt.toFloat(b)
	if err1 == nil && err2 == nil {
		return aFloat > bFloat, nil
	}

	return false, fmt.Errorf("cannot compare %T and %T", a, b)
}

func (rt *Runtime) lessEqual(a, b interface{}) (bool, error) {
	result, err := rt.greaterThan(a, b)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (rt *Runtime) greaterEqual(a, b interface{}) (bool, error) {
	result, err := rt.lessThan(a, b)
	if err != nil {
		return false, err
	}
	return !result, nil
}
