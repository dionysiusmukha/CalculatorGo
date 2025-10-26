package logic

import (
	"Calculator/net"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func IsValidVarName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if !unicode.IsLetter(rune(name[0])) {
		return false
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func looksLikeMath(s string) bool {
	hasOp := false
	for _, r := range s {
		switch r {
		case '+', '-', '*', '/':
			hasOp = true
		case '.', '(', ')':
			
		default:
			if !(unicode.IsDigit(r) || unicode.IsLetter(r) || unicode.IsSpace(r)) {
				return false
			}
		}
	}
	return hasOp
}


func Interpret(cmd string, vars map[string]interface{}) (string, error) {
	// 1) Ссылка на уже сохранённую переменную по имени
	if val, ok := vars[cmd]; ok {
		switch v := val.(type) {
		case float64, int:
			return fmt.Sprintf("%v", v), nil
		case string:
			preview := v
			if len(v) > 300 {
				preview = v[:300] + "... [sliced]"
			}
			return preview, nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}

	// 2) Присваивание: x = ...
	if parts := strings.SplitN(cmd, "=", 2); len(parts) == 2 {
		name := strings.TrimSpace(parts[0])
		rhs := strings.TrimSpace(parts[1])

		if !IsValidVarName(name) {
			return "", fmt.Errorf("invalid variable name: %s", name)
		}

		// Спец-случай: x = curl <url>
		if strings.HasPrefix(rhs, "curl ") {
			url := strings.TrimSpace(strings.TrimPrefix(rhs, "curl "))
			content, err := net.Curl(url)
			if err != nil {
				return "", err
			}
			vars[name] = content
			return fmt.Sprintf("%s = <HTML got>", name), nil
		}

		// Обычная формула
		tokens := tokenize(rhs)
		rpn := toRPN(tokens)
		val, err := evalRPNWithInterface(rpn, vars)
		if err != nil {
			return "", err
		}
		vars[name] = val
		return fmt.Sprintf("%s = %v", name, val), nil
	}

	// 3) Явный curl
	if strings.HasPrefix(cmd, "curl ") {
		url := strings.TrimSpace(strings.TrimPrefix(cmd, "curl "))
		content, err := net.Curl(url)
		if err != nil {
			return "", err
		}
		vars["last_curl"] = content
		preview := strings.TrimSpace(content)
		if len(preview) > 300 {
			preview = preview[:300] + "... [sliced]"
		}

		return fmt.Sprintf("HTML got with %s:\n\n%s", url, preview), nil
	}

	// 4) Классификация/чат ИЛИ арифметика
	if !strings.Contains(cmd, "=") && !strings.HasPrefix(cmd, "curl ") {
		// если это НЕ «похоже на математику», пробуем экстрактор/пайплайн
		if !looksLikeMath(cmd) {
			ex, err := net.ExtractFreeForm(cmd)
			if err == nil && ex != nil && strings.ToLower(ex.Action) == "curl" && strings.TrimSpace(ex.URL) != "" {
				content, err := net.Curl(strings.TrimSpace(ex.URL))
				if err != nil {
					return "", err
				}
				vars["last_curl"] = content

				instr := strings.TrimSpace(ex.Instruction)
				if instr == "" {
					instr = "Сделай краткую сводку содержимого страницы."
				}
				answer, err := net.AskWithContext(instr, content)
				if err != nil {
					return "", fmt.Errorf("DeepSeek err (context): %v", err)
				}
				return answer, nil
			}
			// фолбэк
			resp, err := net.SendToDeepSeek(cmd)
			if err != nil {
				return "", fmt.Errorf("DeepSeek err: %v", err)
			}
			return resp, nil
		}

		// иначе — «похоже на математику»: парсим выражение
		tokens := tokenize(cmd)
		rpn := toRPN(tokens)
		val, err := evalRPNWithInterface(rpn, vars)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", val), nil
	}


	// 5) Остальное — арифметика/выражение
	tokens := tokenize(cmd)
	rpn := toRPN(tokens)
	val, err := evalRPNWithInterface(rpn, vars)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", val), nil
}


func evalRPNWithInterface(tokens []string, vars map[string]interface{}) (float64, error) {
	stack := []float64{}

	for _, tok := range tokens {
		if isOperator(tok) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("missing operand")
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
					return 0, fmt.Errorf("division by zero")
				}
				res = a / b
			default:
				return 0, fmt.Errorf("unknown operator %q", tok)
			}

			stack = append(stack, res)
		} else {
			if val, ok := vars[tok]; ok {
				switch v := val.(type) {
				case float64:
					stack = append(stack, v)
				case int:
					stack = append(stack, float64(v))
				case string:
					return 0, fmt.Errorf("variable %q contains string, not number", tok)
				default:
					return 0, fmt.Errorf("unknown type for variable %q", tok)
				}
			} else {
				f, err := strconv.ParseFloat(tok, 64)
				if err != nil {
					return 0, fmt.Errorf("unknown token %q", tok)
				}
				stack = append(stack, f)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}
	return stack[0], nil
}
