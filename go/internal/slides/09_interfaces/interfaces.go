package interfaces

type operation struct {
	Addition
	a int
	b int
}

func NewOperation(a int, b int) *operation {
	return &operation{
		a: a,
		b: b,
	}
}

type Addition interface {
	Add() int
}
type Subtraction interface {
	Subtract() int
}

func (o *operation) Add() int {
	return o.a + o.b
}
func (o *operation) Subtract() int {
	return o.a - o.b
}

var _ Addition = (*operation)(nil)
var _ Subtraction = (*operation)(nil)
