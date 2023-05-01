package main

import interfaces "fagkveld/internal/slides/09_interfaces"

func main() {
	op := interfaces.NewOperation(2, 1) // *operation
	ComputeSumAndDiff(op, op)
}

// Interface type arguments er alltid pointers
func ComputeSumAndDiff(o1 interfaces.Addition, o2 interfaces.Subtraction) (int, int) {
	return o1.Add(), o2.Subtract()
}
