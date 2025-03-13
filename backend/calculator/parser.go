package calculator

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

const (
	Number TokenType = iota
	Operator
	Unit
	Variable
	Function
	LeftParen
	RightParen
	Comment
	Label
)

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) Evaluate(ctx *Context) Value {
	return Value{Number: n.Value, Type: NumberType}
}

type VariableReference struct {
	Name string
}

func (v *VariableReference) Evaluate(ctx *Context) Value {
	if val, ok := ctx.Variables[v.Name]; ok {
		return val
	}
	return Value{IsError: true, ErrorMsg: "undefined variable: " + v.Name}
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator string
}

func (e *BinaryExpression) Evaluate(ctx *Context) Value {
	left := e.Left.Evaluate(ctx)
	right := e.Right.Evaluate(ctx)

	switch e.Operator {
	case "+":
		return Value{
			Number: left.Number + right.Number,
			Type:   NumberType,
		}
	case "-":
		return Value{
			Number: left.Number - right.Number,
			Type:   NumberType,
		}
	case "*":
		return Value{
			Number: left.Number * right.Number,
			Type:   NumberType,
		}
	case "/":
		if right.Number == 0 {
			return Value{
				IsError:  true,
				ErrorMsg: "division by zero",
			}
		}
		return Value{
			Number: left.Number / right.Number,
			Type:   NumberType,
		}
	case "^":
		return Value{
			Number: math.Pow(left.Number, right.Number),
			Type:   NumberType,
		}
	default:
		return Value{
			IsError:  true,
			ErrorMsg: "unknown operator: " + e.Operator,
		}
	}
}

