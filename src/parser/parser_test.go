package parser

import (
	"math"
	"testing"

	"github.com/pelovett/calc/lexer"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func equalOpTree(a, b OpTree) bool {
	aStack, bStack := []OpTree{a}, []OpTree{b}

	for len(aStack) > 0 {
		curA := aStack[len(aStack)-1]
		curB := bStack[len(bStack)-1]
		aStack = aStack[:len(aStack)-1]
		bStack = bStack[:len(bStack)-1]

		if (curA == nil && curB != nil) || (curA != nil && curB == nil) {
			return false
		} else if curA.kind() != curB.kind() {
			return false
		} else if curA.val() != curB.val() {
			return false
		} else if curA.kind() == NodeOpKind {
			aStack = append(aStack, curA.leftArg())
			aStack = append(aStack, curA.rightArg())
			bStack = append(bStack, curB.leftArg())
			bStack = append(bStack, curB.rightArg())
		}
	}
	return true
}

func TestSimpleInput(t *testing.T) {
	input := []lexer.Token{lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 1.}}
	want := &NodeOp{LeftArg: &LeafOp{Leaf: 1}, RightArg: &LeafOp{Leaf: 1}, Op: "+"}
	out, err := Parse(input)
	if !equalOpTree(want, out) || err != nil {
		t.Fatalf(`Parse(1 + 1) = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestTwoOp(t *testing.T) {
	input := []lexer.Token{lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 1.}}
	want := &NodeOp{
		LeftArg: &NodeOp{
			LeftArg:  &LeafOp{Leaf: 1},
			RightArg: &LeafOp{Leaf: 1},
			Op:       "+",
		},
		RightArg: &LeafOp{Leaf: 1},
		Op:       "+",
	}
	out, err := Parse(input)
	if !equalOpTree(want, out) || err != nil {
		t.Fatalf(`Parse(1 + 1 + 1) = %q, %v, wanted %#q, nil`, OpTreeToString(out), err, OpTreeToString(want))
	}
}

func TestResolveAddition(t *testing.T) {
	input := &NodeOp{LeftArg: &LeafOp{Leaf: 1}, RightArg: &LeafOp{Leaf: 1}, Op: "+"}
	want := 2.
	out, err := input.Resolve()
	if !almostEqual(out, want) || err != nil {
		t.Fatalf(`Parse(1 + 1) = %f, %v, wanted %f, <nil>`, out, err, want)
	}
}

func TestResolveSubtraction(t *testing.T) {
	input := &NodeOp{LeftArg: &LeafOp{Leaf: 1}, RightArg: &LeafOp{Leaf: 1}, Op: "-"}
	want := 0.
	out, err := input.Resolve()
	if !almostEqual(out, want) || err != nil {
		t.Fatalf(`Parse(1 - 1) = %f, %v, wanted %f, <nil>`, out, err, want)
	}
}

func TestResolveMultiplication(t *testing.T) {
	input := &NodeOp{LeftArg: &LeafOp{Leaf: 2}, RightArg: &LeafOp{Leaf: 1}, Op: "*"}
	want := 2.
	out, err := input.Resolve()
	if !almostEqual(out, want) || err != nil {
		t.Fatalf(`Parse(2 * 1) = %f, %v, wanted %f, <nil>`, out, err, want)
	}
}

func TestResolveDivision(t *testing.T) {
	input := &NodeOp{LeftArg: &LeafOp{Leaf: 2}, RightArg: &LeafOp{Leaf: 2}, Op: "/"}
	want := 1.
	out, err := input.Resolve()
	if !almostEqual(out, want) || err != nil {
		t.Fatalf(`Parse(2 / 2) = %f, %v, wanted %f, <nil>`, out, err, want)
	}
}

func TestDivisionByZeroError(t *testing.T) {
	input := &NodeOp{LeftArg: &LeafOp{Leaf: 2}, RightArg: &LeafOp{Leaf: 0}, Op: "/"}
	out, err := input.Resolve()
	if err == nil {
		t.Fatalf(`Parse(2 / 0) = %f, %v, wanted 0, division by zero error`, out, err)
	}
}

func TestOrderOfOps(t *testing.T) {
	input := []lexer.Token{lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 2.},
		lexer.OpToken{Val: "*"}, lexer.NumToken{Val: 2.}}
	want := &NodeOp{
		LeftArg: &LeafOp{Leaf: 1},
		RightArg: &NodeOp{
			LeftArg:  &LeafOp{Leaf: 2},
			RightArg: &LeafOp{Leaf: 2},
			Op:       "*",
		},
		Op: "+",
	}
	out, err := Parse(input)
	if err != nil || !equalOpTree(want, out) {
		t.Fatalf(
			`Parse(1 + 2 * 2) = %q, %v, wanted  %q, <nil>`,
			OpTreeToString(out),
			err,
			OpTreeToString(want),
		)
	}
}

func TestOOOBothSide(t *testing.T) {
	input := []lexer.Token{lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 2.},
		lexer.OpToken{Val: "*"}, lexer.NumToken{Val: 2.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 1.}}
	want := &NodeOp{
		LeftArg: &NodeOp{
			LeftArg: &LeafOp{Leaf: 1},
			RightArg: &NodeOp{
				LeftArg:  &LeafOp{Leaf: 2},
				RightArg: &LeafOp{Leaf: 2},
				Op:       "*",
			},
			Op: "+",
		},
		RightArg: &LeafOp{Leaf: 1},
		Op:       "+",
	}
	out, err := Parse(input)
	if err != nil || !equalOpTree(want, out) {
		t.Fatalf(
			`Parse(1 + 2 * 2 + 1) = %q, %v, wanted  %q, <nil>`,
			OpTreeToString(out),
			err,
			OpTreeToString(want),
		)
	}
}

func TestOOODoubleMult(t *testing.T) {
	input := []lexer.Token{lexer.NumToken{Val: 1.},
		lexer.OpToken{Val: "*"}, lexer.NumToken{Val: 2.},
		lexer.OpToken{Val: "+"}, lexer.NumToken{Val: 2.},
		lexer.OpToken{Val: "*"}, lexer.NumToken{Val: 3.}}
	want := &NodeOp{
		LeftArg: &NodeOp{
			LeftArg:  &LeafOp{Leaf: 1},
			RightArg: &LeafOp{Leaf: 2},
			Op:       "*",
		},
		RightArg: &NodeOp{
			LeftArg:  &LeafOp{Leaf: 2},
			RightArg: &LeafOp{Leaf: 3},
			Op:       "*",
		},
		Op: "+",
	}
	out, err := Parse(input)
	if err != nil || !equalOpTree(want, out) {
		t.Fatalf(
			`Parse(1 * 2 + 2 * 3) = %q, %v, wanted  %q, <nil>`,
			OpTreeToString(out),
			err,
			OpTreeToString(want),
		)
	}
}
