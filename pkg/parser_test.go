package javalanche

import (
	"fmt"
	"testing"
)

func TestParserWithContext(t *testing.T) {
	type testCase struct {
		exprs  []string // strings to evaluate in order
		result Value    // expected Value for the last expression. nil if we expect an error
	}

	var cases = []testCase{

		{
			exprs:  []string{"x = 10\n", "y =  x + 7\n", "\nx * y"},
			result: NewInteger(170),
		},
		{
			exprs:  []string{"x = 10\n", "y =  (x + 7)\n", "\nx * y"},
			result: NewInteger(170),
		},
		{
			exprs:  []string{"x = 20\n", "y =  x + 4\n", "\nx * y"},
			result: NewInteger(480),
		},

		{
			exprs:  []string{`"Hello" + " " + "world"`},
			result: NewString("Hello world"),
		},
		{
			exprs:  []string{"quicmaths = 10 + 9", "quicmaths + 2"},
			result: NewInteger(21),
		},
		{
			exprs:  []string{"x = \"Hello\"", "x + \" world\""},
			result: NewString("Hello world"),
		},
		{
			exprs:  []string{"x = 10"},
			result: NewInteger(10),
		},
		{
			exprs:  []string{"y = 5"},
			result: NewInteger(5),
		},
		{
			exprs:  []string{"x = 10", "x * 2"},
			result: NewInteger(20),
		},
		{
			exprs:  []string{"5^-3"},
			result: NewFloat(0.008),
		},
		{
			exprs:  []string{"3^0"},
			result: NewInteger(1),
		},
		{
			exprs:  []string{`"foo" + "bar"`},
			result: NewString("foobar"),
		},
		{
			exprs:  []string{"1.0 == 1"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"1.0 != 1"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"-5 < -3"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"-5.0 < -3.0"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"true == 1"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"1+1"},
			result: NewInteger(2),
		},
		{
			exprs:  []string{"3>=3"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"6<=4"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"5 / 0"},
			result: nil,
		},
		{
			exprs:  []string{"5.0 / 0.0"},
			result: nil,
		},
		{
			exprs:  []string{"true + false"},
			result: nil,
		},
		{
			exprs:  []string{"true * false"},
			result: nil,
		},
		{
			exprs:  []string{"-1.5 + 2"},
			result: NewFloat(0.5),
		},
		{
			exprs:  []string{"1 == 1"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"(5 >= 5)"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"2 ** 3"},
			result: nil,
		},
		{
			exprs:  []string{"1 != 1"},
			result: NewBoolean(false),
		},
		{
			exprs:  []string{"(5 <= 5)"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"3.14159 * 2"},
			result: NewFloat(6.28318),
		},
		{
			exprs:  []string{"3.14159 / 2"},
			result: NewFloat(1.570795),
		},
		{
			exprs:  []string{"0.1 == 0.1"},
			result: NewBoolean(true),
		},
		{
			exprs:  []string{"0.1 != 0.2"},
			result: NewBoolean(true),
		},
	}

	for _, tc := range cases {
		evaluator := NewEvaluator()
		var res Value
		var err error

		for _, expr := range tc.exprs {
			res, err = EvalString(expr, evaluator)
			fmt.Printf("Evaluated expression %q, result is %v\n", expr, res)

			if err != nil {
				break
			}
		}

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
