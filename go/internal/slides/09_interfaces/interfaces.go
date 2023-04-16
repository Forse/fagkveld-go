package valuereceivers

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

// --------------------------------
// En annen package
//
//lint:ignore U1000 Example
func main() {
	op := NewOperation(2, 1) // *operation
	ComputeSumAndDiff(op, op)
}

// Interface type arguments er alltid pointers
func ComputeSumAndDiff(o1 Addition, o2 Subtraction) (int, int) {
	return o1.Add(), o2.Subtract()
}

var _ Addition = (*operation)(nil)
var _ Subtraction = (*operation)(nil)
