package javalanche

var (
	_ Node      = (*Variable)(nil)
	_ SetValuer = (*Variable)(nil)
)

type Variable struct {
	Name string
}

type SetValuer interface {
	SetValue(ctx *Javalanche, n Value) error
}

// evaluates the variable node by getting its value from the evaluator
func (v *Variable) Eval(ctx *Javalanche) (Value, error) {
	return ctx.GetValue(v.Name)
}

// sets the value of the variable in the evaluator.
func (v *Variable) SetValue(ctx *Javalanche, n Value) error {
	return ctx.SetValue(v.Name, n)
}
