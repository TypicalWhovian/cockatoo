package parser

import (
	"errors"
)

var (
	ErrFailedToParse = errors.New("failed to parse query")
	ErrSyntaxError   = errors.New("syntax error in sql query")
)

const (
	T_SELECT = "SELECT"
	T_FROM   = "FROM"
	T_WHERE  = "WHERE"
	T_LIMIT  = "LIMIT"
	T_CREATE = "CREATE"
	T_TABLE  = "TABLE"
	T_INSERT = "INSERT"
	T_INTO   = "INTO"
	T_VALUES = "VALUES"
	T_AND    = "AND"
	T_OR     = "OR"

	T_INT    = "INT"
	T_BIGINT = "BIGINT"
	T_TEXT   = "TEXT"

	T_COMMA     = ","
	T_SEMICOLON = ";"
	T_LPAREN    = "("
	T_RPAREN    = ")"
	T_STAR      = "*"
	T_EQ        = "="
	T_NEQ       = "!="
	T_GT        = ">"
	T_GTE       = ">="
	T_LT        = "<"
	T_LTE       = "<="
)

var (
	validOps = map[string]struct{}{
		T_EQ:  {},
		T_NEQ: {},
		T_GT:  {},
		T_GTE: {},
		T_LT:  {},
		T_LTE: {},
	}

	validTypes = map[string]struct{}{
		T_INT:    {},
		T_BIGINT: {},
		T_TEXT:   {},
	}
)
