package javalanche

import "testing"

func TestParser(t *testing.T) {
	var cases = []struct {
		expr   string // string to evaluate
		result Value  // expected Value. nil if we expect an error
	}{
		{"\"foo\" + \"bar\"", NewString("foobar")},                 // String concatenation
		{"1.0 == 1", NewBoolean(true)},                             // Float and integer equality
		{"1.0 != 1", NewBoolean(false)},                            // Float and integer inequality
		{"-5 < -3", NewBoolean(true)},                              // Negative integer comparison
		{"-5.0 < -3.0", NewBoolean(true)},                          // Negative float comparison
		{"true == 1", NewBoolean(false)},                           // Boolean and integer comparison
		{"1+1", NewInteger(2)},                                     // Integer addition
		{"3>=3", NewBoolean(true)},                                 // Integer greater or equal comparison
		{"6<=4", NewBoolean(false)},                                // Integer less or equal comparison
		{"5 / 0", nil},                                             // Integer division by zero (error case)
		{"5.0 / 0.0", nil},                                         // Float division by zero (error case)
		{"true + false", nil},                                      // Invalid addition between booleans (error case)
		{"true * false", nil},                                      // Invalid multiplication between booleans (error case)
		{"-1.5 + 2", NewFloat(0.5)},                                // Float addition with negative number
		{"1 == 1", NewBoolean(true)},                               // Integer equality
		{"(5 >= 5)", NewBoolean(true)},                             // Integer greater or equal comparison
		{"2 ** 3", nil},                                            // Invalid exponentiation operation (error case)
		{"1 != 1", NewBoolean(false)},                              // Integer inequality
		{"(5 <= 5)", NewBoolean(true)},                             // Integer less or equal comparison
		{"3.14159 * 2", NewFloat(6.28318)},                         // Float multiplication
		{"3.14159 / 2", NewFloat(1.570795)},                        // Float division
		{"0.1 == 0.1", NewBoolean(true)},                           // Float equality
		{"0.1 != 0.2", NewBoolean(true)},                           // Float inequality
		{"0.1 != 0.1", NewBoolean(false)},                          // Float inequality (false case)
		{"0.5 + 0.1 == 0.6", NewBoolean(true)},                     // Float addition and equality
		{"true ^ true", NewBoolean(false)},                         // Boolean XOR
		{"false ^ true", NewBoolean(true)},                         // Boolean XOR
		{"1 + 1", NewInteger(2)},                                   // Integer addition
		{"4 + 5", NewInteger(9)},                                   // Integer addition
		{"8/(8/8)", NewInteger(8)},                                 // Integer division
		{"3 != 3", NewBoolean(false)},                              // Integer inequality (false case)
		{"\"\"", NewString("")},                                    // Empty string
		{"1.5 + 3.5", NewInteger(5)},                               // Float addition
		{"\"hello\" + \"world\"", NewString("helloworld")},         // String concatenation
		{"true == false", NewBoolean(false)},                       // Boolean equality
		{"true != false", NewBoolean(true)},                        // Boolean inequality
		{"(5 < 10)", NewBoolean(true)},                             // Integer less than comparison
		{"!(5 - 4 > 3 * 2 == !false)", NewBoolean(true)},           // Complex boolean expression
		{"true and true", NewBoolean(true)},                        // Boolean AND operation
		{"false and true", NewBoolean(false)},                      // Boolean AND operation (false case)
		{"0<1 or false", NewBoolean(true)},                         // Boolean OR operation
		{"false or false", NewBoolean(false)},                      // Boolean OR operation (false case)
		{"\"foo\" + \"bar\"", NewString("foobar")},                 // String concatenation
		{"\"10 corgis\" != \"10\" +\"corgis\" ", NewBoolean(true)}, // String concatenation and inequality
		{"-1 + 2", NewInteger(1)},                                  // Integer addition with negative number
	}

	for _, tc := range cases {
		res, err := EvalString(tc.expr)

		switch {
		case err != nil && tc.result == nil:
			// expected fail
			t.Logf("PASS: %s: failed as expected: %s", tc.expr, err)
		case err == nil && tc.result == nil:
			// failed to fail
			t.Errorf("ERROR: %s: should have failed, got %q instead", tc.expr, res)
		case err != nil && tc.result != nil:
			// unexpected error
			t.Errorf("ERROR: %s: was expected to return %q: %s",
				tc.expr,
				tc.result,
				err)
		case tc.result.Equal(res):
			t.Logf("PASS: %s → %q", tc.expr, res)
		default:
			t.Errorf("ERROR: %s: got %q expected %q", tc.expr, res, tc.result)
		}
	}
}
