package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/go-sqllexer"

	"cockatoo/ast"
)

func parseInsertStatement(ts *TokenStream) (*ast.InsertStmt, error) {
	if err := ts.Consume(T_INSERT); err != nil {
		return nil, err
	}

	if err := ts.Consume(T_INTO); err != nil {
		return nil, err
	}

	tableName, err := ts.ConsumeIdentifier()
	if err != nil {
		return nil, fmt.Errorf("%w: expected table name", ErrSyntaxError)
	}

	result := &ast.InsertStmt{
		TableName: tableName,
	}

	if err := ts.Consume(T_VALUES); err != nil {
		return nil, err
	}

	values, err := parseValuesList(ts)
	if err != nil {
		return nil, err
	}
	result.Values = values

	_, val := ts.Current()
	if !ts.IsEOF() && strings.ToUpper(val) != T_SEMICOLON {
		return nil, fmt.Errorf("%w: unexpected token after INSERT statement", ErrSyntaxError)
	}

	return result, nil
}

func parseValuesList(ts *TokenStream) ([]ast.Expr, error) {
	var values []ast.Expr

	if err := ts.Consume(T_LPAREN); err != nil {
		return nil, err
	}

	for {
		value, err := parseValue(ts)
		if err != nil {
			return nil, err
		}
		values = append(values, value)

		if _, val := ts.Current(); val == T_COMMA {
			ts.Next()
			continue
		}
		break
	}

	if err := ts.Consume(T_RPAREN); err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("%w: INSERT statement must have at least one value", ErrSyntaxError)
	}

	return values, nil
}

func parseValue(ts *TokenStream) (ast.Expr, error) {
	tokenType, val := ts.Current()

	switch tokenType {
	case sqllexer.NUMBER:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid integer", ErrSyntaxError)
		}
		ts.Next()
		return &ast.LiteralInt{Value: intVal}, nil
	case sqllexer.STRING:
		ts.Next()
		return &ast.LiteralString{Value: val}, nil
	case sqllexer.NULL:
		ts.Next()
		return &ast.LiteralNull{}, nil
	default:
		return nil, fmt.Errorf("%w: expected value", ErrSyntaxError)
	}
}
