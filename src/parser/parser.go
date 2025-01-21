package parser

import (
	"fmt"
	"os"
	"slices"

	"github.com/pelovett/calc/lexer"
)

type OpTree interface {
	kind() rune
	val() string
	Resolve() (float64, error)
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

func (o NodeOp) Resolve() (float64, error) {
	leftResolve, err := o.LeftArg.Resolve()
	if err != nil {
		return 0, err
	}
	rightResolve, err := o.RightArg.Resolve()
	if err != nil {
		return 0, err
	}
	if o.Op == "+" {
		return leftResolve + rightResolve, nil
	} else if o.Op == "-" {
		return leftResolve - rightResolve, nil
	} else if o.Op == "*" {
		return leftResolve * rightResolve, nil
	} else if o.Op == "/" {
		if rightResolve == 0. {
			return 0, fmt.Errorf("division by zero error")
		}
		return leftResolve / rightResolve, nil
	} else {
		fmt.Printf("Unknown operator: %s\n", o.Op)
		os.Exit(1)
		return 0, fmt.Errorf("Unknown operator: %s\n", o.Op)
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
	Leaf float64
}

const LeafOpKind = 'l'

func (o LeafOp) kind() rune {
	return LeafOpKind
}

func (o LeafOp) val() string {
	return fmt.Sprintf("%f", o.Leaf)
}

func (o *LeafOp) Resolve() (float64, error) {
	return o.Leaf, nil
}

func (o *LeafOp) leftArg() OpTree {
	return nil
}

func (o *LeafOp) rightArg() OpTree {
	return nil
}

func (o *LeafOp) setLeftArg(new OpTree) {}

func (o *LeafOp) setRightArg(new OpTree) {}

func higherPrecedence(first string, second string) bool {
	higher := []string{"*", "/"}
	lesser := []string{"+", "-"}
	return (slices.Contains(higher, first) &&
		slices.Contains(lesser, second))
}

func Parse(tokens []lexer.Token) (OpTree, error) {
	var root OpTree
	var curNode OpTree
	for _, token := range tokens {
		if token.Kind() == lexer.Number {
			if curNode == nil {
				curNode = &LeafOp{Leaf: token.NumVal()}
				root = curNode
			} else if curNode.kind() == NodeOpKind {
				if curNode.rightArg() == nil {
					curNode.setRightArg(&LeafOp{Leaf: token.NumVal()})
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
			} else if curNode.kind() == NodeOpKind && curNode.rightArg() == nil {
				// Can't have two operators in a row
				return nil, fmt.Errorf("syntax error: can't have two operators: %v", token)
			} else {
				// If prev is a leaf, create new node
				if curNode.kind() == LeafOpKind {
					temp := &NodeOp{LeftArg: curNode, Op: token.OpVal()}
					curNode = temp
					root = curNode
					continue
				}
				if higherPrecedence(token.OpVal(), curNode.val()) {
					// If new op is higher precedence, steal the leaf and promote prev
					temp := &NodeOp{Op: token.OpVal()}
					temp.setLeftArg(curNode.rightArg())
					curNode.setRightArg(temp)
					root = curNode
					curNode = temp
				} else {
					// If old op is higher precedence, then attach to root
					temp := &NodeOp{Op: token.OpVal()}
					temp.setLeftArg(root)
					curNode = temp
					root = curNode
				}
			}
		}
	}
	return root, nil
}

func OpTreeToString(tree OpTree) string {
	if tree.kind() == LeafOpKind {
		return tree.val()
	}

	return fmt.Sprintf(
		"%s(%s, %s)",
		tree.val(),
		OpTreeToString(tree.leftArg()),
		OpTreeToString(tree.rightArg()),
	)
}
