package go_interpreter

import (
	"github.com/kovalenkong/go-interpreter/functions"
	"strings"
	"testing"
)

func prepareExecutor(vars map[string]any) *Interpreter {
	return &Interpreter{
		variables: vars,
		functions: map[string]Func{
			"Sum":  functions.Sum,
			"Len":  functions.Len,
			"Mean": functions.Mean,
			"And":  functions.And,
			"Or":   functions.Or,
		},
	}
}

func TestParser_ParseFloat(t *testing.T) {
	var X, Y float64 = 3, 10
	executor := prepareExecutor(map[string]any{
		"X": X,
		"Y": Y,
	})

	cases := map[string]float64{
		`1 + 1`:               2,
		`2 + 2 * 2`:           6,
		`2 + (2 * 2)`:         6,
		`(2 + 2) * 2`:         8,
		`(2 + 2 * 2)`:         6,
		`((2 + 2) * 2)`:       8,
		`20 / 4`:              5,
		`Y / 10`:              Y / 10,
		`Y * X + 2`:           Y*X + 2,
		`2,5 + 2,5`:           5,
		`2 + 2 * 2 - 2 / 2`:   5,
		`Len(1;2;3)`:          3,
		`Len()`:               0,
		`2 * Len(1;2;3) * 3`:  18,
		`Sum()`:               0,
		`1 + Sum() + 1`:       2,
		`Sum(1;1;1)`:          3,
		`Sum(1;1;Sum(1;1;1))`: 5,
		`Len(1;2;Sum(1;2;3))`: 3,
		`Sum(1;2;Len(1;1;1))`: 6,
		`-1 + 2`:              1,
		`--1`:                 1,
		`-1 + (2)`:            1,
		`-1 + (-2)`:           -3,
	}
	for formula, result := range cases {
		lexer := NewLexer()
		tokens, err := lexer.Lex(strings.NewReader(formula))
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		parser := NewParser()
		node, err := parser.Parse(tokens)
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		res, err := executor.execute(node)
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		if res != result {
			t.Fatalf("formula '%s' expected '%f', got '%f'", formula, result, res)
		}
	}
}

func TestParser_ParseBool(t *testing.T) {
	var X, Y float64 = 3, 10
	executor := prepareExecutor(map[string]any{
		"X": X,
		"Y": Y,
	})
	cases := map[string]bool{
		`1 = 1`:  true,
		`1 <= 1`: true,
		`1 <= 2`: true,
		`1 < 2`:  true,
		`1 < 1`:  false,
		`1 <= 0`: false,
	}
	for formula, result := range cases {
		lexer := NewLexer()
		tokens, err := lexer.Lex(strings.NewReader(formula))
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		parser := NewParser()
		node, err := parser.Parse(tokens)
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		res, err := executor.execute(node)
		if err != nil {
			t.Fatalf("formula '%s': expected nil error, got %s", formula, err)
		}
		if res != result {
			t.Fatalf("formula '%s' expected '%v', got '%f'", formula, result, res)
		}
	}
}
