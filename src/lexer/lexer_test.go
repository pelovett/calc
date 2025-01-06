package lexer

import "testing"

func equalTokenSlices(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ta := range a {
		tb := b[i]
		if ta.Kind() != tb.Kind() || ta.Val() != tb.Val() {
			return false
		}
	}
	return true
}

func TestSingleChar(t *testing.T) {
	input := "1"
	want := []Token{NumToken{val: 1.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("1") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestTwoChar(t *testing.T) {
	input := "11"
	want := []Token{NumToken{val: 11.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("11") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestSpaceInNumber(t *testing.T) {
	input := "1 1"
	want := []Token{NumToken{val: 1.}, NumToken{val: 1.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("11") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestSimpleAddition(t *testing.T) {
	input := "1 + 1"
	want := []Token{NumToken{val: 1.},
		OpToken{val: "+"}, NumToken{val: 1.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("1 + 1") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestSimpleSubtraction(t *testing.T) {
	input := "1 - 1"
	want := []Token{NumToken{val: 1.},
		OpToken{val: "-"}, NumToken{val: 1.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("1 - 1") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestNoSpace(t *testing.T) {
	input := "1+1"
	want := []Token{NumToken{val: 1.},
		OpToken{val: "+"}, NumToken{val: 1.}}
	out, err := Tokenize(input)
	if !equalTokenSlices(want, out) || err != nil {
		t.Fatalf(`Tokenize("1+1") = %q, %v, wanted %#q, nil`, out, err, want)
	}
}

func TestIllegalChar(t *testing.T) {
	input := "d"
	out, err := Tokenize(input)
	if out != nil || err == nil {
		t.Fatalf(`Tokenize("d") = %q, %v, wanted nil, illegal starting char at position 1: d`, out, err)
	}
}

func TestIllegalCharAfterNum(t *testing.T) {
	input := "1d"
	out, err := Tokenize(input)
	if out != nil || err == nil {
		t.Fatalf(`Tokenize("1d") = %q, %v, wanted nil, illegal starting char at position 2: d`, out, err)
	}
}
