package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

func main() {
	color.Cyan("欢迎使用 g-calc 命令行计算器！")
	color.Cyan("请输入表达式（支持加、减、乘、除和括号），输入'q'退出。")
	color.Cyan("使用左右箭头键移动光标。")

	rl, err := readline.New("> ")
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if input == "q" {
			break
		}
		result, err := evaluate(input)
		if err != nil {
			color.Red("错误: %v", err)
		} else {
			color.Yellow("结果: %v", result.Text('f', 10))
		}
	}
	color.Cyan("谢谢使用，再见！")
}

func evaluate(expr string) (*big.Float, error) {
	// 替换中文括号为英文括号
	expr = strings.ReplaceAll(expr, "（", "(")
	expr = strings.ReplaceAll(expr, "）", ")")

	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}
	return evalRPN(toRPN(tokens))
}

func tokenize(expr string) ([]string, error) {
	var tokens []string
	var num strings.Builder
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch {
		case c == ' ':
			continue
		case c >= '0' && c <= '9' || c == '.':
			num.WriteByte(c)
		case c == '+' || c == '-' || c == '*' || c == '/' || c == '(' || c == ')':
			if num.Len() > 0 {
				tokens = append(tokens, num.String())
				num.Reset()
			}
			tokens = append(tokens, string(c))
		default:
			return nil, fmt.Errorf("无效字符: %c", c)
		}
	}
	if num.Len() > 0 {
		tokens = append(tokens, num.String())
	}
	return tokens, nil
}

func toRPN(tokens []string) []string {
	var output []string
	var stack []string
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
			}
		default:
			output = append(output, token)
		}
	}
	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func evalRPN(tokens []string) (*big.Float, error) {
	var stack []*big.Float
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return nil, fmt.Errorf("无效的表达式")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result big.Float
			switch token {
			case "+":
				result.Add(a, b)
			case "-":
				result.Sub(a, b)
			case "*":
				result.Mul(a, b)
			case "/":
				if b.Sign() == 0 {
					return nil, fmt.Errorf("除数不能为零")
				}
				result.Quo(a, b)
			}
			stack = append(stack, &result)
		default:
			f, _, err := new(big.Float).Parse(token, 10)
			if err != nil {
				return nil, fmt.Errorf("无效的数字: %s", token)
			}
			stack = append(stack, f)
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("无效的表达式")
	}
	return stack[0], nil
}
