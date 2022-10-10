package go_interpreter

import (
	"fmt"
	"github.com/kovalenkong/go-interpreter/functions"
	"testing"
)

func BenchmarkInterpreter_ExecuteFormula(b *testing.B) {
	interpreter := NewInterpreter(
		map[string]any{
			"X": 10.0,
			"Y": 20.0,
		},
		map[string]Func{
			"Sum": functions.Sum,
			"Len": functions.Len,
		},
	)
	for i := 0; i < b.N; i++ {
		res, err := interpreter.Execute(`X + Y * 72 / Sum(1;2;3)^Len(1;10)`)
		if err != nil {
			b.Fatalf("expected nil error, got %v", err)
		}
		val, ok := res.(float64)
		if !ok {
			b.Fatalf("expected float64 result, got %T", res)
		}
		if val != 50 {
			b.Fatalf("expected 50, got %f", val)
		}
	}
}

func BenchmarkInterpreter_ExecuteNumber(b *testing.B) {
	interpreter := NewInterpreter(nil, nil)
	var expected float64 = 1
	for i := 0; i < b.N; i++ {
		res, err := interpreter.Execute(`1`)
		if err != nil {
			b.Fatalf("expected nil error, got %v", err)
		}
		val, ok := res.(float64)
		if !ok {
			b.Fatalf("expected float64, got %T", res)
		}
		if val != expected {
			b.Fatalf("expected %f, got %f", expected, val)
		}
	}
}

func BenchmarkInterpreter_ExecuteSimpleAdd(b *testing.B) {
	interpreter := NewInterpreter(nil, nil)
	var expected float64 = 2
	for i := 0; i < b.N; i++ {
		res, err := interpreter.Execute(`1+1`)
		if err != nil {
			b.Fatalf("expected nil error, got %v", err)
		}
		val, ok := res.(float64)
		if !ok {
			b.Fatalf("expected float64, got %T", res)
		}
		if val != expected {
			b.Fatalf("expected %f, got %f", expected, val)
		}
	}
}

func BenchmarkInterpreter_ExecuteBig(b *testing.B) {
	interpreter := NewInterpreter(
		map[string]any{
			"X": 10.0,
			"Y": 100.0,
		},
		map[string]Func{
			"SUM": func(args ...any) (any, error) {
				var total float64
				for _, el := range args {
					val, ok := el.(float64)
					if !ok {
						return nil, fmt.Errorf("expected float64, got %T", el)
					}
					total += val
				}
				return total, nil
			},
			"And": func(args ...any) (any, error) {
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
			},
			"IF": func(args ...any) (any, error) {
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
			},
		},
	)
	var expected float64 = 123
	for i := 0; i < b.N; i++ {
		res, err := interpreter.Execute(`IF(AND(SUM(1;2;3)=6;X^2=Y);123;0)`)
		if err != nil {
			b.Fatalf("expected nil error, got %v", err)
		}
		val, ok := res.(float64)
		if !ok {
			b.Fatalf("expected float64, got %T", res)
		}
		if val != expected {
			b.Fatalf("expected %f, got %f", expected, val)
		}
	}
}
