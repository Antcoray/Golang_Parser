package main

import (
	"fmt"
	"math"
)

func Test() {
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
		fmt.Println(a)
	}
	x++
	x = x + y - 4*5
	fmt.Println(math.Acos(0.5))
	var a int = 4
L1:
	a--
	switch a {
	case 1:
		{
			aboba := "aboba"
			fmt.Println([]byte(aboba))
		}
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
}
