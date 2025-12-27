package parser

import (
	"fmt"
	"strings"

	"github.com/DataDog/go-sqllexer"

	"cockatoo/ast"
)

type TokenStream struct {
	lexer       *sqllexer.Lexer
	currentType sqllexer.TokenType
	currentVal  string
	query       string
	initialized bool
	atEOF       bool
}

func NewTokenStream(query string) *TokenStream {
	lexer := sqllexer.New(query)
	return &TokenStream{
		lexer: lexer,
		query: query,
	}
}

func (ts *TokenStream) Initialize() {
	if !ts.initialized {
		// read first token
		t := ts.lexer.Scan()
		if t.Type != sqllexer.EOF {
			ts.currentVal = t.Value
			ts.currentType = t.Type
		} else {
			ts.currentVal = ""
			ts.atEOF = true
		}

		ts.initialized = true
	}
}

func (ts *TokenStream) Next() (sqllexer.TokenType, string) {
	if !ts.initialized {
		ts.Initialize()
		return ts.currentType, ts.currentVal
	}

	if ts.atEOF {
		return ts.currentType, ts.currentVal
	}

	t := ts.lexer.Scan()
	for t.Type == sqllexer.SPACE {
		t = ts.lexer.Scan()
		continue
	}

	ts.currentType = t.Type
	ts.currentVal = t.Value

	if ts.currentType == sqllexer.EOF {
		ts.atEOF = true
	}

	return ts.currentType, ts.currentVal
}

func (ts *TokenStream) Current() (sqllexer.TokenType, string) {
	if !ts.initialized {
		ts.Initialize()
	}
	return ts.currentType, ts.currentVal
}

func (ts *TokenStream) Consume(expected string) error {
	_, val := ts.Current()
	upperVal := strings.ToUpper(val)
	upperExpected := strings.ToUpper(expected)

	if upperVal != upperExpected {
		return fmt.Errorf("%w: expected %s, got %q",
			ErrSyntaxError, expected, val)
	}

	ts.Next()
	return nil
}

func (ts *TokenStream) ConsumeIdentifier() (string, error) {
	tokenType, val := ts.Current()
	if tokenType != sqllexer.IDENT {
		return "", fmt.Errorf("%w: expected identifier, got %q", ErrSyntaxError, val)
	}

	identifier := val
	ts.Next()
	return identifier, nil
}

func (ts *TokenStream) ConsumeNumber() (string, error) {
	tokenType, val := ts.Current()
	if tokenType != sqllexer.NUMBER {
		return "", fmt.Errorf("%w: expected number, got %q", ErrSyntaxError, val)
	}

	number := val
	ts.Next()
	return number, nil
}

func (ts *TokenStream) ConsumeString() (string, error) {
	tokenType, val := ts.Current()
	if tokenType != sqllexer.STRING {
		return "", fmt.Errorf("%w: expected string, got %q", ErrSyntaxError, val)
	}

	str := val
	ts.Next()
	return str, nil
}

func (ts *TokenStream) IsEOF() bool {
	tokenType, _ := ts.Current()
	return tokenType == sqllexer.EOF
}

func ParseQuery(query string) (ast.Wrapper, error) {
	ts := NewTokenStream(query)
	ts.Initialize()

	_, val := ts.Current()
	upperVal := strings.ToUpper(val)

	// Determine the statement type based on the first token
	if upperVal == T_SELECT {
		return parseSelectStatement(ts)
	} else if upperVal == T_CREATE {
		return parseCreateTableStatement(ts)
	} else if upperVal == T_INSERT {
		return parseInsertStatement(ts)
	} else {
		return nil, fmt.Errorf("%w: unsupported statement type: %q", ErrSyntaxError, val)
	}
}

func QueryToAst(query string) (ast.Wrapper, error) {
	return ParseQuery(query)
}
