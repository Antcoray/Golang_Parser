package main

import (
	"fmt"
	"math"
)

// 1 - operators, 2 - operands
func main() { // 1) func, main, { 2)
	fmt.Println("Pizza") // 1) ., Println, () 2) fmt, "Pizza"
	var x int = 15       // 1) var, = 2) x, 15
	var y int = 20       // 1) var, = 2) y, 20
	f := "dsfkdsfjk"     // 1) := 2) f, "dsfkdsfjk"
	fmt.Printf("%s", f)  // 1) ., Printf, (), , 2) fmt, "%s", f

	if x > y { // 1) if, >, { 2) x, y
		fmt.Println("x > y") // 1) ., Println, () 2) fmt, "x > y"
	} else if 6 < 5 { // 1) else, if, <, { 2) 6, 5
		fmt.Println("never") // 1) ., Println, () 2) fmt, "never"
	} else { // 1) else, {
		fmt.Println("x <= y") // 1) ., Println, () 2) fmt, "x <= y"
	} // 1) }

	if 6 != 5 { // 1) if, !=, { 2) 6, 5
		fmt.Println("aboba") // 1) ., Println, () 2) fmt, "aboba"
	} // 1) }

	for a := 0; a < 10; a++ { // 1) for, :=, <, ++, { 2) a, 0, 10
		for u := 5; u < 10; u++ { // 1) for, :=, <, ++, { 2) u, 5, 10
			fmt.Println(a * u) // 1) ., Println, (), * 2) fmt, a, u
			if 8 == 8 {        // 1) if, ==, { 2) 8, 8
				break // 1) break
			} // 1) }
		} // 1) }
	} // 1) }

	x++                         // 1) ++ 2) x
	x = x + y - 4*5             // 1) =, +, -, * 2) x, x, y, 4, 5
	fmt.Println(math.Acos(0.5)) // 1) ., Println, (), ., Acos, () 2) fmt, math, 0.5

	var a int = 4 // 1) var, = 2) a, 4
L1:
	a--        // 1) -- 2) a
	switch a { // 1) switch, { 2) a
	case 1: // 1) case, : 2) 1
		aboba := "aboba"           // 1) := 2) aboba, "aboba"
		fmt.Println([]byte(aboba)) // 1) ., Println, (), [] 2) fmt, aboba
		fallthrough                // 1) fallthrough
	case 2: // 1) case, : 2) 2
		{ // 1) {
			goto L1 // 1) goto 2)
		} // 1) }
	case 3: // 1) case, : 2) 3
		{
			goto L1 // 1) goto 2)
		}
	case 4: // 1) case, : 2) 4
		{
			goto L1 // 1) goto 2)
		}
	} // 1) }
	fmt.Println(a) // 1) ., Println, () 2) fmt, a

	var numbers [500]float64 = [500]float64{1, 2.0, 3.0, 4.0, 5.0} // 1) var, =, [], {}, , (запятые) 2) numbers, 1, 2.0, 3.0, 4.0, 5.0
	fmt.Println(len(numbers))                                      // 1) ., Println, (), len, () 2) fmt, numbers

	for _, value := range numbers { // 1) for, range, :=, { 2) _, value, numbers
		if value != 0 { // 1) if, !=, { 2) value, 0
			fmt.Println(value) // 1) ., Println, () 2) fmt, value
		} // 1) }
	} // 1) }

	fmt.Println(Add(int(numbers[0]), int(numbers[1]), int(numbers[5]))) // 1) ., Println, (), Add, (), [], , 2) fmt, numbers, 0, 1, 5
	fmt.Println(Multiply(1, 2, 3))                                      // 1) ., Println, (), Multiply, (), , 2) fmt, 1, 2, 3
	fmt.Println(Multiply(1, 2, 3, 4))                                   // 1) ., Println, (), Multiply, (), , 2) fmt, 1, 2, 3, 4

	var p *int = &x             // 1) var, *, =, & 2) p, x
	fmt.Println(*p&45*3/78 | 4) // 1) ., Println, (), *, &, |, *, / 2) fmt, p, 45, 3, 78, 4
	var b = 45 &^ 6 & 8932 << 4 // 1) var, =, &^, &, << 2) b, 45, 6, 8932, 4
	fmt.Println(b)              // 1) ., Println, () 2) fmt, b

	type chuvak struct { // 1) type, struct, { 2)
		Age  int
		Cool bool
		Cash int
		Name string
	} // 1) }
	var Chuvachok chuvak = chuvak{18, false, 0, "Durak"} // 1) var, =, {}, , 2) Chuvachok, 18, false, 0, "Durak"
	fmt.Println(Chuvachok)                               // 1) ., Println, () 2) fmt, Chuvachok
}

func Add(x1 int, x2 int, x3 int) int { // 1) func, Add, (, ), , { 2) x1, x2, x3
	return x1 + x2 + x3 // 1) return, +, + 2) x1, x2, x3
} // 1) }

func Multiply(numbers ...float64) float64 { // 1) func, Multiply, (, ), ..., , { 2) numbers
	var sum float64 = 1              // 1) var, = 2) sum, 1
	for _, number := range numbers { // 1) for, range, :=, { 2) _, number, numbers
		sum *= number // 1) *= 2) sum, number
	} // 1) }
	return sum // 1) return 2) sum
} // 1) }
