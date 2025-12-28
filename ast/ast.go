package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Wrapper interface {
	String() string
}

type SelectStmt struct {
	Projections []ProjectionItem
	From        TableRef
	Selection   Expr
	Limit       *uint64
}

type CreateTableStmt struct {
	TableName string
	Columns   []ColumnDef
}

type InsertStmt struct {
	TableName string
	Values    []Expr
}

type ProjectionItem struct {
	Expression Expr
	IsWildcard bool
}

type TableRef struct {
	Name string
}

type ColumnDef struct {
	Name string
	Type string
}

type Expr interface {
	ExprString() string
}

type ColumnRef struct {
	Name string
}

type LiteralInt struct {
	Value int64
}

type LiteralString struct {
	Value string
}

type LiteralNull struct {
}

type ComparisonOp struct {
	Left     Expr
	Right    Expr
	Operator string
}

type LogicalOp struct {
	Left     Expr
	Right    Expr
	Operator string // and, or
}

func (s *SelectStmt) String() string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	encoder.Encode(s)
	return buffer.String()
}

func (s *CreateTableStmt) String() string {
	res, _ := json.MarshalIndent(s, "", "  ")
	return string(res)
}

func (s *InsertStmt) String() string {
	res, _ := json.MarshalIndent(s, "", "  ")
	return string(res)
}

func (c *ColumnRef) ExprString() string {
	return c.Name
}

func (l *LiteralInt) ExprString() string {
	return fmt.Sprintf("%d", l.Value)
}

func (l *LiteralString) ExprString() string {
	return fmt.Sprintf("'%s'", l.Value)
}

func (b *ComparisonOp) ExprString() string {
	return fmt.Sprintf("%s %s %s", b.Left.ExprString(), b.Operator, b.Right.ExprString())
}

func (l *LogicalOp) ExprString() string {
	return fmt.Sprintf("%s %s %s", l.Left.ExprString(), l.Operator, l.Right.ExprString())
}

func (n *LiteralNull) ExprString() string {
	return "null"
}
