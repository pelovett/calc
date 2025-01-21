package lexer

import (
	"fmt"
	"os"
	"strconv"
)

const (
	Number   rune = 'N'
	Operator rune = 'O'
)

type Token interface {
	Kind() rune
	OpVal() string // This is bad abstraction
	NumVal() float64
}

type NumToken struct {
	Val float64
}

func (t NumToken) Kind() rune {
	return Number
}

func (t NumToken) OpVal() string {
	fmt.Printf("Do NOT call OpVal on NumToken, exiting...\n")
	os.Exit(1)
	return ""
}

func (t NumToken) NumVal() float64 {
	return t.Val
}

type OpToken struct {
	Val string
}

func (t OpToken) Kind() rune {
	return Operator
}

func (t OpToken) OpVal() string {
	return t.Val
}

func (t OpToken) NumVal() float64 {
	fmt.Printf("Do NOT call NumVal on OpToken, exiting...\n")
	os.Exit(1)
	return 0
}

func createToken(tokenType rune, input string) (Token, error) {
	if tokenType == Number {
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return nil, err
		}
		return NumToken{Val: val}, nil
	} else if tokenType == Operator {
		return OpToken{Val: input}, nil
	} else {
		return nil, fmt.Errorf("unknown token type: %c for input: %s", tokenType, input)
	}
}

func isSomething(something []rune) func(rune) bool {
	return func(r rune) bool {
		for _, val := range something {
			if r == val {
				return true
			}
		}
		return false
	}
}

var isDigit = isSomething([]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'})
var isWhitespace = isSomething([]rune{' ', '\t', '\n', '\r'})
var isOp = isSomething([]rune{'+', '-', '*', '/'})

func Tokenize(input string) ([]Token, error) {
	output := []Token{}

	var curKind rune
	var curStart int
	for i, char := range input {
		// If we're just starting we need to start a digit (TODO support parens)
		switch curKind {
		case 0:
			if isDigit(char) {
				curKind = Number
				curStart = i
			} else if isWhitespace(char) {
				// Skip whitespace for now
				continue
			} else if isOp(char) {
				curKind = Operator
				curStart = i
			} else {
				return nil, fmt.Errorf("illegal starting char at position %d: %c", i+1, char)
			}
		case Number:
			if isDigit(char) {
				continue
			} else if isWhitespace(char) {
				// End of cur Token, emit
				tok, err := createToken(curKind, input[curStart:i])
				if err != nil {
					return nil, err
				}
				output = append(output, tok)
				curKind = 0
				curStart = 0
			} else if isOp(char) {
				// End of cur Token, emit then start op Token
				tok, err := createToken(curKind, input[curStart:i])
				if err != nil {
					return nil, err
				}
				output = append(output, tok)
				curKind = Operator
				curStart = i
			} else {
				return nil, fmt.Errorf("illegal char in number at position %d: %c", i+1, char)
			}
		case Operator:
			if isDigit(char) {
				// End of cur Token, emit then start number Token
				tok, err := createToken(curKind, input[curStart:i])
				if err != nil {
					return nil, err
				}
				output = append(output, tok)
				curKind = Number
				curStart = i
			} else if isWhitespace(char) {
				// End of cur Token, emit
				tok, err := createToken(curKind, input[curStart:i])
				if err != nil {
					return nil, err
				}
				output = append(output, tok)
				curKind = 0
				curStart = 0
			} else if isOp(char) {
				// Ops are currenty only one char wide so error out
				return nil, fmt.Errorf("double operator at position %d: %c", i+1, char)
			} else {
				return nil, fmt.Errorf("illegal char in operator at position %d: %c", i+1, char)
			}
		}
	}

	// Cleanup final Token
	if curKind != 0 {
		tok, err := createToken(curKind, input[curStart:])
		if err != nil {
			return nil, err
		}
		output = append(output, tok)
	}
	return output, nil
}
