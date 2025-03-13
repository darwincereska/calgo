package calculator

import (
	"math"
	"strings"
	"fmt"
)

type UnitConversion struct {
	Value    Expression
	FromUnit string
	ToUnit   string
}

type FunctionCall struct {
	Name string
	Args []Expression
}



func (e *UnitConversion) Evaluate(ctx *Context) Value {
	val := e.Value.Evaluate(ctx)
	if val.IsError {
		return val
	}

	fromUnit := strings.ToLower(e.FromUnit)
	toUnit := strings.ToLower(e.ToUnit)

	// Handle currency conversions
	if fromRate, ok := ctx.Variables[fromUnit]; ok {
		if toRate, ok := ctx.Variables[toUnit]; ok {
			rate := toRate.Number / fromRate.Number
			return Value{
				Number: val.Number * rate,
				Unit:   toUnit,
				Type:   CurrencyType,
			}
		}
	}

	// Handle metric/imperial conversions
	for _, category := range unitConversions {
		if fromRatio, ok := category[fromUnit]; ok {
			if toRatio, ok := category[toUnit]; ok {
				result := val.Number * (fromRatio / toRatio)
				return Value{
					Number: result,
					Unit:   toUnit,
					Type:   UnitType,
				}
			}
		}
	}

	return Value{
		IsError:  true,
		ErrorMsg: "unsupported unit conversion: " + fromUnit + " to " + toUnit,
	}
}

func (e *FunctionCall) Evaluate(ctx *Context) Value {
	args := make([]Value, len(e.Args))
	for i, arg := range e.Args {
		args[i] = arg.Evaluate(ctx)
		if args[i].IsError {
			return args[i]
		}
	}

	switch e.Name {
	case "sqrt":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "sqrt requires 1 argument"}
		}
		if args[0].Number < 0 {
			return Value{IsError: true, ErrorMsg: "cannot take square root of negative number"}
		}
		return Value{Number: math.Sqrt(args[0].Number), Type: NumberType}

	case "abs":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "abs requires 1 argument"}
		}
		return Value{Number: math.Abs(args[0].Number), Type: NumberType}

	case "sin":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "sin requires 1 argument"}
		}
		return Value{Number: math.Sin(args[0].Number), Type: NumberType}

	case "cos":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "cos requires 1 argument"}
		}
		return Value{Number: math.Cos(args[0].Number), Type: NumberType}

	case "tan":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "tan requires 1 argument"}
		}
		return Value{Number: math.Tan(args[0].Number), Type: NumberType}

	case "log":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "log requires 1 argument"}
		}
		return Value{Number: math.Log10(args[0].Number), Type: NumberType}

	case "ln":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "ln requires 1 argument"}
		}
		return Value{Number: math.Log(args[0].Number), Type: NumberType}

	case "round":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "round requires 1 argument"}
		}
		return Value{Number: math.Round(args[0].Number), Type: NumberType}

	case "floor":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "floor requires 1 argument"}
		}
		return Value{Number: math.Floor(args[0].Number), Type: NumberType}

	case "ceil":
		if len(args) != 1 {
			return Value{IsError: true, ErrorMsg: "ceil requires 1 argument"}
		}
		return Value{Number: math.Ceil(args[0].Number), Type: NumberType}

	default:
		return Value{IsError: true, ErrorMsg: "unknown function: " + e.Name}
	}
}

func handleUnitConversion(value float64, fromUnit string, toUnit string) (Value, error) {
    // Check if it's a temperature conversion
    if fromUnit == "C" && toUnit == "F" {
        return Value{Number: value * 9/5 + 32, Unit: "F", Type: TemperatureType}, nil
    }
    if fromUnit == "F" && toUnit == "C" {
        return Value{Number: (value - 32) * 5/9, Unit: "C", Type: TemperatureType}, nil
    }

    // Try standard unit conversion
    result, err := convertUnit(value, fromUnit, toUnit)
    if err == nil {
        return Value{Number: result, Unit: toUnit, Type: UnitType}, nil
    }

    return Value{IsError: true, ErrorMsg: fmt.Sprintf("Cannot convert from %s to %s", fromUnit, toUnit)}, nil
}