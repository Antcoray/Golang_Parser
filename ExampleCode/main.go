package main

import (
	"fmt"
)

// 1 - operators, 2 - operands
func main() { // 1) func, main, { 2)
	fmt.Println("Pizza")      // 1) ., Println, () 2) fmt, "Pizza"
	var x int = 15            // 1) var, = 2) x, 15
	var y int = 20            // 1) var, = 2) y, 20
	f := "dsfkdsfjk"          // 1) := 2) f, "dsfkdsfjk"
	fmt.Printf("%s", f, x, y) // 1) ., Printf, (), , 2) fmt, "%s", f
	type chuvak struct {
		Age  int
		Cool bool
		Cash int
		Name string
	}
	var Chuvachok chuvak = chuvak{18, false, 0, "Durak"}
	fmt.Printf("%s", Chuvachok)

	for a := 0; a < 10; a++ {
		for u := 5; u < 10; u++ {
			fmt.Println(a * u)
			if 8 == 8 {
				break
			}
		}
	}

}

func Multiply(numbers ...float64) float64 { // 1) func, Multiply, (, ), ..., , { 2) numbers
	var sum float64 = 1 // 1) var, = 2) sum, 1
	for _, number := range numbers {
		sum *= number // 1) *= 2) sum, number
	} // 1) }
	return sum // 1) return 2) sum
} // 1) }
