package main

import (
	"fmt"
	"os"

	"github.com/pelovett/calc/lexer"
	"github.com/pelovett/calc/parser"
)

func main() {
	// Tokenization
	tokens, err := lexer.Tokenize(os.Args[1])
	if err != nil {
		fmt.Printf("Encountered error tokenizing: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Tokens:\n")
	for _, tok := range tokens {
		if tok.Kind() == lexer.Number {
			fmt.Printf("%f ", tok.NumVal())
		} else if tok.Kind() == lexer.Operator {
			fmt.Printf("%s ", tok.OpVal())
		}
	}
	fmt.Println()

	// Parser
	nodeTree, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("Encountered error parsing: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parse Tree:\n%v\n", parser.OpTreeToString(nodeTree))
	fmt.Printf("Result: %f\n", nodeTree.Resolve())
}
