package parser

import (
	"fmt"
	"os"

	"github.com/pelovett/calc/lexer"
)

type OpTree interface {
	kind() rune
	val() string
	resolve() float64
	leftArg() OpTree
	rightArg() OpTree
	setLeftArg(OpTree)
	setRightArg(OpTree)
}

type NodeOp struct {
	LeftArg  OpTree
	RightArg OpTree
	Op       string
}

const NodeOpKind = 'n'

func (o NodeOp) kind() rune {
	return NodeOpKind
}

func (o NodeOp) val() string {
	return o.Op
}

func (o NodeOp) resolve() float64 {
	if o.Op == "+" {
		return o.LeftArg.resolve() + o.RightArg.resolve()
	} else if o.Op == "-" {
		return o.LeftArg.resolve() - o.RightArg.resolve()
	} else {
		fmt.Printf("Unknown operator: %s\n", o.Op)
		os.Exit(1)
		return 0
	}
}

func (o *NodeOp) leftArg() OpTree {
	return o.LeftArg
}

func (o *NodeOp) rightArg() OpTree {
	return o.RightArg
}

func (o *NodeOp) setLeftArg(new OpTree) {
	o.LeftArg = new
}

func (o *NodeOp) setRightArg(new OpTree) {
	o.RightArg = new
}

type LeafOp struct {
	leaf float64
}

const LeafOpKind = 'l'

func (o LeafOp) kind() rune {
	return LeafOpKind
}

func (o LeafOp) val() string {
	return fmt.Sprintf("%f", o.leaf)
}

func (o *LeafOp) resolve() float64 {
	return o.leaf
}

func (o *LeafOp) leftArg() OpTree {
	return nil
}

func (o *LeafOp) rightArg() OpTree {
	return nil
}

func (o *LeafOp) setLeftArg(new OpTree) {}

func (o *LeafOp) setRightArg(new OpTree) {}

func Parse(tokens []lexer.Token) (OpTree, error) {
	var curNode OpTree
	for _, token := range tokens {
		if token.Kind() == lexer.Number {
			if curNode == nil {
				curNode = &LeafOp{leaf: token.NumVal()}
			} else if curNode.kind() == NodeOpKind {
				if curNode.rightArg() == nil {
					curNode.setRightArg(&LeafOp{leaf: token.NumVal()})
				} else {
					return nil, fmt.Errorf("syntax error two values to right of operator: %v", token)
				}
			} else {
				// Trying to parse two numbers in a row, error
				return nil, fmt.Errorf("syntax error, two numbers without opp between: %v", token)
			}
		} else if token.Kind() == lexer.Operator {
			if curNode == nil {
				// Can't start with operator
				return nil, fmt.Errorf("syntax error: can't start with operator: %v", token)
			} else if curNode.kind() != LeafOpKind && curNode.rightArg() == nil {
				// Can't have two operators in a row
				return nil, fmt.Errorf("syntax error: can't have two operators: %v", token)
			} else {
				temp := &NodeOp{LeftArg: curNode, Op: token.OpVal()}
				curNode = temp
			}
		}
	}
	return curNode, nil
}

func OpTreeToString(tree OpTree) string {
	if tree.kind() == LeafOpKind {
		return tree.val()
	}

	return fmt.Sprintf("%s(%s, %s)",
		tree.val(),
		OpTreeToString(tree.leftArg()),
		OpTreeToString(tree.rightArg()))

}
