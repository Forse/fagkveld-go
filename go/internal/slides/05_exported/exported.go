package exported

//lint:ignore U1000 Example
type operation struct {
	a int
	b int
}

func NewOperation(a int, b int) *operation {
	return &operation{
		a: a,
		b: b,
	}
}
