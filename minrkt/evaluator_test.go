package minrkt

import (
	"testing"
)

func TestEvaluator(t *testing.T) {
	var want float64
	// (+ 2 3)
	var tokens = []Token{tokenLP, tokenAdd, token2, token3, tokenRP}
	root, _ := Parse(tokens)
	want = 5
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected  evaluated result is", want, " but got", result)
	}
	// (* 2 3 (+ 2))
	tokens = []Token{tokenLP, tokenMUL, token2, token3, tokenLP, tokenAdd, token2, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	want = 12
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected  evaluated result is", want, " but got", result)
	}
	// (/ 2 (* 2 3) (+ 2))
	tokens = []Token{tokenLP, tokenDIV, token2, tokenLP, tokenMUL, token2, token3, tokenRP, tokenLP, tokenAdd, token2, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	want = 0.16666666666666666
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected evaluated result is", want, " but got", result)
	}
	// (+)
	tokens = []Token{tokenLP, tokenAdd, tokenRP}
	root, _ = Parse(tokens)
	want = 0
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected evaluated result is", want, " but got", result)
	}
	// (*)
	tokens = []Token{tokenLP, tokenMUL, tokenRP}
	root, _ = Parse(tokens)
	want = 1
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected evaluated result is", want, " but got", result)
	}
	// (/ 2 (- 3 3))  -> divide by 0
	tokens = []Token{tokenLP, tokenDIV, token2, tokenLP, tokenSub, token3, token3, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(division by 0) doesn't show up for expression (/ 2 (- 3 3))")
	}
	// (not 4)
	tokens = []Token{tokenLP, tokenNOT, token4, tokenRP}
	root, _ = Parse(tokens)
	if result, _, _ := root.Eval(); result.(bool) != false {
		t.Error("expected evaluated result is false but got", result)
	}
	// (and false (/ 4 0))
	tokens = []Token{tokenLP, tokenAND, tokenFalse, tokenLP, tokenDIV, token4, token0, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	if result, _, _ := root.Eval(); result.(bool) != false {
		t.Error("expected evaluated result is false but got", result)
	}
	// (or (not true) (<= 2 3))
	tokens = []Token{tokenLP, tokenOR, tokenLP, tokenNOT, tokenTrue, tokenRP, tokenLP, tokenLessEqual, token2, token3, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	if result, _, _ := root.Eval(); result.(bool) != true {
		t.Error("expected evaluated result is true but got", result)
	}
	// (if (and (>= 1 2) (= 3 4)) (/ 1 0) (or true false))
	tokens = []Token{tokenLP, tokenIf, tokenLP, tokenAND, tokenLP, tokenLargeEqual, token1, token2, tokenRP, tokenLP,
		tokenEqual, token3, token4, tokenRP, tokenRP, tokenLP, tokenDIV, token1, token0, tokenRP, tokenLP, tokenOR, tokenTrue,
		tokenFalse, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	if result, _, _ := root.Eval(); result.(bool) != true {
		t.Error("expected evaluated result is true but got", result)
	}
	// (if 4 (/ 1 0) (or true false))
	tokens = []Token{tokenLP, tokenIf, token4, tokenLP, tokenDIV, token1, token0, tokenRP, tokenLP, tokenOR, tokenTrue,
		tokenFalse, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(division by 0) doesn't show up for expression (if 4 (/ 1 0) (or true false))")
	}
	// (>= 4 true)
	tokens = []Token{tokenLP, tokenLargeEqual, token4, tokenTrue, tokenRP}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(division by 0) doesn't show up for expression (>= 4 true)")
	}
}

func TestVariableAndFunctionEvaluator(t *testing.T) {
	var want float64
	// (define x (+ 1 2))
	var tokens = []Token{tokenLP, tokenDefine, tokenIdentifierX, tokenLP, tokenAdd, token1, token2, tokenRP, tokenRP}
	root, _ := Parse(tokens)
	root.Eval()
	root, _ = Parse([]Token{tokenIdentifierX})
	want = 3
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected  evaluated result is", want, " but got", result)
	}

	// (define (fib x) (if (<= x 1) x (+ (fib (- x 1)) (fib (- x 2)))))
	tokens = []Token{tokenLP, tokenDefine, tokenLP, tokenIdentifierFib, tokenIdentifierX, tokenRP, tokenLP, tokenIf, tokenLP,
		tokenLessEqual, tokenIdentifierX, token1, tokenRP, tokenIdentifierX, tokenLP, tokenAdd, tokenLP, tokenIdentifierFib, tokenLP, tokenSub, tokenIdentifierX,
		token1, tokenRP, tokenRP, tokenLP, tokenIdentifierFib, tokenLP, tokenSub, tokenIdentifierX, token2, tokenRP, tokenRP,
		tokenRP, tokenRP, tokenRP}
	root, _ = Parse(tokens)
	root.Eval()
	root, _ = Parse([]Token{tokenLP, tokenIdentifierFib, token4, tokenRP})
	want = 3
	if result, _, _ := root.Eval(); result.(float64) != want {
		t.Error("expected  evaluated result is", want, " but got", result)
	}

	// (x)
	tokens = []Token{tokenLP, tokenIdentifierX, tokenRP}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(application: not a procedure) doesn't show up for expression (x)")
	}

	// z
	tokens = []Token{tokenIdentifierY}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(yL undefined) doesn't show up for expression y")
	}

	// (fib 2 3)
	tokens = []Token{tokenLP, tokenIdentifierFib, token2, token3, tokenRP}
	root, _ = Parse(tokens)
	if _, _, err := root.Eval(); err == nil {
		t.Error("expected evaluation error(fib: arity mismatch) doesn't show up for expression (fib 2 3)")
	}

}
