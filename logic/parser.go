package logic

import (
	"fmt"
	"unicode"
)

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func tokenize(input string) []string {
	tokens := []string{}
	current := ""

	addToken := func() {
		if current != "" {
			tokens = append(tokens, current)
			current = ""
		}
	}

	for _, ch := range input {
		if unicode.IsSpace(ch) {
			addToken()
			continue
		}
		if unicode.IsDigit(ch) || string(ch) == "." {
			current += string(ch)
			continue
		}
		if unicode.IsLetter(ch) || string(ch) == "_" {
			current += string(ch)
			continue
		}
		if string(ch) == "-" {
			if current == "" && (len(tokens) == 0 || tokens[len(tokens)-1] == "(" || isOperator(tokens[len(tokens)-1])) {
				current += "-"
				continue
			}
			addToken()
			tokens = append(tokens, "-")
			continue
		}
		addToken()
		tokens = append(tokens, string(ch))
	}
	addToken()
	return tokens
}

func toRPN(tokens []string) []string {
	output := []string{}
	stack := []string{}

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	for _, tok := range tokens {
		_, ok := precedence[tok]
		if !ok && tok != "(" && tok != ")" {
			output = append(output, tok)
		} else if tok == "(" {
			stack = append(stack, tok)
		} else if tok == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				fmt.Println("Mismatched parentheses")
			} else {
				stack = stack[:len(stack)-1]
			}
		} else {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" || precedence[top] < precedence[tok] {
					break
				}
				output = append(output, top)
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, tok)
		}
	}
	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}
