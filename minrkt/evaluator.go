package minrkt

import "fmt"

func (e *ExpNum) Print() string {
	// fmt.Println(e.val)
	return fmt.Sprintf("%.2f ", e.val)
}

func (e *ExpBool) Print() string {
	// fmt.Println(e.val)
	return fmt.Sprintf("%t ", e.val)
}

func (e *ExpOperator) Print() string {
	// fmt.Println(e.opeType)
	var result = e.opeType + " "
	for _, child := range e.operands {
		result += child.Print()
	}
	return result
}

func (e *ExpIdentifier) Print() string {
	// fmt.Println(e.val)
	return fmt.Sprintf("%s ", e.val)
}

func (e *ExpBool) Eval(p Params) (interface{}, TypeEnum, error) {
	return e.val, TYPE_BOOLEAN, nil
}

func (e *ExpNum) Eval(p Params) (interface{}, TypeEnum, error) {
	return e.val, TYPE_FLOAT64, nil
}

func (e *ExpIdentifier) Eval(p Params) (interface{}, TypeEnum, error) {
	var val interface{}
	var ok = false
	argsMap := p.CallStack[len(p.CallStack)-1]
	if val, ok = argsMap[e.val]; !ok {
		if val, ok = p.MapIdentifier[e.val]; !ok {
			return nil, TYPE_ERROR, fmt.Errorf("%s: undifined", e.val)
		}
	}
	switch got := val.(type) {
	case float64:
		return got, TYPE_FLOAT64, nil
	case bool:
		return got, TYPE_BOOLEAN, nil
	default:
		return "#<procedure:" + e.val + ">", TYPE_NOTIFICATION, nil
	}

}