func tokenize(input string) []Token {
	var tokens []Token
	input = strings.TrimSpace(input)

	for len(input) > 0 {
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			break
		}

		// Match numbers
		if match := regexp.MustCompile(`^-?\d*\.?\d+`).FindString(input); match != "" {
			tokens = append(tokens, Token{Type: Number, Value: match})
			input = input[len(match):]
			continue
		}

		// Match operators
		if match := regexp.MustCompile(`^[+\-*/^]`).FindString(input); match != "" {
			tokens = append(tokens, Token{Type: Operator, Value: match})
			input = input[len(match):]
			continue
		}

		// Match parentheses
		if input[0] == '(' {
			tokens = append(tokens, Token{Type: LeftParen, Value: "("})
			input = input[1:]
			continue
		}
		if input[0] == ')' {
			tokens = append(tokens, Token{Type: RightParen, Value: ")"})
			input = input[1:]
			continue
		}

		// Match variables, functions, units
		if match := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*`).FindString(input); match != "" {
			// Check if it's followed by parentheses (function)
			if len(input) > len(match) && input[len(match)] == '(' {
				tokens = append(tokens, Token{Type: Function, Value: match})
			} else if strings.HasSuffix(match, ":") {
				tokens = append(tokens, Token{Type: Label, Value: match})
			} else {
				// Check if it's a unit or variable
				isUnit := false
				// Add unit detection logic here
				if isUnit {
					tokens = append(tokens, Token{Type: Unit, Value: match})
				} else {
					tokens = append(tokens, Token{Type: Variable, Value: match})
				}
			}
			input = input[len(match):]
			continue
		}

		// Skip unknown characters
		input = input[1:]
	}

	return tokens
}

func parse(tokens []Token) (Expression, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	// Parse expressions with precedence
	return parseExpression(tokens, 0)
}

func parseExpression(tokens []Token, precedence int) (Expression, error) {
	var left Expression
	var err error

	left, tokens, err = parseTerm(tokens)
	if err != nil {
		return nil, err
	}

	for len(tokens) > 0 && isOperator(tokens[0]) && getOperatorPrecedence(tokens[0].Value) >= precedence {
		operator := tokens[0]
		tokens = tokens[1:]

		right, remainingTokens, err := parseTerm(tokens)
		if err != nil {
			return nil, err
		}

		left = &BinaryExpression{
			Left:     left,
			Right:    right,
			Operator: operator.Value,
		}
		tokens = remainingTokens
	}

	return left, nil
}

func parseTerm(tokens []Token) (Expression, []Token, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("unexpected end of expression")
	}

	token := tokens[0]
	tokens = tokens[1:]

	switch token.Type {
	case Number:
		val, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, tokens, fmt.Errorf("invalid number: %s", token.Value)
		}
		return &NumberLiteral{Value: val}, tokens, nil

	case Variable:
		return &VariableReference{Name: token.Value}, tokens, nil

	case Function:
		if len(tokens) == 0 || tokens[0].Type != LeftParen {
			return nil, tokens, fmt.Errorf("expected ( after function name")
		}
		tokens = tokens[1:] // Skip (

		var args []Expression
		for len(tokens) > 0 && tokens[0].Type != RightParen {
			arg, newTokens, err := parseTerm(tokens)
			if err != nil {
				return nil, tokens, err
			}
			args = append(args, arg)
			tokens = newTokens

			if len(tokens) > 0 && tokens[0].Type == RightParen {
				break
			}
		}

		if len(tokens) == 0 || tokens[0].Type != RightParen {
			return nil, tokens, fmt.Errorf("unclosed function call")
		}
		tokens = tokens[1:] // Skip )

		fc := &FunctionCall{Name: token.Value, Args: args}

		// Implement Expression interface
		var _ Expression = fc

		return fc, tokens, nil

	case LeftParen:
		expr, err := parseExpression(tokens, 0)
		if err != nil {
			return nil, tokens, err
		}
		if len(tokens) == 0 || tokens[0].Type != RightParen {
			return nil, tokens, fmt.Errorf("unclosed parenthesis")
		}
		tokens = tokens[1:] // Skip )
		return expr, tokens, nil

	default:
		return nil, tokens, fmt.Errorf("unexpected token: %v", token)
	}
}

func isOperator(token Token) bool {
	return token.Type == Operator
}

func getOperatorPrecedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	default:
		return 0
	}
}

func parseNumber(s string) (float64, error) {
	// Handle hex numbers
	if strings.HasPrefix(s, "0x") {
		val, err := strconv.ParseInt(s[2:], 16, 64)
		return float64(val), err
	}
	// Handle binary numbers
	if strings.HasPrefix(s, "0b") {
		val, err := strconv.ParseInt(s[2:], 2, 64)
		return float64(val), err
	}
	// Handle octal numbers
	if strings.HasPrefix(s, "0o") {
		val, err := strconv.ParseInt(s[2:], 8, 64)
		return float64(val), err
	}
	return strconv.ParseFloat(s, 64)
}

func isUnit(word string) bool {
	units := map[string]bool{
		// Length
		"m": true, "cm": true, "mm": true, "km": true, "in": true, "ft": true,
		"inch": true, "feet": true, "yard": true, "yd": true, "mile": true,
		// Weight
		"kg": true, "g": true, "mg": true, "lb": true, "oz": true,
		// Currency
		"USD": true, "EUR": true, "GBP": true, "CAD": true,
		// Time
		"second": true, "minute": true, "hour": true, "day": true, "week": true,
		// Temperature
		"C": true, "F": true, "K": true,
	}
	return units[strings.ToLower(word)]
}

func isCurrency(word string) bool {
 currencies := map[string]bool{
  "USD": true, "EUR": true, "GBP": true, "CAD": true,
  "AUD": true, "JPY": true, "CHF": true, "CNY": true,
  "HKD": true, "NZD": true, "SEK": true, "KRW": true,
  "SGD": true, "NOK": true, "MXN": true, "INR": true,
  "RUB": true, "ZAR": true, "TRY": true, "BRL": true,
  "TWD": true, "DKK": true, "PLN": true, "THB": true,
  "IDR": true, "HUF": true, "CZK": true, "ILS": true,
  "CLP": true, "PHP": true, "AED": true, "COP": true,
  "SAR": true, "MYR": true, "RON": true,
 }
 return currencies[strings.ToUpper(word)]
}