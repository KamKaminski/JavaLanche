package javalanche

import (
	"strings"
	"testing"
)

func TestParserWithContext(t *testing.T) {
	type testCase struct {
		exprs  []string // strings to evaluate in order
		result Value    // expected Value for the last expression. nil if we expect an error
	}

	var cases = []testCase{
		{
			exprs:  []string{"3"},
			result: NewInteger(3),
		},

		{
			exprs:  []string{"3 - 3"},
			result: NewInteger(0),
		},
		{
			exprs:  []string{"3 + 3"},
			result: NewInteger(6),
		},
		{
			exprs:  []string{"3 * 3"},
			result: NewInteger(9),
		},
		{
			exprs:  []string{"3 * 3 - 2"},
			result: NewInteger(7),
		},
		{
			exprs:  []string{"-2"},
			result: NewInteger(-2),
		},
		{
			exprs:  []string{"-2 + 3"},
			result: NewInteger(1),
		},
		{
			exprs:  []string{"2 + -3"},
			result: NewInteger(-1),
		},
		{
			exprs:  []string{"2 - -3"},
			result: NewInteger(5),
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
			exprs:  []string{"x = 10", "x++", "x"},
			result: NewInteger(11),
		},
		{
			exprs:  []string{"\"a\"+ \"b\""},
			result: NewString("ab"),
		},
		{
			exprs:  []string{"true and true"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"true && true"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"x = 3", "x++", "x >= 4"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"true == true"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"true != true"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"true != false"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"2 + 3 * 4"},
			result: NewInteger(14),
		},
		{
			exprs:  []string{"2 + 3 == 5"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"2 + 3 != 5"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"2 + 3 < 5"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"2 +3 > 5"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"2 + 3 >= 5"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"2 + 3 >= 5"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"3 - 4 * 8"},
			result: NewInteger(-29),
		},
		{
			exprs:  []string{"8+3/3"},
			result: NewInteger(9),
		},
		{
			exprs:  []string{"x = 10", "x--", "x"},
			result: NewInteger(9),
		},
		{
			exprs:  []string{"x = 10", "y = x + 5", "x * y"},
			result: NewInteger(150),
		},
		{
			exprs:  []string{"!true"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"\"hello\"+ \"world\""},
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
		{
			exprs:  []string{"!true == !false"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"(123)"},
			result: NewInteger(123),
		},
		{
			exprs:  []string{"(1 + 2)"},
			result: NewInteger(3),
		},
		{
			exprs:  []string{"2 + (2 - 3)"},
			result: NewInteger(1),
		},
		{
			exprs:  []string{"(2 - 3) * 3"},
			result: NewInteger(-3),
		},
		{
			exprs:  []string{"(1 + 2) == 3"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"(4 * 8) - 3 * 3"},
			result: NewInteger(23),
		},
		{
			exprs:  []string{"4 - 1 < (3 * 4) - (5 + 3)"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"!(5 - 4 > 3 * 2 == !false)"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"x = 10", "if (x == 10)", "isTen = true", "else isTen = false", "end", "isTen"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"x = 10", "for (x < 25)", "x++", "end", "x"},
			result: NewInteger(25),
		},
		{
			exprs:  []string{"2%2"},
			result: NewInteger(0),
		},
		{
			exprs:  []string{"x=0", "for (x <= 100)", "x++", "x", "end", "x"},
			result: NewInteger(101),
		},
	}

	for _, tc := range cases {
		var retrying bool

	retryLoop:
		for {
			var failed bool

			ctx := New()

			exprs := strings.Join(tc.exprs, "\n")
			res, err := ctx.EvalLine(tc.exprs...)

			switch {
			case err != nil && tc.result == nil:
				t.Logf("PASS: %q: failed as expected: %s", exprs, err)
			case err == nil && tc.result == nil:
				failed = true
				t.Errorf("ERROR: %q: should have failed, got %q instead", exprs, res)
			case err != nil && tc.result != nil:
				failed = true
				t.Errorf("ERROR: %q: was expected to return %q: %s", exprs, tc.result, err)
			case tc.result.Equal(res):
				t.Logf("PASS: %q → %q", exprs, res)
			default:
				failed = true
				t.Errorf("ERROR: %q: got %q expected %q", exprs, res, tc.result)
			}

			switch {
			case !failed:
				// PASS
				debugStage = false
				retrying = false
				break retryLoop
			case retrying:
				// not again
				debugStage = false
				break retryLoop
			default:
				// oops, again but with logs
				debugStage = true
				retrying = true
			}
		}
	}
}
