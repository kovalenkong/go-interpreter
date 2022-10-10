package functions

import (
	"fmt"
	"math"
)

func Sum(args ...any) (any, error) {
	var result float64
	for _, el := range args {
		switch val := el.(type) {
		case float64:
			result += val
		case []float64:
			for _, subarg := range val {
				result += subarg
			}
		default:
			return nil, fmt.Errorf("expected float64, got %T", el)
		}
	}
	return result, nil
}

func Len(args ...any) (any, error) {
	return float64(len(args)), nil
}

func And(args ...any) (any, error) {
	for _, el := range args {
		cond, ok := el.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", el)
		}
		if !cond {
			return false, nil
		}
	}
	return true, nil
}

func Or(args ...any) (any, error) {
	for _, el := range args {
		cond, ok := el.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", el)
		}
		if cond {
			return true, nil
		}
	}
	return false, nil
}

func If(args ...any) (any, error) {
	if length := len(args); length != 3 {
		return nil, fmt.Errorf("expected 3 args, got %d", length)
	}
	cond, ok := args[0].(bool)
	if !ok {
		return nil, fmt.Errorf("expected bool, got %T", args[0])
	}
	if cond {
		return args[1], nil
	}
	return args[2], nil
}

func Round(args ...any) (any, error) {
	if length := len(args); length == 0 || length > 2 {
		return nil, fmt.Errorf("expected 1 or 2 args, got %d", length)
	}
	val, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("expected float64, got %T", args[0])
	}
	var precision float64
	if len(args) == 2 {
		var ok bool
		precision, ok = args[1].(float64)
		if !ok {
			return nil, fmt.Errorf("expected float64, got %T", args[1])
		}
	}
	if precision < 0 {
		return nil, fmt.Errorf("round precision should be >= 0, got %f", precision)
	}
	ratio := math.Pow(10, precision)
	return math.Round(val*ratio) / ratio, nil
}

func Mean(args ...any) (any, error) {
	var total float64
	var count float64
	for _, arg := range args {
		switch val := arg.(type) {
		case float64:
			total += val
			count++
		case []float64:
			for _, sub := range val {
				total += sub
				count++
			}
		default:
			return nil, fmt.Errorf("expected float64, got %T", arg)
		}
	}
	return total / count, nil
}

func Min(args ...any) (any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected 1 or more args, got 0")
	}
	min := math.Inf(1)
	for _, arg := range args {
		val, ok := arg.(float64)
		if !ok {
			return nil, fmt.Errorf("expected float64, got %T", arg)
		}
		if val < min {
			min = val
		}
	}
	return min, nil
}

func Max(args ...any) (any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected 1 or more args, got 0")
	}
	max := math.Inf(-1)
	for _, arg := range args {
		val, ok := arg.(float64)
		if !ok {
			return nil, fmt.Errorf("expected float64, got %T", arg)
		}
		if val > max {
			max = val
		}
	}
	return max, nil
}

func Ifs(args ...any) (any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected >=2 arguments, got 0")
	}
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("expected even number of arguments, got %d", len(args))
	}
	for i := 0; i < len(args); i += 2 {
		cond, ok := args[i].(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool condition, got %T", args[i])
		}
		if cond {
			return args[i+1], nil
		}
	}

	return nil, fmt.Errorf("none of the conditions turned out to be true")
}
