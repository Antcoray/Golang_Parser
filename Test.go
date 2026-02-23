package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("Pizza")
	var x int = 15
	var y int = 20
	f := "dsfkdsfjk"
	fmt.Printf("%s", f)
	if x > y {
		fmt.Println("x > y")
	} else if 6 < 5 {
		fmt.Println("never")
	} else {
		fmt.Println("x <= y")
	}
	if 6 != 5 {
		fmt.Println("aboba")
	}
	for a := 0; a < 10; a++ {
		for u := 5; u < 10; u++ {
			fmt.Println(a * u)
			if 8 == 8 {
				break
			}
		}
	}
	x++
	x = x + y - 4*5
	fmt.Println(math.Acos(0.5))
	var a int = 4
L1:
	a--
	switch a {
	case 1:
		aboba := "aboba"
		fmt.Println([]byte(aboba))
		fallthrough
	case 2:
		{
			goto L1
		}
	case 3:
		{
			goto L1
		}
	case 4:
		{
			goto L1
		}
	}
	fmt.Println(a)
	var numbers [500]float64 = [500]float64{1, 2.0, 3.0, 4.0, 5.0}
	fmt.Println(len(numbers))
	for _, value := range numbers {
		if value != 0 {
			fmt.Println(value)
		}
	}
	fmt.Println(numbers)
	fmt.Println(Add(int(numbers[0]), int(numbers[1]), int(numbers[5])))
	fmt.Println(Multiply(1, 2, 3))
	fmt.Println(Multiply(1, 2, 3, 4))
	var p *int = &x
	fmt.Println(*p&45*3/78 | 4)
	var b = 45 &^ 6 & 8932 << 4
	fmt.Println(b)
	type chuvak struct {
		Age  int
		Cool bool
		Cash int
		Name string
	}
	var Chuvachok chuvak = chuvak{18, false, 0, "Durak"}
	fmt.Println(Chuvachok)
}
func Add(x1 int, x2 int, x3 int) int {
	return x1 + x2 + x3
}
func Multiply(numbers ...float64) float64 {
	var sum float64 = 1
	for _, number := range numbers {
		sum *= number
	}
	return sum
}
