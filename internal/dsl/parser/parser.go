package parser

import (
	"fmt"
	"strconv"

	"github.com/periplon/bract/internal/dsl/ast"
)

// Parser parses DSL tokens into an AST
type Parser struct {
	tokens  []Token
	current int
}

// NewParser creates a new parser
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse parses the tokens into an AST
func (p *Parser) Parse() (*ast.Script, error) {
	script := &ast.Script{
		Statements: []ast.Statement{},
	}

	for !p.isAtEnd() {
		// Skip newlines at the top level
		if p.match(TokenNewline) {
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			script.Statements = append(script.Statements, stmt)
		}
	}

	return script, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	// Skip newlines
	p.consumeNewlines()

	switch {
	case p.match(TokenConnect):
		return p.parseConnect()
	case p.match(TokenCall):
		return p.parseCall()
	case p.match(TokenAssert):
		return p.parseAssert()
	case p.match(TokenWait):
		return p.parseWait()
	case p.match(TokenLoop):
		return p.parseLoop()
	case p.match(TokenIf):
		return p.parseIf()
	case p.match(TokenSet):
		return p.parseSet()
	case p.match(TokenPrint):
		return p.parsePrint()
	case p.match(TokenDefine):
		return p.parseDefine()
	case p.match(TokenRun):
		return p.parseRun()
	case p.check(TokenIdentifier):
		// Could be a variable assignment or a call
		if p.checkAhead(TokenAssign) {
			return p.parseSet()
		}
		return nil, fmt.Errorf("unexpected identifier '%s' at line %d", p.peek().Value, p.peek().Line)
	default:
		if p.isAtEnd() {
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected token '%s' at line %d", p.peek().Value, p.peek().Line)
	}
}

func (p *Parser) parseConnect() (ast.Statement, error) {
	stmt := &ast.ConnectStatement{
		Args:    []ast.Expression{},
		Options: make(map[string]ast.Expression),
	}

	// Parse server expression
	server, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected server expression after 'connect': %w", err)
	}
	stmt.Server = server

	// Parse optional arguments
	for !p.checkNewlineOrEOF() && !p.check(TokenLeftBrace) {
		arg, err := p.parseExpression()
		if err != nil {
			break
		}
		stmt.Args = append(stmt.Args, arg)
	}

	// Parse optional options block
	if p.match(TokenLeftBrace) {
		for !p.check(TokenRightBrace) && !p.isAtEnd() {
			p.consumeNewlines()
			
			if p.check(TokenRightBrace) {
				break
			}

			// Parse option name
			if !p.check(TokenIdentifier) {
				return nil, fmt.Errorf("expected option name at line %d", p.peek().Line)
			}
			optionName := p.advance().Value

			if !p.match(TokenColon) {
				return nil, fmt.Errorf("expected ':' after option name at line %d", p.peek().Line)
			}

			// Parse option value
			value, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("expected value for option '%s': %w", optionName, err)
			}

			stmt.Options[optionName] = value
			p.consumeNewlines()
		}

		if !p.match(TokenRightBrace) {
			return nil, fmt.Errorf("expected '}' at line %d", p.peek().Line)
		}
	}

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseCall() (ast.Statement, error) {
	stmt := &ast.CallStatement{}

	// Parse tool name
	if !p.check(TokenIdentifier) && !p.check(TokenString) {
		return nil, fmt.Errorf("expected tool name after 'call' at line %d", p.peek().Line)
	}
	
	toolToken := p.advance()
	stmt.Tool = toolToken.Value

	// Parse optional arguments
	if !p.checkNewlineOrEOF() && !p.check(TokenArrow) {
		args, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		stmt.Arguments = args
	}

	// Parse optional result variable
	if p.match(TokenArrow) {
		if !p.check(TokenIdentifier) {
			return nil, fmt.Errorf("expected variable name after '->' at line %d", p.peek().Line)
		}
		stmt.Variable = p.advance().Value
	}

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseAssert() (ast.Statement, error) {
	stmt := &ast.AssertStatement{}

	// Parse condition
	expr, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected expression after 'assert': %w", err)
	}
	stmt.Expression = expr

	// Parse optional message
	if p.match(TokenComma) {
		if !p.check(TokenString) {
			return nil, fmt.Errorf("expected string message after ',' at line %d", p.peek().Line)
		}
		stmt.Message = p.advance().Value
	}

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseWait() (ast.Statement, error) {
	stmt := &ast.WaitStatement{}

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected condition after 'wait': %w", err)
	}
	stmt.Condition = condition

	// Parse optional timeout
	if p.match(TokenComma) {
		timeout, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("expected timeout expression: %w", err)
		}
		stmt.Timeout = timeout

		// Parse optional interval
		if p.match(TokenComma) {
			interval, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("expected interval expression: %w", err)
			}
			stmt.Interval = interval
		}
	}

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseLoop() (ast.Statement, error) {
	stmt := &ast.LoopStatement{}

	// Parse iterator variable
	if !p.check(TokenIdentifier) {
		return nil, fmt.Errorf("expected iterator variable after 'loop' at line %d", p.peek().Line)
	}
	stmt.Iterator = p.advance().Value

	// Expect 'in'
	if !p.match(TokenIn) {
		return nil, fmt.Errorf("expected 'in' after iterator variable at line %d", p.peek().Line)
	}

	// Parse collection
	collection, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected collection expression: %w", err)
	}
	stmt.Collection = collection

	// Parse body
	body, err := p.parseBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to parse loop body: %w", err)
	}
	stmt.Body = body

	return stmt, nil
}