func (e *ExpOperator) Eval(p Params) (interface{}, TypeEnum, error) {
	var sum float64 = 0
	switch e.opeType {
	case "+":
		for _, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return got, t, err
			} else if t == TYPE_FLOAT64 {
				sum += got.(float64)
			} else {
				return got, t, fmt.Errorf("operand for + should be number")
			}
		}
		return sum, TYPE_FLOAT64, nil
	case "-":
		for i, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return got, t, err
			} else if t == TYPE_FLOAT64 {
				if i == 0 {
					sum += got.(float64)
				} else {
					sum -= got.(float64)
				}
			} else {
				return got, t, fmt.Errorf("operand for + should be number")
			}
		}
		// (- 9) -> -9
		if len(e.operands) == 1 {
			sum = -sum
		}
		return sum, TYPE_FLOAT64, nil
	case "*":
		sum = 1
		for _, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return got, t, err
			} else if t == TYPE_FLOAT64 {
				sum *= got.(float64)
			} else {
				return got, t, fmt.Errorf("operand for + should be number")
			}
		}
		return sum, TYPE_FLOAT64, nil
	case "/":
		for i, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return got, t, err
			} else if t == TYPE_FLOAT64 {
				if i == 0 {
					sum = got.(float64)
				} else {
					divisor := got.(float64)
					if divisor == 0 {
						return nil, TYPE_FLOAT64, fmt.Errorf("divide by zero error")
					}
					sum /= divisor
				}
			} else {
				return got, t, fmt.Errorf("operand for + should be number")
			}
		}
		// ( / 9) -> 1/9
		if len(e.operands) == 1 {
			sum = 1 / sum
		}
		return sum, TYPE_FLOAT64, nil
	case "and":
		res := true
		for _, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return nil, t, err
			} else if t == TYPE_BOOLEAN {
				res = res && got.(bool)
			} else if t == TYPE_FLOAT64 {
				continue
			} else {
				return nil, t, fmt.Errorf("operand for and should be boolean")
			}
			// short circuit
			if !res {
				break
			}
		}
		return res, TYPE_BOOLEAN, nil
	case "or":
		res := false
		for _, c := range e.operands {
			if got, t, err := c.Eval(p); err != nil {
				return nil, t, err
			} else if t == TYPE_BOOLEAN {
				value := got.(bool)
				res = res || value
			} else if t == TYPE_FLOAT64 {
				res = true
			} else {
				return nil, t, fmt.Errorf("operand for and should be boolean")
			}
			// short circuit
			if res {
				break
			}
		}
		return res, TYPE_BOOLEAN, nil
	case "not":
		// check there is only 1 operand for not
		if got, t, err := e.operands[0].Eval(p); err != nil {
			return got, t, err
		} else if t == TYPE_BOOLEAN {
			value := got.(bool)
			return !value, TYPE_BOOLEAN, nil
		} else if t == TYPE_FLOAT64 {
			// number stands for true
			// not true is false
			return false, TYPE_BOOLEAN, nil
		} else {
			return nil, t, fmt.Errorf("operand for and should be boolean")
		}
	case ">":
		firstNum, secondNum, t, err := getTwoNum(e, p)
		if err != nil {
			return nil, t, err
		}
		if firstNum > secondNum {
			return true, TYPE_BOOLEAN, nil
		} else {
			return false, TYPE_BOOLEAN, nil
		}
	case ">=":
		firstNum, secondNum, t, err := getTwoNum(e, p)
		if err != nil {
			return nil, t, err
		}
		if firstNum >= secondNum {
			return true, TYPE_BOOLEAN, nil
		} else {
			return false, TYPE_BOOLEAN, nil
		}
	case "=":
		firstNum, secondNum, t, err := getTwoNum(e, p)
		if err != nil {
			return nil, t, err
		}
		if firstNum == secondNum {
			return true, TYPE_BOOLEAN, nil
		} else {
			return false, TYPE_BOOLEAN, nil
		}
	case "<":
		firstNum, secondNum, t, err := getTwoNum(e, p)
		if err != nil {
			return nil, t, err
		}
		if firstNum < secondNum {
			return true, TYPE_BOOLEAN, nil
		} else {
			return false, TYPE_BOOLEAN, nil
		}
	case "<=":
		firstNum, secondNum, t, err := getTwoNum(e, p)
		if err != nil {
			return nil, t, err
		}
		if firstNum <= secondNum {
			return true, TYPE_BOOLEAN, nil
		} else {
			return false, TYPE_BOOLEAN, nil
		}
	case "if":
		if len(e.operands) != 3 {
			return nil, TYPE_ERROR, fmt.Errorf("if statement should have three expressions")
		}
		// check the first expression
		if got, t, err := e.operands[0].Eval(p); err != nil {
			return got, t, err
		} else {
			expression := false
			switch v := got.(type) {
			case bool:
				if v {
					expression = true
				}
			case float64:
				expression = true
			}

			if expression {
				// the first expression is true, get the second result
				return e.operands[1].Eval(p)
			} else {
				// get the third result
				return e.operands[2].Eval(p)
			}
		}
	case "define":
		if len(e.operands) != 2 {
			return nil, TYPE_ERROR, fmt.Errorf("define statement should have an identifier and an expression")
		}
		// get identifier
		switch v := e.operands[0].(type) {
		case *ExpIdentifier:
			if expression, t, err := e.operands[1].Eval(p); err != nil {
				return expression, t, err
			} else {
				p.MapIdentifier[v.val] = expression
			}
		case *ExpOperator:
			// the first operator token after define is function
			// value of map for function should be (args, function body)
			argsNum := len(v.operands)
			args := make([]string, argsNum)
			for i := 0; i < argsNum; i++ {
				args[i] = v.operands[i].(*ExpIdentifier).val
			}
			p.MapIdentifier[v.opeType] = functionValue{args: args, body: e.operands[1]}
		default:
			// type is not identifier
			return nil, TYPE_ERROR, fmt.Errorf("define statement should followed by an identifier")
		}
		return nil, TYPE_DEFINE, nil
	default:
		// function invocation will fall into here
		// (fib 2)
		argsNum := len(e.operands)
		args := make([]interface{}, argsNum)
		for i := 0; i < argsNum; i++ {
			// expressions are bound to function arguments
			args[i], _, _ = e.operands[i].Eval(p)
		}
		// check function name exist in environmnet
		if value, ok := p.MapIdentifier[e.opeType]; ok {
			// the identifier is function not variable
			if fv, ok := value.(functionValue); ok {
				// match the input parameters
				if len(fv.args) != argsNum {
					return nil, TYPE_ERROR, fmt.Errorf("arity mismatch")
				}
				// before we modify the argsMap, put the existing argsMap into stack

				argsMap := make(map[string]interface{})
				for i := 0; i < argsNum; i++ {
					argsMap[fv.args[i]] = args[i]
				}
				p.CallStack = append(p.CallStack, argsMap)
				// execute the function body
				res, t, err := fv.body.Eval(p)
				// pop up the upmost local variable
				p.CallStack = p.CallStack[:len(p.CallStack)-1]
				return res, t, err
			} else {
				return nil, TYPE_ERROR, fmt.Errorf("application: not a procedure")
			}
		} else {
			return nil, TYPE_ERROR, fmt.Errorf("%s undifined", e.opeType)
		}
	}
}

/*
for comparison operators like =, >=, >, <=, <, there operands are supposed to be numbers.
so get these two numbers otherwie return error.
*/
func getTwoNum(e *ExpOperator, p Params) (float64, float64, TypeEnum, error) {
	var firstNum, secondNum float64
	if len(e.operands) != 2 {
		return 0, 0, TYPE_ERROR, fmt.Errorf("arithmetic comparison should have two operands")
	}
	if got, t, err := e.operands[0].Eval(p); err != nil {
		return 0, 0, t, err
	} else if t == TYPE_FLOAT64 {
		firstNum = got.(float64)
	} else {
		return 0, 0, t, fmt.Errorf("operand for arithmetic comparison should be number")
	}
	if got, t, err := e.operands[1].Eval(p); err != nil {
		return 0, 0, t, err
	} else if t == TYPE_FLOAT64 {
		secondNum = got.(float64)
	} else {
		return 0, 0, t, fmt.Errorf("operand for arithmetic comparison should be number")
	}
	return firstNum, secondNum, TYPE_FLOAT64, nil
}
