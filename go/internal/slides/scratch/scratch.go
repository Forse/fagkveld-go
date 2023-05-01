package main

import "fmt"

func f(c chan string) {
	c <- "ping"
}
func main() {
	c := make(chan string, 4)
	f(c)
	f(c)
	f(c)
	f(c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}