func (p *Parser) parseIf() (ast.Statement, error) {
	stmt := &ast.IfStatement{}

	// Parse condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected condition after 'if': %w", err)
	}
	stmt.Condition = condition

	// Parse then block
	thenBlock, err := p.parseBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to parse if body: %w", err)
	}
	stmt.Then = thenBlock

	// Parse optional else
	if p.match(TokenElse) {
		if p.check(TokenIf) {
			// else if - parse as nested if
			elseIf, err := p.parseIf()
			if err != nil {
				return nil, err
			}
			stmt.Else = []ast.Statement{elseIf}
		} else {
			// else block
			elseBlock, err := p.parseBlock()
			if err != nil {
				return nil, fmt.Errorf("failed to parse else body: %w", err)
			}
			stmt.Else = elseBlock
		}
	}

	return stmt, nil
}

func (p *Parser) parseSet() (ast.Statement, error) {
	stmt := &ast.SetStatement{}

	// Skip 'set' if present
	p.match(TokenSet)

	// Parse variable name
	if !p.check(TokenIdentifier) {
		return nil, fmt.Errorf("expected variable name at line %d", p.peek().Line)
	}
	stmt.Variable = p.advance().Value

	// Expect '='
	if !p.match(TokenAssign) {
		return nil, fmt.Errorf("expected '=' after variable name at line %d", p.peek().Line)
	}

	// Parse value
	value, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected value expression: %w", err)
	}
	stmt.Value = value

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parsePrint() (ast.Statement, error) {
	stmt := &ast.PrintStatement{}

	// Parse expression
	expr, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("expected expression after 'print': %w", err)
	}
	stmt.Expression = expr

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseDefine() (ast.Statement, error) {
	stmt := &ast.DefineStatement{
		Parameters: []string{},
	}

	// Parse automation name
	if !p.check(TokenIdentifier) {
		return nil, fmt.Errorf("expected automation name after 'define' at line %d", p.peek().Line)
	}
	stmt.Name = p.advance().Value

	// Parse optional parameters
	if p.match(TokenLeftParen) {
		for !p.check(TokenRightParen) && !p.isAtEnd() {
			if !p.check(TokenIdentifier) {
				return nil, fmt.Errorf("expected parameter name at line %d", p.peek().Line)
			}
			stmt.Parameters = append(stmt.Parameters, p.advance().Value)

			if !p.check(TokenRightParen) {
				if !p.match(TokenComma) {
					return nil, fmt.Errorf("expected ',' or ')' at line %d", p.peek().Line)
				}
			}
		}

		if !p.match(TokenRightParen) {
			return nil, fmt.Errorf("expected ')' at line %d", p.peek().Line)
		}
	}

	// Parse body
	body, err := p.parseBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to parse automation body: %w", err)
	}
	stmt.Body = body

	return stmt, nil
}

