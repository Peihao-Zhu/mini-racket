package minrkt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// “ raw string literal, no special character here
var tokenRegexList = []string{
	`^(\()`,
	`^(\))`,
	`^([\-\+]?0\.[0-9]+)`,
	`^(0)`,
	`^([\-\+]?[1-9][0-9]*(?:\.[0-9]*)?)`,
	`^(\+)`,
	`^(\-)`,
	`^(\*)`,
	`^(\/)`,
	`^(and)`,
	`^(or)`,
	`^(not)`,
	`^(=)`,
	`^(>=)`,
	`^(<=)`,
	`^(>)`,
	`^(<)`,
	`^(true)`,
	`^(false)`,
	`^(if)`,
	`^(define)`,
	`^([a-zA-Z]([a-zA-Z]|[0-9]|_)*)`,
}

type Token struct {
	tokenType TokenType
	num       float64
	val       string
}

type TokenType int

const (
	TOK_INVALID TokenType = iota // increament from 0 to 18
	TOK_LPAREN
	TOK_RPAREN
	TOK_NUM
	TOK_ADD
	TOK_SUB
	TOK_MUL
	TOK_DIV
	TOK_EOF
	TOK_AND
	TOK_OR
	TOK_NOT
	TOK_EQUAL
	TOK_LARGE
	TOK_LARGEEQUAL // >=
	TOK_LESS
	TOK_LESSEQUAL // <=
	TOK_TRUE
	TOK_FALSE
	TOK_IF
	TOK_DEFINE
	TOK_IDENTIFIER
)

var re = regexp.MustCompile(strings.Join(tokenRegexList, "|"))
var wsRe = regexp.MustCompile(`^\s+`)

func NextToken(remainder string, preToken Token) (Token, string, error) {
	var hasWhitespaces bool = false
	// strip off whitespaces
	if ws := wsRe.FindStringSubmatch(remainder); ws != nil {
		remainder = remainder[len(ws[0]):]
		hasWhitespaces = true
	}
	if len(remainder) == 0 {
		return Token{tokenType: TOK_EOF}, remainder, nil
	}
	matched_arr := re.FindStringSubmatch(remainder)
	var tokenType TokenType
	if matched_arr == nil {
		return Token{tokenType: TOK_INVALID}, remainder, fmt.Errorf("invalid token: %s", remainder)
	}
	// to match which token type it corresponds to
	matched_token := matched_arr[0]
	for i := 1; i < 23; i++ {
		if matched_token == matched_arr[i] {
			tokenType = getTokenType(i)
			break
		}
	}

	remainder = remainder[len(matched_token):]
	value, _ := strconv.ParseFloat(matched_token, 64)
	curToken := Token{tokenType: tokenType, num: value, val: matched_token}
	if err := checkInvalidToken(preToken, curToken, hasWhitespaces); err != nil {
		return Token{tokenType: TOK_INVALID}, remainder, err
	}
	return curToken, remainder, nil
}

func checkInvalidToken(preToken, curToken Token, hasWhitespaces bool) error {
	preTokenType := preToken.tokenType
	if curToken.tokenType == TOK_NUM && !hasWhitespaces && (preTokenType == TOK_ADD || preTokenType == TOK_SUB ||
		preTokenType == TOK_MUL || preTokenType == TOK_DIV) {
		return fmt.Errorf("invalid token: %s", preToken.val+curToken.val)
	}
	return nil
}

func getTokenType(value int) TokenType {
	switch value {
	case 1:
		return TOK_LPAREN
	case 2:
		return TOK_RPAREN
	case 3, 4, 5:
		return TOK_NUM
	case 6:
		return TOK_ADD
	case 7:
		return TOK_SUB
	case 8:
		return TOK_MUL
	case 9:
		return TOK_DIV
	case 10:
		return TOK_AND
	case 11:
		return TOK_OR
	case 12:
		return TOK_NOT
	case 13:
		return TOK_EQUAL
	case 14:
		return TOK_LARGEEQUAL
	case 15:
		return TOK_LESSEQUAL
	case 16:
		return TOK_LARGE
	case 17:
		return TOK_LESS
	case 18:
		return TOK_TRUE
	case 19:
		return TOK_FALSE
	case 20:
		return TOK_IF
	case 21:
		return TOK_DEFINE
	case 22:
		return TOK_IDENTIFIER
	}
	return TOK_INVALID
}

func Tokenize(line string) ([]Token, error) {
	remainder := line
	var tokens []Token
	var preToken Token
	for {
		token, newRemainder, err := NextToken(remainder, preToken)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
		if len(newRemainder) == 0 {
			return tokens, nil
		}
		remainder = newRemainder
		preToken = token
	}
}
