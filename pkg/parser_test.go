package javalanche

import "testing"

func TestParser(t *testing.T) {
	var cases = []struct {
		expr   string // string to evaluate
		result Value  // expected Value. nil if we expect an error
	}{
		{"1 + 1", NewInteger(2)},
		{"\"\"", NewString("")},
		{"\"hello\" + \"world\"", NewString("helloworld")},
		{"let x = 5\nx * 2", NewInteger(10)},
		{"true == false", NewBoolean(false)},
		{"true != false", NewBoolean(true)},
		{"(5 < 10)", NewBoolean(true)},
		{"!(5 - 4 > 3 * 2 == !false)", NewBoolean((true))},
		{"true and true", NewBoolean(true)},
		{"false and true", NewBoolean(false)},
		{("0<1 or false"), NewBoolean(true)},
		{"false or false", NewBoolean(false)},
		{"\"foo\" + \"bar\"", NewString("foobar")},
		{"\"10 corgis\" != \"10\" +\"corgis\" ", NewBoolean(true)},
		{"let quickMaths = 9 + 10", NewInteger(19)},
	}

	for _, tc := range cases {
		res, err := EvalString(tc.expr)

		switch {
		case err != nil && tc.result == nil:
			// expected fail
			t.Logf("PASS: %s: failed as expected: %s", tc.expr, err)
		case err == nil && tc.result == nil:
			// failed to fail
			t.Errorf("ERROR: %s: should have failed, got %s instead", tc.expr, res)
		case err != nil && tc.result != nil:
			// unexpected error
			t.Errorf("ERROR: %s: was expected to return %s: %s",
				tc.expr,
				tc.result,
				err)
		case tc.result.Equal(res):
			t.Logf("PASS: %s → %s", tc.expr, res)
		default:
			t.Errorf("ERROR: %s: got %s expected %s", tc.expr, res, tc.result)
		}
	}
}
