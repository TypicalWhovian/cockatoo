package parser

import (
	"fmt"
	"strings"

	"cockatoo/ast"
)

func parseCreateTableStatement(ts *TokenStream) (*ast.CreateTableStmt, error) {
	if err := ts.Consume(T_CREATE); err != nil {
		return nil, err
	}

	if err := ts.Consume(T_TABLE); err != nil {
		return nil, err
	}

	tableName, err := ts.ConsumeIdentifier()
	if err != nil {
		return nil, fmt.Errorf("%w: expected table name", ErrSyntaxError)
	}

	result := &ast.CreateTableStmt{
		TableName: tableName,
	}

	if err := ts.Consume(T_LPAREN); err != nil {
		return nil, err
	}

	columns, err := parseColumnDefinitions(ts)
	if err != nil {
		return nil, err
	}
	result.Columns = columns

	if err := ts.Consume(T_RPAREN); err != nil {
		return nil, err
	}

	_, val := ts.Current()
	if !ts.IsEOF() && strings.ToUpper(val) != T_SEMICOLON {
		return nil, fmt.Errorf("%w: unexpected token after CREATE TABLE statement", ErrSyntaxError)
	}

	return result, nil
}

func parseColumnDefinitions(ts *TokenStream) ([]ast.ColumnDef, error) {
	var columns []ast.ColumnDef

	for {
		columnName, err := ts.ConsumeIdentifier()
		if err != nil {
			return nil, fmt.Errorf("%w: expected column name", ErrSyntaxError)
		}

		columnType, err := ts.ConsumeIdentifier()
		if err != nil {
			return nil, fmt.Errorf("%w: expected column type", ErrSyntaxError)
		}

		columnType = strings.ToUpper(columnType)

		if _, ok := validTypes[columnType]; !ok {
			return nil, fmt.Errorf("%w: unsupported column type %s", ErrSyntaxError, columnType)
		}

		columns = append(columns, ast.ColumnDef{
			Name: columnName,
			Type: columnType,
		})

		if _, val := ts.Current(); val == T_COMMA {
			ts.Next()
			continue
		}
		break
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("%w: table must have at least one column", ErrSyntaxError)
	}

	return columns, nil
}
