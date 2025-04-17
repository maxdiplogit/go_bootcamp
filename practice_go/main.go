package main

import (
	"fmt"

	"practice_go/simple"
)

func main() {
	var x, y int

	fmt.Print("Enter x: ")
	fmt.Scanf("%d", &x)

	fmt.Print("Enter y: ")
	fmt.Scanf("%d", &y)

	res := simple.Add(x, y)
	fmt.Printf("Result: %v\n", res)
}
