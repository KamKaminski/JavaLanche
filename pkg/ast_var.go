package javalanche

import "fmt"

var (
	_ Node           = (*Variable)(nil)
	_ SetValuer      = (*Variable)(nil)
	_ fmt.Stringer   = (*Variable)(nil)
	_ fmt.GoStringer = (*Variable)(nil)
)

type Variable struct {
	Name string
}

type SetValuer interface {
	SetValue(ctx *Javalanche, n Value) error
}

func NewVariable(n string) *Variable {
	return &Variable{Name: n}
}

// evaluates the variable node by getting its value from the evaluator
func (v *Variable) Eval(ctx *Javalanche) (Value, error) {
	return ctx.GetValue(v.Name)
}

// sets the value of the variable in the evaluator.
func (v *Variable) SetValue(ctx *Javalanche, n Value) error {
	return ctx.SetValue(v.Name, n)
}

func (v Variable) String() string {
	return v.Name
}

func (v Variable) GoString() string {
	return fmt.Sprintf("NewVariable(%q)", v.Name)
}
