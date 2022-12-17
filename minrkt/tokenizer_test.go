package minrkt

import (
	"testing"
)

// go test . -test.v // run the test file
func TestTokenizer_NextToken(t *testing.T) {
	var preToken Token
	// test correct use case
	if token, newRemainder, _ := NextToken("(+ 2 3)", preToken); token.tokenType != TOK_LPAREN {
		t.Error("expected token type is 1 but got ", token.tokenType)
	} else if newRemainder != "+ 2 3)" {
		t.Error("expected remaining string is ", "+ 2 3)", " but got ", newRemainder)
	}

	if token, newRemainder, _ := NextToken("+ 2.0 3)", preToken); token.tokenType != TOK_ADD {
		t.Error("expected token type is 4 but got ", token.tokenType)
	} else if newRemainder != " 2.0 3)" {
		t.Error("expected remaining string is ", "2.0 3)", " but got ", newRemainder)
	}

	if token, newRemainder, _ := NextToken("2.0 3)", preToken); token.tokenType != TOK_NUM {
		t.Error("expected token type is 3 but got ", token.tokenType)
	} else if newRemainder != " 3)" {
		t.Error("expected remaining string is ", " 3)", " but got ", newRemainder)
	}

	if token, newRemainder, _ := NextToken(")", preToken); token.tokenType != TOK_RPAREN {
		t.Error("expected token type is 2 but got ", token.tokenType)
	} else if newRemainder != "" {
		t.Error("expected remaining string is ", "", " but got ", newRemainder)
	}

	if token, newRemainder, _ := NextToken("-2.345", preToken); token.tokenType != TOK_NUM {
		t.Error("expected token type is 3 but got ", token.tokenType)
	} else if newRemainder != "" {
		t.Error("expected remaining string is ", "", " but got ", newRemainder)
	}

	if token, newRemainder, _ := NextToken("+4124.1 )", preToken); token.tokenType != TOK_NUM {
		t.Error("expected token type is 3 but got ", token.tokenType)
	} else if newRemainder != " )" {
		t.Error("expected remaining string is ", " )", " but got ", newRemainder)
	}

	if token, _, _ := NextToken("", preToken); token.tokenType != TOK_EOF {
		t.Error("expected token type is 8 but got ", token.tokenType)
	}

	if token, _, _ := NextToken("if 2 3 4", preToken); token.tokenType != TOK_IF {
		t.Error("expected token type is 19 but got ", token.tokenType)
	}

	if token, _, _ := NextToken("and true", preToken); token.tokenType != TOK_AND {
		t.Error("expected token type is 9 but got ", token.tokenType)
	}

	if token, _, _ := NextToken("or true", preToken); token.tokenType != TOK_OR {
		t.Error("expected token type is 10 but got ", token.tokenType)
	}

	if token, _, _ := NextToken("not", preToken); token.tokenType != TOK_NOT {
		t.Error("expected token type is 11 but got ", token.tokenType)
	}

	// test error use case
	if _, _, err := NextToken("ab", preToken); err == nil {
		t.Error("expected error doesn't show up: ", err)
	}

	// (++2 4)
	preToken = Token{tokenType: TOK_ADD, val: "+"}
	if _, _, err := NextToken("+2)", preToken); err == nil {
		t.Error("expected error doesn't show up: ", err)
	}
}

func TestTokenizer_Tokenize(t *testing.T) {
	tokenLP := Token{TOK_LPAREN, 0, "("}
	tokenAdd := Token{TOK_ADD, 0, "+"}
	token2 := Token{TOK_NUM, 2, "2"}
	token3 := Token{TOK_NUM, 3, "3"}
	tokenRP := Token{TOK_RPAREN, 0, ")"}

	var want = []Token{tokenLP, tokenAdd, token2, token3, tokenRP}
	if tokens, _ := Tokenize("(+ 2 3)"); !compareTokens(tokens, want) {
		t.Error("expected token type is", want, " but got", tokens)
	}

	tokenPlus2 := Token{TOK_NUM, 2, "+2"}
	tokenMinus3 := Token{TOK_NUM, -3, "-3.0"}
	tokenSub := Token{TOK_SUB, 0, "-"}
	want = []Token{tokenLP, tokenAdd, tokenPlus2, tokenLP, tokenSub, tokenMinus3, tokenRP, tokenRP}
	if tokens, _ := Tokenize("(+ +2 (- -3.0))"); !compareTokens(tokens, want) {
		t.Error("expected token type is", want, " but got", tokens)
	}

}

func compareTokens(token1, token2 []Token) bool {
	for i, t1 := range token1 {
		if t1 != token2[i] {
			return false
		}
	}
	return true
}
