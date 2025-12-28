# Cockatoo SQL Parser

Cockatoo is a SQL parser written in Go that converts SQL queries into an Abstract Syntax Tree (AST). It currently supports parsing SELECT, CREATE TABLE, and INSERT statements.

## Features

- Parse SQL SELECT statements
- Parse SQL CREATE TABLE statements
- Parse SQL INSERT statements
- Convert SQL queries to AST representations
- Display the AST structure for debugging

## Prerequisites

- Go 1.24+ installed
- Git (for cloning the repository)

## Building the Project

To build the Cockatoo binary, run:

```bash
# Build the project
go build cockatoo
```

This will create a binary executable named `cockatoo` in the current directory.

## Running the Project

After building, you can run the executable with the following syntax:

```bash
./cockatoo --query "SQL_QUERY" [--debug-ast]
```

Where:
- `--query` is a required parameter followed by the SQL query to parse
- `--debug-ast` is an optional flag that prints out the full AST structure

### Examples

```bash
# Parse a simple SELECT query
./cockatoo --query "SELECT * FROM users"

# Parse a CREATE TABLE query with debug AST output
./cockatoo --query "CREATE TABLE users (id INT, name TEXT)" --debug-ast

# Parse an INSERT query with debug AST output
./cockatoo --query "INSERT INTO users VALUES (1, 'John')" --debug-ast
```

## Running Tests

The project includes tests in the parser package. To run tests:

```bash
# Run all tests in the project
go test ./...

# Run tests in the parser package only
go test ./parser

# Run a specific test function
go test -v ./parser -run TestSelectQueries

# Run tests with verbose output
go test -v ./...
```

## Project Structure

```
/
├── ast/
│   └── ast.go             # Contains AST node definitions for SQL syntax
├── parser/
│   ├── constants.go       # SQL language constants
│   ├── create.go          # Parser for CREATE TABLE statements
│   ├── insert.go          # Parser for INSERT statements
│   ├── lexer.go           # SQL lexer and token stream handling
│   ├── parser_test.go     # Test cases for parsing different SQL statements
│   ├── select.go          # Parser for SELECT statements
│   └── syntax_test.go     # Additional syntax tests
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── main.go                # Main application entry point
```

## Dependencies

- [github.com/DataDog/go-sqllexer](https://github.com/DataDog/go-sqllexer) v0.1.10 (for SQL lexical analysis)