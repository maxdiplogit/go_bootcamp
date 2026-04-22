package slices

import (
	"fmt"
)

func SlicesPractice() {
	slice :=  []int{1, 2, 3}

	slice = append(slice, 4)

	fmt.Printf("slice: %#v\nslice capacity: %d\nslice length: %d\n", slice, cap(slice), len(slice))

	slicing := slice[0:2]	// slicing shares memory with slice
	slicing[0] = 55

	fmt.Printf("slicing: %v\nslice: %v\n", slicing, slice)
}

func SliceToArray() {
	slice_0 := make([]int, 5, 10)
	slice_1 := []int{1, 2, 3, 4, 5}

	fmt.Printf("slice_0: %#v\n", slice_0)
	fmt.Printf("slice_1: %#v\n", slice_1)

	// Conversion to an Array Pointer
	arrPtr_0 := (*[3]int)(slice_0)	// pointer to first 3 elements as [3]int
	arrPtr_1 := (*[4]int)(slice_1)	// [1 2 3]

	
	fmt.Printf("arrPtr_0: %#v\n", arrPtr_0)	
	fmt.Printf("arrPtr_1: %#v\n", arrPtr_1)

	// Conversion to an Array Value
	arrVal_0 := [3]int(slice_0)	// copy first 3 elements into a new [3]int
	arrVal_0[2] = 99		// slice_0 is unchanged

	
	fmt.Printf("arrVal_0: %#v\n", arrVal_0)	
}
