package main

import "fmt"

func main() {
	a := 3
	b := 4
	b = b + 1
	a += b
	fmt.Println(a)
}
