package calculator

import (
    "fmt"
    // "math"
    // "strconv"
    "strings"
    // "regexp"
)

func Calculate(expression string) (Value, error) {
    // Handle empty input
    expression = strings.TrimSpace(expression)
    if expression == "" {
        return Value{Number: 0, Type: NumberType}, nil
    }

    // Handle comments
    if strings.HasPrefix(expression, "#") || strings.HasPrefix(expression, "//") {
        return Value{Type: CommentType}, nil
    }

    // Create default context
    ctx := &Context{
        Variables: make(map[string]Value),
        TimeZone: "UTC", 
        PPI: 96,
        EmSize: 16,
    }

    // Handle variable assignment
    if strings.Contains(expression, "=") {
        parts := strings.SplitN(expression, "=", 2)
        varName := strings.TrimSpace(parts[0])
        varExpr := strings.TrimSpace(parts[1])
        
        result, err := Calculate(varExpr)
        if err != nil {
            return Value{IsError: true, ErrorMsg: err.Error()}, err
        }
        
        ctx.Variables[varName] = result
        return result, nil
    }

    // Tokenize
    tokens := tokenize(expression)

    // Parse
    expr, err := parse(tokens)
    if err != nil {
        return Value{IsError: true, ErrorMsg: err.Error()}, err
    }

    // Evaluate
    result := expr.Evaluate(ctx)
    if result.IsError {
        return result, fmt.Errorf(result.ErrorMsg)
    }

    return result, nil
}