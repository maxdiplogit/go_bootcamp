package arrays

import (
	"fmt"
)

func ArrayPractice() {
	var a [3]int
	a[0] = 10

	fmt.Printf("a: %#v\n", a)
	
	b := [3]int{1, 2, 3}

	
	fmt.Printf("b: %v\n", b)

	c := b
	
	fmt.Printf("c: %v\n", c)

	c[0] = 99
	
	fmt.Printf("c: %#v\n", c)

	d := [...]int{1, 2, 3, 4}	// compiler counts → [4]int
	fmt.Printf("d: %#v\n", d)

}
