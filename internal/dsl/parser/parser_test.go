package parser

import (
	"testing"

	"github.com/periplon/bract/internal/dsl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:  "simple connect",
			input: `connect "server"`,
			expected: []TokenType{
				TokenConnect,
				TokenString,
				TokenEOF,
			},
		},
		{
			name:  "call with arrow",
			input: `call tool_name -> result`,
			expected: []TokenType{
				TokenCall,
				TokenIdentifier,
				TokenArrow,
				TokenIdentifier,
				TokenEOF,
			},
		},
		{
			name:  "operators",
			input: `== != < > <= >= && || !`,
			expected: []TokenType{
				TokenEquals,
				TokenNotEquals,
				TokenLess,
				TokenGreater,
				TokenLessEqual,
				TokenGreaterEqual,
				TokenAnd,
				TokenOr,
				TokenNot,
				TokenEOF,
			},
		},
		{
			name:  "object literal",
			input: `{foo: "bar", num: 42}`,
			expected: []TokenType{
				TokenLeftBrace,
				TokenIdentifier,
				TokenColon,
				TokenString,
				TokenComma,
				TokenIdentifier,
				TokenColon,
				TokenNumber,
				TokenRightBrace,
				TokenEOF,
			},
		},
		{
			name:  "keywords",
			input: `if else loop in assert wait set print define run true false null`,
			expected: []TokenType{
				TokenIf,
				TokenElse,
				TokenLoop,
				TokenIn,
				TokenAssert,
				TokenWait,
				TokenSet,
				TokenPrint,
				TokenDefine,
				TokenRun,
				TokenTrue,
				TokenFalse,
				TokenNull,
				TokenEOF,
			},
		},
		{
			name:  "comments",
			input: "# This is a comment\nset x = 1",
			expected: []TokenType{
				TokenNewline,
				TokenSet,
				TokenIdentifier,
				TokenAssign,
				TokenNumber,
				TokenEOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()
			require.NoError(t, err)

			tokenTypes := make([]TokenType, len(tokens))
			for i, token := range tokens {
				tokenTypes[i] = token.Type
			}

			assert.Equal(t, tt.expected, tokenTypes)
		})
	}
}

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, script *ast.Script)
		wantErr bool
	}{
		{
			name:  "connect statement",
			input: `connect "/path/to/server" arg1 arg2`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				conn, ok := script.Statements[0].(*ast.ConnectStatement)
				require.True(t, ok)

				// Check server
				server, ok := conn.Server.(*ast.StringLiteral)
				require.True(t, ok)
				assert.Equal(t, "/path/to/server", server.Value)

				// Check args
				require.Len(t, conn.Args, 2)
			},
		},
		{
			name:  "call statement with result",
			input: `call navigate {url: "https://example.com"} -> page`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				call, ok := script.Statements[0].(*ast.CallStatement)
				require.True(t, ok)

				assert.Equal(t, "navigate", call.Tool)
				assert.Equal(t, "page", call.Variable)

				// Check arguments
				obj, ok := call.Arguments.(*ast.ObjectLiteral)
				require.True(t, ok)
				url, ok := obj.Fields["url"].(*ast.StringLiteral)
				require.True(t, ok)
				assert.Equal(t, "https://example.com", url.Value)
			},
		},
		{
			name:  "assert statement",
			input: `assert result.status == "ok", "Expected status to be ok"`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				assertStmt, ok := script.Statements[0].(*ast.AssertStatement)
				require.True(t, ok)

				assert.Equal(t, "Expected status to be ok", assertStmt.Message)

				// Check expression
				binOp, ok := assertStmt.Expression.(*ast.BinaryOp)
				require.True(t, ok)
				assert.Equal(t, "==", binOp.Operator)
			},
		},
		{
			name: "if-else statement",
			input: `if x > 0 {
				print "positive"
			} else {
				print "negative"
			}`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				ifStmt, ok := script.Statements[0].(*ast.IfStatement)
				require.True(t, ok)

				// Check condition
				cond, ok := ifStmt.Condition.(*ast.BinaryOp)
				require.True(t, ok)
				assert.Equal(t, ">", cond.Operator)

				// Check branches
				require.Len(t, ifStmt.Then, 1)
				require.Len(t, ifStmt.Else, 1)
			},
		},
		{
			name: "loop statement",
			input: `loop item in [1, 2, 3] {
				print item
			}`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				loop, ok := script.Statements[0].(*ast.LoopStatement)
				require.True(t, ok)

				assert.Equal(t, "item", loop.Iterator)

				// Check collection
				arr, ok := loop.Collection.(*ast.ArrayLiteral)
				require.True(t, ok)
				assert.Len(t, arr.Elements, 3)

				// Check body
				require.Len(t, loop.Body, 1)
			},
		},
		{
			name: "define and run",
			input: `define test_login(username, password) {
				call login {user: username, pass: password}
			}
			run test_login("admin", "secret")`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 2)

				// Check define
				def, ok := script.Statements[0].(*ast.DefineStatement)
				require.True(t, ok)
				assert.Equal(t, "test_login", def.Name)
				assert.Equal(t, []string{"username", "password"}, def.Parameters)
				require.Len(t, def.Body, 1)

				// Check run
				run, ok := script.Statements[1].(*ast.RunStatement)
				require.True(t, ok)
				assert.Equal(t, "test_login", run.Name)
				assert.Len(t, run.Arguments, 2)
			},
		},
		{
			name:  "complex expression",
			input: `set result = (a + b) * c / d[0].field`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 1)
				set, ok := script.Statements[0].(*ast.SetStatement)
				require.True(t, ok)
				assert.Equal(t, "result", set.Variable)

				// Value should be a complex expression tree
				assert.NotNil(t, set.Value)
			},
		},
		{
			name: "function calls",
			input: `set length = len(items)
			print json(result)`,
			check: func(t *testing.T, script *ast.Script) {
				require.Len(t, script.Statements, 2)

				// Check len() call
				set, ok := script.Statements[0].(*ast.SetStatement)
				require.True(t, ok)
				fn, ok := set.Value.(*ast.FunctionCall)
				require.True(t, ok)
				assert.Equal(t, "len", fn.Name)
				assert.Len(t, fn.Arguments, 1)

				// Check json() call
				print, ok := script.Statements[1].(*ast.PrintStatement)
				require.True(t, ok)
				fn2, ok := print.Expression.(*ast.FunctionCall)
				require.True(t, ok)
				assert.Equal(t, "json", fn2.Name)
			},
		},
		{
			name:    "syntax error - missing brace",
			input:   `if true { print "yes"`,
			wantErr: true,
		},
		{
			name:    "syntax error - invalid token",
			input:   `set x = @invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens, err := lexer.Tokenize()

			if tt.wantErr && err != nil {
				// Lexer error is acceptable for syntax error tests
				return
			}

			require.NoError(t, err)

			parser := NewParser(tokens)
			script, err := parser.Parse()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.check(t, script)
			}
		})
	}
}

func TestLexerEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "only whitespace",
			input: "   \t\n\r   ",
		},
		{
			name:  "only comments",
			input: "# comment 1\n# comment 2",
		},
		{
			name:  "string with escapes",
			input: `"hello\nworld\t\"quoted\""`,
		},
		{
			name:  "decimal numbers",
			input: "3.14 0.5 42.0",
		},
		{
			name:    "unterminated string",
			input:   `"unterminated`,
			wantErr: true,
		},
		{
			name:    "invalid character",
			input:   `@invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			_, err := lexer.Tokenize()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
