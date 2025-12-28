package parser

import (
	"reflect"
	"testing"

	"cockatoo/ast"
)

func TestSelectQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected ast.Wrapper
		wantErr  bool
	}{
		{
			name:  "simple select all",
			query: "SELECT * FROM users",
			expected: &ast.SelectStmt{
				Projections: []ast.ProjectionItem{
					{IsWildcard: true},
				},
				From: ast.TableRef{Name: "users"},
			},
			wantErr: false,
		},
		{
			name:  "select specific columns",
			query: "SELECT id, name FROM users",
			expected: &ast.SelectStmt{
				Projections: []ast.ProjectionItem{
					{Expression: &ast.ColumnRef{Name: "id"}, IsWildcard: false},
					{Expression: &ast.ColumnRef{Name: "name"}, IsWildcard: false},
				},
				From: ast.TableRef{Name: "users"},
			},
			wantErr: false,
		},
		{
			name:  "select with where clause",
			query: "SELECT id FROM users WHERE age > 18",
			expected: &ast.SelectStmt{
				Projections: []ast.ProjectionItem{
					{Expression: &ast.ColumnRef{Name: "id"}, IsWildcard: false},
				},
				From: ast.TableRef{Name: "users"},
				Selection: &ast.ComparisonOp{
					Left:     &ast.ColumnRef{Name: "age"},
					Operator: ">",
					Right:    &ast.LiteralInt{Value: 18},
				},
			},
			wantErr: false,
		},
		{
			name:  "select with where and limit",
			query: "SELECT id FROM users WHERE age > 18 LIMIT 10",
			expected: &ast.SelectStmt{
				Projections: []ast.ProjectionItem{
					{Expression: &ast.ColumnRef{Name: "id"}, IsWildcard: false},
				},
				From: ast.TableRef{Name: "users"},
				Selection: &ast.ComparisonOp{
					Left:     &ast.ColumnRef{Name: "age"},
					Operator: ">",
					Right:    &ast.LiteralInt{Value: 18},
				},
				Limit: &[]uint64{10}[0],
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := QueryToAst(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryToAst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Check the type first
			if reflect.TypeOf(result) != reflect.TypeOf(tt.expected) {
				t.Errorf("QueryToAst() result type = %T, expected type %T", result, tt.expected)
				return
			}

			// Convert to JSON strings for comparison
			resultJSON := result.String()
			expectedJSON := tt.expected.String()

			// Compare JSON strings
			if resultJSON != expectedJSON {
				t.Errorf("QueryToAst() result = %s\nexpected = %s", resultJSON, expectedJSON)
			}
		})
	}
}

func TestCreateTableQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected ast.Wrapper
		wantErr  bool
	}{
		{
			name:  "simple create table",
			query: "CREATE TABLE users (id INT, name TEXT)",
			expected: &ast.CreateTableStmt{
				TableName: "users",
				Columns: []ast.ColumnDef{
					{Name: "id", Type: "INT"},
					{Name: "name", Type: "TEXT"},
				},
			},
			wantErr: false,
		},
		{
			name:  "create table with multiple columns",
			query: "CREATE TABLE products (id INT, name TEXT, price BIGINT, description TEXT)",
			expected: &ast.CreateTableStmt{
				TableName: "products",
				Columns: []ast.ColumnDef{
					{Name: "id", Type: "INT"},
					{Name: "name", Type: "TEXT"},
					{Name: "price", Type: "BIGINT"},
					{Name: "description", Type: "TEXT"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := QueryToAst(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryToAst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Check the type first
			if reflect.TypeOf(result) != reflect.TypeOf(tt.expected) {
				t.Errorf("QueryToAst() result type = %T, expected type %T", result, tt.expected)
				return
			}

			// Convert to JSON strings for comparison
			resultJSON := result.String()
			expectedJSON := tt.expected.String()

			// Compare JSON strings
			if resultJSON != expectedJSON {
				t.Errorf("QueryToAst() result = %s\nexpected = %s", resultJSON, expectedJSON)
			}
		})
	}
}

func TestInsertQueries(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected ast.Wrapper
		wantErr  bool
	}{
		{
			name:  "simple insert",
			query: "INSERT INTO users VALUES (1, 'Alice')",
			expected: &ast.InsertStmt{
				TableName: "users",
				Values: []ast.Expr{
					&ast.LiteralInt{Value: 1},
					&ast.LiteralString{Value: "Alice"},
				},
			},
			wantErr: false,
		},
		{
			name:  "insert with multiple integer values",
			query: "INSERT INTO numbers VALUES (1, 2, 3)",
			expected: &ast.InsertStmt{
				TableName: "numbers",
				Values: []ast.Expr{
					&ast.LiteralInt{Value: 1},
					&ast.LiteralInt{Value: 2},
					&ast.LiteralInt{Value: 3},
				},
			},
			wantErr: false,
		},
		{
			name:  "insert with mixed values",
			query: "INSERT INTO products VALUES (1, 'Laptop', 1200, 'High performance laptop')",
			expected: &ast.InsertStmt{
				TableName: "products",
				Values: []ast.Expr{
					&ast.LiteralInt{Value: 1},
					&ast.LiteralString{Value: "Laptop"},
					&ast.LiteralInt{Value: 1200},
					&ast.LiteralString{Value: "High performance laptop"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := QueryToAst(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryToAst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// Check the type first
			if reflect.TypeOf(result) != reflect.TypeOf(tt.expected) {
				t.Errorf("QueryToAst() result type = %T, expected type %T", result, tt.expected)
				return
			}

			// Convert to JSON strings for comparison
			resultJSON := result.String()
			expectedJSON := tt.expected.String()

			// Compare JSON strings
			if resultJSON != expectedJSON {
				t.Errorf("QueryToAst() result = %s\nexpected = %s", resultJSON, expectedJSON)
			}
		})
	}
}
