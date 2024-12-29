package main

import "fmt"

type floatMap map[string]float64

func main() {
	courseRatings := make(floatMap, 3)

	courseRatings.output()
}

func (m floatMap) output() {
	fmt.Println(m)
}
