package main

import (
	"fmt"
	"os"

	"cockatoo/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./cockatoo --query \"sql query\" [--debug-ast]")
		os.Exit(1)
	}

	var query string
	var debugAst bool

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--query" && i+1 < len(os.Args) {
			query = os.Args[i+1]
			i++
		} else if os.Args[i] == "--debug-ast" {
			debugAst = true
		} else if i > 1 && os.Args[i-1] == "--query" {
			continue
		} else {
			fmt.Fprintf(os.Stderr, "error: unknown arg %s\n", os.Args[i])
			fmt.Println("usage: ./program --query \"sql query\" [--debug-ast]")
			os.Exit(1)
		}
	}

	if query == "" {
		fmt.Fprintf(os.Stderr, "error: --query is required\n")
		fmt.Println("usage: ./program --query \"sql query\" [--debug-ast]")
		os.Exit(1)
	}

	tree, err := parser.QueryToAst(query)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("sql:", query)
	if debugAst {
		fmt.Println("ast:")
		fmt.Println(tree.String())
	} else {
		fmt.Println("query parsed ok")
	}
}
