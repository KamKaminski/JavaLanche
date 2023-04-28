package main

import "testing"

func TestParser(t *testing.T) {
	var cases = []struct {
		expr   string
		result Value
	}{
		{"1 + 1", &IntegerLiteral{Value: 1}},
	}

	for _, tc := range cases {
		res, err := EvalString(tc.expr)

		switch {
		case err != nil && tc.result == nil:
			// expected fail
			t.Logf("%q: failed as expected: %s", tc.expr, err)
		case err == nil && tc.result == nil:
			// failed to fail
			t.Errorf("%q: should have failed, got %q instead", res)
		case err != nil && tc.result != nil:
			// unexpected error
			t.Errorf("%q: was expected to return %q: %s",
				tc.expr,
				tc.result,
				err)
		case tc.result.Equal(res):
			t.Log("%q: %q", tc.expr, res)
		default:
			t.Errorf("%q: got %q expected %q", tc.expr, tc.result, res)
		}
	}
}
