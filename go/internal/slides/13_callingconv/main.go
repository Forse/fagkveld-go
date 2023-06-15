package main

import "fmt"

// Add three numbers
//
//go:noinline
func Add(x, y int) int {
	return x + y
}

// Entrypoint to illustrate calling convention
func main() {
	x := 1
	y := 2

	r := Add(x, y)

	fmt.Println(r)
}
