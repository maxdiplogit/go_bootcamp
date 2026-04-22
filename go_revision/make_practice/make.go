package make_practice

import (
	"fmt"
)

func MakeExample() {
	slice := make([]int, 2, 5)
	fmt.Printf("make slice: %#v\n", slice)

	// slice[4] = 98

	slice[1] = 99
	fmt.Printf("make slice: %#v\n", slice)
}
