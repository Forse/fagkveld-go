package main

import generics "fagkveld/internal/slides/10_generics"

func main() {
	op := generics.NewOperation(2, 1) // *operation[int]
	ComputeSumAndDiff[int](op)        // second type parameter is inferred
}

// Addition og Subtraction blir "embedded" i AdditionAndSubtraction
type AdditionAndSubtraction[T generics.Number] interface {
	generics.Addition[T]
	generics.Subtraction[T]
}

func ComputeSumAndDiff[T generics.Number, O AdditionAndSubtraction[T]](o O) (T, T) {
	return o.Add(), o.Subtract()
}
