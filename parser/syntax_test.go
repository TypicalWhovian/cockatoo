package parser

import (
	"errors"
	"testing"
)

func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectedErr error
	}{
		{
			name:        "missing column list in SELECT",
			query:       "SELECT FROM users",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "missing values in INSERT",
			query:       "INSERT INTO users",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "empty column list in CREATE TABLE",
			query:       "CREATE TABLE t()",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "invalid SQL syntax",
			query:       "SELEC name FROM users",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "missing FROM clause",
			query:       "SELECT name",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "incomplete WHERE clause",
			query:       "SELECT name FROM users WHERE",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "unmatched parentheses",
			query:       "SELECT name FROM users WHERE (age > 18",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "invalid comparison operator",
			query:       "SELECT name FROM users WHERE age <=> 18",
			expectedErr: ErrSyntaxError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := QueryToAst(tt.query)
			if err == nil {
				t.Errorf("Expected error for query: %s, but got none", tt.query)
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %q, but got %q", tt.expectedErr, err)
			}
		})
	}
}

// Test for basic structural validation
func TestStructuralValidation(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectedErr error
	}{
		{
			name:        "SELECT without columns",
			query:       "SELECT FROM users",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "INSERT without values",
			query:       "INSERT INTO users VALUES",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "CREATE TABLE without columns",
			query:       "CREATE TABLE t()",
			expectedErr: ErrSyntaxError,
		},
		{
			name:        "WHERE without expression",
			query:       "SELECT * FROM users WHERE",
			expectedErr: ErrSyntaxError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := QueryToAst(tt.query)
			if err == nil {
				t.Errorf("Expected error for query: %s, but got none", tt.query)
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %q, but got %q", tt.expectedErr, err)
			}
		})
	}
}
