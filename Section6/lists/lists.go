package lists

import "fmt"

func main() {
	var prices []float64 = []float64{10.99, 0.99, 2.99}
	fmt.Println(prices, cap(prices))

	prices = append(prices, 1.99)
	fmt.Println(prices, cap(prices))

	prices = append(prices, 6.99)
	fmt.Println(prices, cap(prices))

	prices = prices[1:]
	fmt.Println(prices, cap(prices))
}
