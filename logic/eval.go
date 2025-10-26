package logic

import (
	"fmt"
	"strconv"
)

func evalRPN(tokens []string, vars map[string]float64) (float64, error) {
	var stack []float64

	for _, tok := range tokens {
		if isOperator(tok) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Missing operand")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch tok {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("Division by zero")
				}
				res = a / b
			}
			stack = append(stack, res)
			continue
		}
		num, err := strconv.ParseFloat(tok, 64)
		if err == nil {
			stack = append(stack, num)
		} else if v, ok := vars[tok]; ok {
			stack = append(stack, v)
		} else {
			return 0, fmt.Errorf("Unknown token %q", tok)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("Invalid expression: leftover stack")
	}
	return stack[0], nil
}
