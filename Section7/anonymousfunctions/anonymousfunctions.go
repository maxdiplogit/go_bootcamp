package anonymousfunctions

import "fmt"

type transformFn func(int) int

func main() {
	numbers := []int{1, 2, 3, 4}

	doubleTransformerFunction := createTransformerFunction(2)

	transformed := transformNumbers(&numbers, doubleTransformerFunction)

	fmt.Println(transformed)
}

func transformNumbers(numbers *[]int, transform transformFn) []int {
	dNumbers := []int{}
	for _, value := range *numbers {
		dNumbers = append(dNumbers, transform(value))
	}

	return dNumbers
}

func createTransformerFunction(factor int) transformFn {
	return func(number int) int {
		return number * factor
	}
}
