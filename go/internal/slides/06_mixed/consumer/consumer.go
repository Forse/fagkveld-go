package main

import (
	mixed "fagkveld/internal/slides/06_mixed"
	"fmt"
)

func main() {
	a := mixed.NewOperation(1, 2)

	b := a.Add()
	fmt.Println(b)
}
