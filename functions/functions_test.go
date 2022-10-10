package functions

import "testing"

func TestLen(t *testing.T) {
	cases := map[float64][]any{
		0: {},
		2: {1, "s"},
		1: {"s"},
		3: {"s", "s", "s"},
	}
	for result, args := range cases {
		res, err := Len(args...)
		if err != nil {
			t.Errorf("expected nil error, got '%v'", err)
			continue
		}
		if result != res {
			t.Errorf("expected %f, got %f", result, res)
		}
	}
}

func TestSum(t *testing.T) {
	type Case struct {
		result float64
		args   []any
	}
	cases := []Case{
		{0, []any{0.}},
		{0, []any{}},
		{1, []any{1.}},
		{1, []any{1., 0., 0.}},
		{6, []any{1., 2., 3.}},
	}
	for _, testCase := range cases {
		res, err := Sum(testCase.args...)
		if err != nil {
			t.Errorf("expected nil error, got '%v'", err)
			continue
		}
		if res != testCase.result {
			t.Errorf("expected %f, got %f", testCase.result, res)
		}
	}
}
