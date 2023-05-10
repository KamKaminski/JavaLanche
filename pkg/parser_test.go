package javalanche

import (
	"testing"
)

func TestParserWithContext(t *testing.T) {
	type testCase struct {
		exprs  []string // strings to evaluate in order
		result Value    // expected Value for the last expression. nil if we expect an error
	}

	var cases = []testCase{

		{
			exprs:  []string{"x = 10", "y = x + 5", "x * y"},
			result: NewInteger(150),
		},
		{
			exprs:  []string{"x = 3", "y = 4", "x * y"},
			result: NewInteger(12),
		},
		{
			exprs:  []string{"x = 7", "y = 8", "y+x"},
			result: NewInteger(15),
		},
		{
			exprs:  []string{"x = 7", "y = x + 8", "y * x"},
			result: NewInteger(105),
		},
		{
			exprs:  []string{"true and true == true"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"!(5 - 4 > 3 * 2 == !false)"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"3 == 3"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"3 <= 3"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"4 > 5"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"\"hello\"+ \"world"},
			result: NewString("helloworld"),
		},
		{
			exprs:  []string{"x = 2", "y = 3", "x * y + 4 <= 12"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"x = 2", "y = 3", "x * y + 4 > 12"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"x = 3", "y = 4", "x * y / 2 + 1"},
			result: NewInteger(7),
		},
		{
			exprs:  []string{"x = 10", "y = 2", "x ^ y"},
			result: NewInteger(100),
		},
		{
			exprs:  []string{"x = 10", "y = 5", "x / y"},
			result: NewInteger(2),
		},
		{
			exprs:  []string{"x = true", "y = false", "x and y"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"x = true", "y = false", "x or y"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"x = \"hello\"", "y = \"world\"", "x + y"},
			result: NewString("helloworld"),
		},
		{
			exprs:  []string{"x = 2", "y = 3", "x * y > 6 and x + y < 7"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"x = 2", "y = 3", "x * y > 6 or x + y < 7"},
			result: NewBoolean(true),
		},
	}

	for _, tc := range cases {
		evaluator := NewEvaluator()
		res, err := EvalString(evaluator, tc.exprs...)

		switch {
		case err != nil && tc.result == nil:
			t.Logf("PASS: %v: failed as expected: %s", tc.exprs, err)
		case err == nil && tc.result == nil:
			t.Errorf("ERROR: %v: should have failed, got %q instead", tc.exprs, res)
		case err != nil && tc.result != nil:
			t.Errorf("ERROR: %v: was expected to return %q: %s", tc.exprs, tc.result, err)
		case tc.result.Equal(res):
			t.Logf("PASS: %v → %q", tc.exprs, res)
		default:
			t.Errorf("ERROR: %v: got %q expected %q", tc.exprs, res, tc.result)
		}
	}
}
