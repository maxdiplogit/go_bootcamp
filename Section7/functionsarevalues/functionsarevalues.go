package functionsarevalues

import "fmt"

type transformFn func(int) int

func main() {
	numbers := []int{1, 2, 3, 4}
	x := transformNumbers(&numbers, double)

	fmt.Println(x)
}

func transformNumbers(numbers *[]int, transform transformFn) []int {
	dNumbers := []int{}
	for _, value := range *numbers {
		dNumbers = append(dNumbers, transform(value))
	}

	return dNumbers
}

func getTransformerFunction() transformFn {
	return double
}

func double(number int) int {
	return number * 2
}
