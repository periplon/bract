package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Literals
	TokenEOF TokenType = iota
	TokenString
	TokenNumber
	TokenBoolean
	TokenIdentifier

	// Keywords
	TokenConnect
	TokenCall
	TokenAssert
	TokenWait
	TokenLoop
	TokenIn
	TokenIf
	TokenElse
	TokenSet
	TokenPrint
	TokenDefine
	TokenRun
	TokenTrue
	TokenFalse
	TokenNull

	// Operators
	TokenEquals
	TokenNotEquals
	TokenLess
	TokenGreater
	TokenLessEqual
	TokenGreaterEqual
	TokenAnd
	TokenOr
	TokenNot
	TokenAssign
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide

	// Delimiters
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenLeftBracket
	TokenRightBracket
	TokenComma
	TokenColon
	TokenDot
	TokenArrow
	TokenNewline

	// Comments
	TokenComment
)

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Value    string
	Line     int
	Column   int
}

// Lexer tokenizes DSL input
type Lexer struct {
	input   string
	pos     int
	line    int
	column  int
	tokens  []Token
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		column: 1,
	}
}

// Tokenize converts the input into tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		l.skipWhitespace()
		if l.pos >= len(l.input) {
			break
		}

		if err := l.nextToken(); err != nil {
			return nil, err
		}
	}

	l.tokens = append(l.tokens, Token{
		Type:   TokenEOF,
		Line:   l.line,
		Column: l.column,
	})

	return l.tokens, nil
}

func (l *Lexer) nextToken() error {
	ch := l.input[l.pos]

	// Comments
	if ch == '#' {
		l.skipComment()
		return nil
	}

	// String literals
	if ch == '"' || ch == '\'' {
		return l.readString(ch)
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) {
		return l.readNumber()
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readIdentifier()
	}

	// Operators and delimiters
	switch ch {
	case '=':
		if l.peek() == '=' {
			l.addToken(TokenEquals, "==")
			l.advance()
			l.advance()
		} else {
			l.addToken(TokenAssign, "=")
			l.advance()
		}
	case '!':
		if l.peek() == '=' {
			l.addToken(TokenNotEquals, "!=")
			l.advance()
			l.advance()
		} else {
			l.addToken(TokenNot, "!")
			l.advance()
		}
	case '<':
		if l.peek() == '=' {
			l.addToken(TokenLessEqual, "<=")
			l.advance()
			l.advance()
		} else {
			l.addToken(TokenLess, "<")
			l.advance()
		}
	case '>':
		if l.peek() == '=' {
			l.addToken(TokenGreaterEqual, ">=")
			l.advance()
			l.advance()
		} else {
			l.addToken(TokenGreater, ">")
			l.advance()
		}
	case '&':
		if l.peek() == '&' {
			l.addToken(TokenAnd, "&&")
			l.advance()
			l.advance()
		} else {
			return fmt.Errorf("unexpected character '&' at line %d, column %d", l.line, l.column)
		}
	case '|':
		if l.peek() == '|' {
			l.addToken(TokenOr, "||")
			l.advance()
			l.advance()
		} else {
			return fmt.Errorf("unexpected character '|' at line %d, column %d", l.line, l.column)
		}
	case '-':
		if l.peek() == '>' {
			l.addToken(TokenArrow, "->")
			l.advance()
			l.advance()
		} else {
			l.addToken(TokenMinus, "-")
			l.advance()
		}
	case '+':
		l.addToken(TokenPlus, "+")
		l.advance()
	case '*':
		l.addToken(TokenMultiply, "*")
		l.advance()
	case '/':
		l.addToken(TokenDivide, "/")
		l.advance()
	case '(':
		l.addToken(TokenLeftParen, "(")
		l.advance()
	case ')':
		l.addToken(TokenRightParen, ")")
		l.advance()
	case '{':
		l.addToken(TokenLeftBrace, "{")
		l.advance()
	case '}':
		l.addToken(TokenRightBrace, "}")
		l.advance()
	case '[':
		l.addToken(TokenLeftBracket, "[")
		l.advance()
	case ']':
		l.addToken(TokenRightBracket, "]")
		l.advance()
	case ',':
		l.addToken(TokenComma, ",")
		l.advance()
	case ':':
		l.addToken(TokenColon, ":")
		l.advance()
	case '.':
		l.addToken(TokenDot, ".")
		l.advance()
	case '\n':
		l.addToken(TokenNewline, "\n")
		l.advance()
		l.line++
		l.column = 0
	default:
		return fmt.Errorf("unexpected character '%c' at line %d, column %d", ch, l.line, l.column)
	}

	return nil
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) skipComment() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.advance()
	}
}

func (l *Lexer) readString(quote byte) error {
	start := l.pos
	l.advance() // Skip opening quote

	var value strings.Builder
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		ch := l.input[l.pos]
		if ch == '\\' && l.pos+1 < len(l.input) {
			l.advance()
			switch l.input[l.pos] {
			case 'n':
				value.WriteByte('\n')
			case 't':
				value.WriteByte('\t')
			case 'r':
				value.WriteByte('\r')
			case '\\':
				value.WriteByte('\\')
			case '"':
				value.WriteByte('"')
			case '\'':
				value.WriteByte('\'')
			default:
				value.WriteByte(l.input[l.pos])
			}
		} else {
			value.WriteByte(ch)
		}
		l.advance()
	}

	if l.pos >= len(l.input) {
		return fmt.Errorf("unterminated string at line %d, column %d", l.line, start)
	}

	l.advance() // Skip closing quote
	l.addToken(TokenString, value.String())
	return nil
}

func (l *Lexer) readNumber() error {
	start := l.pos
	hasDecimal := false

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsDigit(rune(ch)) {
			l.advance()
		} else if ch == '.' && !hasDecimal && l.pos+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.pos+1])) {
			hasDecimal = true
			l.advance()
		} else {
			break
		}
	}

	l.addToken(TokenNumber, l.input[start:l.pos])
	return nil
}

func (l *Lexer) readIdentifier() error {
	start := l.pos

	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	value := l.input[start:l.pos]
	tokenType := l.getKeywordType(value)
	l.addToken(tokenType, value)
	return nil
}

func (l *Lexer) getKeywordType(value string) TokenType {
	keywords := map[string]TokenType{
		"connect": TokenConnect,
		"call":    TokenCall,
		"assert":  TokenAssert,
		"wait":    TokenWait,
		"loop":    TokenLoop,
		"in":      TokenIn,
		"if":      TokenIf,
		"else":    TokenElse,
		"set":     TokenSet,
		"print":   TokenPrint,
		"define":  TokenDefine,
		"run":     TokenRun,
		"true":    TokenTrue,
		"false":   TokenFalse,
		"null":    TokenNull,
	}

	if tokenType, ok := keywords[strings.ToLower(value)]; ok {
		return tokenType
	}
	return TokenIdentifier
}

func (l *Lexer) advance() {
	l.pos++
	l.column++
}

func (l *Lexer) peek() byte {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

func (l *Lexer) addToken(tokenType TokenType, value string) {
	l.tokens = append(l.tokens, Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	})
}