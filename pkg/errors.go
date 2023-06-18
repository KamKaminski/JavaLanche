package javalanche

import "fmt"

var (
	_ error = (*ErrInvalidToken)(nil)
	_ error = (*ErrInvalidValue)(nil)
)

type ErrInvalidToken struct {
	Token  *Token
	Reason string
}

func (e ErrInvalidToken) Error() string {
	return fmt.Sprintf("InvalidToken: %#v, Reason: %s", e.Token, e.Reason)
}

type ErrInvalidValue struct {
	Value string
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("InvalidValue: %q", e.Value)
}
