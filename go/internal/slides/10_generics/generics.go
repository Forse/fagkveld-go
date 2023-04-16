package generics

// This package is 0.0.0 - experimental
import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type operation[T Number] struct {
	a T
	b T
}

func NewOperation[T Number](a T, b T) *operation[T] {
	return &operation[T]{
		a: a,
		b: b,
	}
}

type Addition[T Number] interface {
	Add() T
}
type Subtraction[T Number] interface {
	Subtract() T
}

func (o *operation[T]) Add() T {
	return o.a + o.b
}
func (o *operation[T]) Subtract() T {
	return o.a - o.b
}

var _ Addition[int] = (*operation[int])(nil)
var _ Subtraction[int] = (*operation[int])(nil)

// --------------------------------
// En annen package
//
//lint:ignore U1000 Example
func main() {
	op := NewOperation(2, 1)   // *operation
	ComputeSumAndDiff[int](op) // second type parameter is inferred
}

// Addition og Subtraction blir "embedded" i AdditionAndSubtraction
type AdditionAndSubtraction[T Number] interface {
	Addition[T]
	Subtraction[T]
}

func ComputeSumAndDiff[T Number, O AdditionAndSubtraction[T]](o O) (T, T) {
	return o.Add(), o.Subtract()
}
