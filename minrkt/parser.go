package minrkt

import (
	"fmt"
)

// will be the value of MapIdentifier if the key is a function name
type functionValue struct {
	args []string
	body Exp
}

type Params struct {
	MapIdentifier map[string]interface{}
	CallStack     []map[string]interface{}
}

type TypeEnum int

const (
	TYPE_FLOAT64 TypeEnum = iota
	TYPE_BOOLEAN
	TYPE_ERROR
	TYPE_DEFINE
	TYPE_NOTIFICATION
)

type Exp interface {
	Eval(p Params) (interface{}, TypeEnum, error)
	Print() string
}

type ExpOperator struct {
	opeType  string
	operands []Exp
}

type ExpNum struct {
	val float64
}

type ExpBool struct {
	val bool
}

// can be variable or function
type ExpIdentifier struct {
	val string
}

type ExpFunction struct {
}

func newExpOperator(ope string) *ExpOperator {
	return &ExpOperator{opeType: ope}
}

func newExpNum(num float64) *ExpNum {
	return &ExpNum{val: num}
}

func newExpBool(val bool) *ExpBool {
	return &ExpBool{val: val}
}

func newExpIdentifier(val string) *ExpIdentifier {
	return &ExpIdentifier{val: val}
}

var idx = 0

func Parse(tokens []Token) (Exp, error) {
	idx = 0
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression, you should input an expression")
	}
	// the initial token for expression should be a number, ( , true or false
	if len(tokens) == 1 {
		if tokens[0].tokenType == TOK_NUM {
			return newExpNum(tokens[idx].num), nil
		} else if tokens[0].tokenType == TOK_TRUE {
			return newExpBool(true), nil
		} else if tokens[0].tokenType == TOK_FALSE {
			return newExpBool(false), nil
		} else if tokens[0].tokenType == TOK_IDENTIFIER {
			return newExpIdentifier(tokens[0].val), nil
		} else {
			return nil, fmt.Errorf("for expression with single length, the token should be number ,true or false")
		}
	} else if tokens[0].tokenType == TOK_LPAREN {
		if root, err := buildPasedTree(tokens); err != nil {
			return nil, err
		} else if idx != len(tokens) {
			return nil, fmt.Errorf("there shouldn't have any expression outside paired parentheses")
		} else {
			return root, nil
		}
	} else {
		return nil, fmt.Errorf("for multiple length expression, first token should be (")
	}
}

func isOperator(tokenType TokenType) bool {
	if tokenType == TOK_ADD || tokenType == TOK_SUB || tokenType == TOK_MUL || tokenType == TOK_DIV ||
		tokenType == TOK_AND || tokenType == TOK_OR || tokenType == TOK_NOT || tokenType == TOK_LARGE ||
		tokenType == TOK_LARGEEQUAL || tokenType == TOK_EQUAL || tokenType == TOK_LESS || tokenType == TOK_LESSEQUAL ||
		tokenType == TOK_IF || tokenType == TOK_DEFINE {
		return true
	}
	return false
}

/*
Each time we encounter (, invoke the buildPasedTree() to construct the expression
wrapped by the parentheses
*/
func buildPasedTree(tokens []Token) (Exp, error) {
	idx++ // initially the token is (, we don't need to check anymore
	if idx >= len(tokens) {
		return nil, fmt.Errorf("expression end too early")
	}
	// note first identifier after ( can be function name, e.g (addx x)
	if !isOperator(tokens[idx].tokenType) && tokens[idx].tokenType != TOK_IDENTIFIER {
		return nil, fmt.Errorf("left parentheses should always followed by an operator")
	}
	root := newExpOperator(tokens[idx].val)
	idx++
	// quickly check the token after operator is not an operator
	if idx < len(tokens) && isOperator(tokens[idx].tokenType) {
		fmt.Println(tokens[idx-1], tokens[idx])
		return nil, fmt.Errorf("operator shouldn't followed by an operator")
	} else if idx < len(tokens) && (root.opeType == "-" || root.opeType == "/") && tokens[idx].tokenType == TOK_RPAREN {
		return nil, fmt.Errorf("- or / shouldn't followed by )")
	}
	// check all operands with the oparator
	for idx < len(tokens) {
		curToken := tokens[idx]
		idx++
		if curToken.tokenType == TOK_NUM {
			root.operands = append(root.operands, newExpNum(curToken.num))
		} else if curToken.tokenType == TOK_TRUE {
			root.operands = append(root.operands, newExpBool(true))
		} else if curToken.tokenType == TOK_FALSE {
			root.operands = append(root.operands, newExpBool(false))
		} else if curToken.tokenType == TOK_LPAREN {
			idx--
			if subOperand, err := buildPasedTree(tokens); err != nil {
				return nil, err
			} else {
				root.operands = append(root.operands, subOperand)
			}
		} else if curToken.tokenType == TOK_RPAREN {
			return root, nil
		} else if curToken.tokenType == TOK_IDENTIFIER {
			root.operands = append(root.operands, newExpIdentifier(curToken.val))
		}
	}
	return nil, fmt.Errorf("you miss the right parentheses")
}