func (p *Parser) parseRun() (ast.Statement, error) {
	stmt := &ast.RunStatement{
		Arguments: []ast.Expression{},
	}

	// Parse automation name
	if !p.check(TokenIdentifier) {
		return nil, fmt.Errorf("expected automation name after 'run' at line %d", p.peek().Line)
	}
	stmt.Name = p.advance().Value

	// Parse optional arguments
	if p.match(TokenLeftParen) {
		for !p.check(TokenRightParen) && !p.isAtEnd() {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("failed to parse argument: %w", err)
			}
			stmt.Arguments = append(stmt.Arguments, arg)

			if !p.check(TokenRightParen) {
				if !p.match(TokenComma) {
					return nil, fmt.Errorf("expected ',' or ')' at line %d", p.peek().Line)
				}
			}
		}

		if !p.match(TokenRightParen) {
			return nil, fmt.Errorf("expected ')' at line %d", p.peek().Line)
		}
	}

	p.consumeNewlines()
	return stmt, nil
}

func (p *Parser) parseBlock() ([]ast.Statement, error) {
	statements := []ast.Statement{}

	p.consumeNewlines()
	if !p.match(TokenLeftBrace) {
		return nil, fmt.Errorf("expected '{' at line %d", p.peek().Line)
	}
	p.consumeNewlines()

	for !p.check(TokenRightBrace) && !p.isAtEnd() {
		if p.match(TokenNewline) {
			continue
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	if !p.match(TokenRightBrace) {
		return nil, fmt.Errorf("expected '}' at line %d", p.peek().Line)
	}
	p.consumeNewlines()

	return statements, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (ast.Expression, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.match(TokenOr) {
		op := p.previous().Value
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseAnd() (ast.Expression, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.match(TokenAnd) {
		op := p.previous().Value
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseEquality() (ast.Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(TokenEquals, TokenNotEquals) {
		op := p.previous().Value
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseComparison() (ast.Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(TokenLess, TokenGreater, TokenLessEqual, TokenGreaterEqual) {
		op := p.previous().Value
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseTerm() (ast.Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(TokenPlus, TokenMinus) {
		op := p.previous().Value
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseFactor() (ast.Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.match(TokenMultiply, TokenDivide) {
		op := p.previous().Value
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	if p.match(TokenNot, TokenMinus) {
		op := p.previous().Value
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{
			Operator: op,
			Operand:  operand,
		}, nil
	}

	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (ast.Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(TokenDot) {
			// Field access
			if !p.check(TokenIdentifier) {
				return nil, fmt.Errorf("expected field name after '.' at line %d", p.peek().Line)
			}
			field := p.advance().Value
			expr = &ast.FieldAccess{
				Object: expr,
				Field:  field,
			}
		} else if p.match(TokenLeftBracket) {
			// Index access
			index, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("expected index expression: %w", err)
			}
			if !p.match(TokenRightBracket) {
				return nil, fmt.Errorf("expected ']' at line %d", p.peek().Line)
			}
			expr = &ast.IndexAccess{
				Object: expr,
				Index:  index,
			}
		} else if p.check(TokenLeftParen) && p.previous().Type == TokenIdentifier {
			// Function call (only for identifiers)
			if ident, ok := expr.(*ast.Variable); ok {
				p.advance() // consume '('
				args := []ast.Expression{}
				
				for !p.check(TokenRightParen) && !p.isAtEnd() {
					arg, err := p.parseExpression()
					if err != nil {
						return nil, fmt.Errorf("failed to parse function argument: %w", err)
					}
					args = append(args, arg)

					if !p.check(TokenRightParen) {
						if !p.match(TokenComma) {
							return nil, fmt.Errorf("expected ',' or ')' at line %d", p.peek().Line)
						}
					}
				}

				if !p.match(TokenRightParen) {
					return nil, fmt.Errorf("expected ')' at line %d", p.peek().Line)
				}

				expr = &ast.FunctionCall{
					Name:      ident.Name,
					Arguments: args,
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) parsePrimary() (ast.Expression, error) {
	// String literal
	if p.match(TokenString) {
		return &ast.StringLiteral{Value: p.previous().Value}, nil
	}

	// Number literal
	if p.match(TokenNumber) {
		value, err := strconv.ParseFloat(p.previous().Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %w", err)
		}
		return &ast.NumberLiteral{Value: value}, nil
	}

	// Boolean literal
	if p.match(TokenTrue, TokenFalse) {
		return &ast.BooleanLiteral{Value: p.previous().Type == TokenTrue}, nil
	}

	// Null literal
	if p.match(TokenNull) {
		return &ast.StringLiteral{Value: ""}, nil // Represent null as empty string for now
	}

	// Object literal
	if p.match(TokenLeftBrace) {
		fields := make(map[string]ast.Expression)
		
		for !p.check(TokenRightBrace) && !p.isAtEnd() {
			p.consumeNewlines()
			
			if p.check(TokenRightBrace) {
				break
			}

			// Parse field name
			var fieldName string
			if p.check(TokenString) {
				fieldName = p.advance().Value
			} else if p.check(TokenIdentifier) {
				fieldName = p.advance().Value
			} else {
				return nil, fmt.Errorf("expected field name at line %d", p.peek().Line)
			}

			if !p.match(TokenColon) {
				return nil, fmt.Errorf("expected ':' after field name at line %d", p.peek().Line)
			}

			// Parse field value
			value, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("expected field value: %w", err)
			}

			fields[fieldName] = value

			// Optional comma
			p.match(TokenComma)
			p.consumeNewlines()
		}

		if !p.match(TokenRightBrace) {
			return nil, fmt.Errorf("expected '}' at line %d", p.peek().Line)
		}

		return &ast.ObjectLiteral{Fields: fields}, nil
	}

	// Array literal
	if p.match(TokenLeftBracket) {
		elements := []ast.Expression{}
		
		for !p.check(TokenRightBracket) && !p.isAtEnd() {
			elem, err := p.parseExpression()
			if err != nil {
				return nil, fmt.Errorf("expected array element: %w", err)
			}
			elements = append(elements, elem)

			if !p.check(TokenRightBracket) {
				if !p.match(TokenComma) {
					return nil, fmt.Errorf("expected ',' or ']' at line %d", p.peek().Line)
				}
			}
		}

		if !p.match(TokenRightBracket) {
			return nil, fmt.Errorf("expected ']' at line %d", p.peek().Line)
		}

		return &ast.ArrayLiteral{Elements: elements}, nil
	}

	// Parenthesized expression
	if p.match(TokenLeftParen) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.match(TokenRightParen) {
			return nil, fmt.Errorf("expected ')' at line %d", p.peek().Line)
		}
		return expr, nil
	}

	// Variable or function call
	if p.check(TokenIdentifier) {
		name := p.advance().Value
		return &ast.Variable{Name: name}, nil
	}

	return nil, fmt.Errorf("expected expression at line %d", p.peek().Line)
}

// Helper methods

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) checkAhead(t TokenType) bool {
	if p.current+1 >= len(p.tokens) {
		return false
	}
	return p.tokens[p.current+1].Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenEOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) checkNewlineOrEOF() bool {
	return p.check(TokenNewline) || p.isAtEnd()
}

func (p *Parser) consumeNewlines() {
	for p.match(TokenNewline) {
		// consume all newlines
	}
}