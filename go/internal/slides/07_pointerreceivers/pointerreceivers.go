package valuereceivers

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

func (o *operation) Add() int {
	return o.a + o.b
}
