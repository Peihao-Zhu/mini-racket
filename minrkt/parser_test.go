package minrkt

import (
	"testing"
)

var tokenLP = Token{TOK_LPAREN, 0, "("}
var tokenAdd = Token{TOK_ADD, 0, "+"}
var token0 = Token{TOK_NUM, 0, "0"}
var token1 = Token{TOK_NUM, 1, "1"}
var token2 = Token{TOK_NUM, 2, "2"}
var token3 = Token{TOK_NUM, 3, "3"}
var token4 = Token{TOK_NUM, 4, "4"}
var tokenRP = Token{TOK_RPAREN, 0, ")"}
var tokenPlus2 = Token{TOK_NUM, 2, "+2"}
var tokenMinus3 = Token{TOK_NUM, -3, "-3.0"}
var tokenSub = Token{TOK_SUB, 0, "-"}
var tokenMUL = Token{TOK_MUL, 0, "*"}
var tokenDIV = Token{TOK_DIV, 0, "/"}
var tokenAND = Token{TOK_AND, 0, "and"}
var tokenOR = Token{TOK_OR, 0, "or"}
var tokenNOT = Token{TOK_NOT, 0, "not"}
var tokenLargeEqual = Token{TOK_LARGEEQUAL, 0, ">="}
var tokenLessEqual = Token{TOK_LESSEQUAL, 0, "<="}
var tokenEqual = Token{TOK_EQUAL, 0, "=="}
var tokenIf = Token{TOK_IF, 0, "if"}
var tokenTrue = Token{TOK_TRUE, 0, "true"}
var tokenFalse = Token{TOK_FALSE, 0, "false"}
var tokenDefine = Token{TOK_DEFINE, 0, "define"}
var tokenIdentifierX = Token{TOK_IDENTIFIER, 0, "x"}
var tokenIdentifierY = Token{TOK_IDENTIFIER, 0, "y"}
var tokenIdentifierFib = Token{TOK_IDENTIFIER, 0, "fib"}

func TestParse(t *testing.T) {

	// (+ 2 3)
	var tokens = []Token{tokenLP, tokenAdd, token2, token3, tokenRP}
	root, _ := Parse(tokens)
	want := "+ 2.00 3.00 "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}
	// (+ +2 (- -3))
	tokens = []Token{tokenLP, tokenAdd, tokenPlus2, tokenLP, tokenSub, tokenMinus3, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	want = "+ 2.00 - -3.00 "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

	// 2
	tokens = []Token{token2}
	root, _ = Parse(tokens)
	want = "2.00 "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

	// (+)
	tokens = []Token{tokenLP, tokenAdd, tokenRP}
	root, _ = Parse(tokens)
	want = "+ "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

	// (*)
	tokens = []Token{tokenLP, tokenMUL, tokenRP}
	root, _ = Parse(tokens)
	want = "* "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}
	// (and false (/ 4 0))
	tokens = []Token{tokenLP, tokenAND, tokenFalse, tokenLP, tokenDIV, token4, token0, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	want = "and false / 4.00 0.00 "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}
	// (>= 4 true)
	tokens = []Token{tokenLP, tokenLargeEqual, token4, tokenTrue, tokenRP}
	root, _ = Parse(tokens)
	want = ">= 4.00 true "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

	// detect error
	// (-)
	tokens = []Token{tokenLP, tokenSub, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (-)")
	}
	// (/)
	tokens = []Token{tokenLP, tokenDIV, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (/)")
	}
	// (2)
	tokens = []Token{tokenLP, token2, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (2)")
	}
	// (+ + 4)
	tokens = []Token{tokenLP, tokenAdd, tokenAdd, token2, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (+ + 4)")
	}
	// (- 2
	tokens = []Token{tokenLP, tokenSub, token2}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (- 2")
	}
	// + 2 )
	tokens = []Token{tokenAdd, token2, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression + 2)")
	}
	// (+ 2 - 2 3)
	tokens = []Token{tokenRP, tokenAdd, token2, tokenSub, token2, token3, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (+ 2 - 2 3)")
	}
	// ( )
	tokens = []Token{tokenLP, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression ( )")
	}

	// (if and 2 3)
	tokens = []Token{tokenLP, tokenIf, tokenAND, token2, token3, tokenRP}
	if _, err := Parse(tokens); err == nil {
		t.Error("expected parser error doesn't show up for expression (if and 2 3)")
	}

}

func TestVariableAndFunctionParser(t *testing.T) {
	// (define x 2)
	var tokens = []Token{tokenLP, tokenDefine, tokenIdentifierX, token2, tokenRP}
	root, _ := Parse(tokens)
	want := "define x 2.00 "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

	// (define y (* x x))
	tokens = []Token{tokenLP, tokenDefine, tokenIdentifierY, tokenLP, tokenMUL, tokenIdentifierX, tokenIdentifierX, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	want = "define y * x x "
	if result := root.Print(); result != want {
		t.Error("expected parsed tree is", want, " but got", result)
	}

}
