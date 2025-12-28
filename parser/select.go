package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/go-sqllexer"

	"cockatoo/ast"
)

func parseSelectStatement(ts *TokenStream) (*ast.SelectStmt, error) {
	result := &ast.SelectStmt{}

	if err := ts.Consume(T_SELECT); err != nil {
		return nil, err
	}

	projections, err := parseProjectionList(ts)
	if err != nil {
		return nil, err
	}
	result.Projections = projections

	if err := ts.Consume(T_FROM); err != nil {
		return nil, err
	}

	tableName, err := parseTableName(ts)
	if err != nil {
		return nil, err
	}
	result.From = tableName

	for {
		_, val := ts.Current()
		upperVal := strings.ToUpper(val)
		switch upperVal {
		case T_WHERE:
			if result.Limit != nil {
				return nil, fmt.Errorf("%w: WHERE clause must come before LIMIT", ErrSyntaxError)
			}
			ts.Next()

			expr, err := parseExpression(ts)
			if err != nil {
				return nil, err
			}
			result.Selection = expr
		case T_LIMIT:
			ts.Next()

			limitStr, err := ts.ConsumeNumber()
			if err != nil {
				return nil, fmt.Errorf("%w: expected number after LIMIT", ErrSyntaxError)
			}

			limitVal, err := strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: invalid LIMIT value", ErrSyntaxError)
			}
			result.Limit = &limitVal
		case T_SEMICOLON:
			return result, nil
		default:
			if ts.IsEOF() {
				return result, nil
			}
			return nil, fmt.Errorf("%w: unexpected token after SELECT statement", ErrSyntaxError)
		}
	}
}

func parseProjectionList(ts *TokenStream) ([]ast.ProjectionItem, error) {
	var projections []ast.ProjectionItem

	for {
		token, val := ts.Current()

		if val == T_STAR {
			projections = append(projections, ast.ProjectionItem{IsWildcard: true})
			ts.Next()
			return projections, nil
		}

		if token == sqllexer.IDENT {
			columnName := val
			projection := ast.ProjectionItem{
				Expression: &ast.ColumnRef{Name: columnName},
				IsWildcard: false,
			}
			ts.Next()

			token, val = ts.Current()

			projections = append(projections, projection)

			if _, val = ts.Current(); val == T_COMMA {
				ts.Next()
				continue
			}
			break
		}

		return nil, fmt.Errorf("%w: expected column name or *", ErrSyntaxError)
	}

	return projections, nil
}

func parseTableName(ts *TokenStream) (ast.TableRef, error) {
	tableName, err := ts.ConsumeIdentifier()
	if err != nil {
		return ast.TableRef{}, fmt.Errorf("%w: expected table name", ErrSyntaxError)
	}

	result := ast.TableRef{
		Name: tableName,
	}

	return result, nil
}

func parseExpression(ts *TokenStream) (ast.Expr, error) {
	left, err := parseSimpleExpression(ts)
	if err != nil {
		return nil, err
	}

	_, val := ts.Current()
	upperVal := strings.ToUpper(val)
	if upperVal == T_AND || upperVal == T_OR {
		operator := upperVal
		ts.Next() // consume and/or

		right, err := parseExpression(ts)
		if err != nil {
			return nil, err
		}

		return &ast.LogicalOp{
			Left:     left,
			Right:    right,
			Operator: operator,
		}, nil
	}

	return left, nil
}

func parseSimpleExpression(ts *TokenStream) (ast.Expr, error) {
	columnName, err := ts.ConsumeIdentifier()
	if err != nil {
		return nil, fmt.Errorf("%w: expected column name", ErrSyntaxError)
	}

	left := &ast.ColumnRef{Name: columnName}

	_, val := ts.Current()
	if _, ok := validOps[val]; !ok {
		return nil, fmt.Errorf("%w: expected comparison operator (>, <, =, !=, >=, <=) in 'WHERE' clause, got %q", ErrSyntaxError, val)
	}

	operator := val
	ts.Next()

	tokenType, val := ts.Current()
	var right ast.Expr

	switch tokenType {
	case sqllexer.NUMBER:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid integer", ErrSyntaxError)
		}
		right = &ast.LiteralInt{Value: intVal}
		ts.Next()
	case sqllexer.STRING:
		right = &ast.LiteralString{Value: val}
		ts.Next()
	case sqllexer.IDENT:
		right = &ast.ColumnRef{Name: val}
		ts.Next()
	default:
		return nil, fmt.Errorf("%w: expected value", ErrSyntaxError)
	}

	return &ast.ComparisonOp{
		Left:     left,
		Right:    right,
		Operator: operator,
	}, nil
}
