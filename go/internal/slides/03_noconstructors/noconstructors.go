package noconstructors

//lint:ignore U1000 Example
type operation struct {
	a int
	b int
}

//lint:ignore U1000 Example
func newOperation(a int, b int) *operation {
	return &operation{
		a: a,
		b: b,
	}
}
