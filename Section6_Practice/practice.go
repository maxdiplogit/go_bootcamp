package main

import "fmt"

type Product struct {
	title string
	id    int
	price float64
}

func main() {
	var hobbies [3]string = [3]string{"Innovating", "Football", "WorkingOut"}

	var x = hobbies[:2]
	fmt.Println(x, cap(x))

	x = x[1:3]
	fmt.Println(x, cap(x))

	var courseGoals []string = []string{"LearnGO", "Expertise in GO"}
	updatedCourseGoals := append(courseGoals, "Make backend in GO")
	fmt.Println(updatedCourseGoals, cap(updatedCourseGoals))

	courseGoals[1] = "Goal Changed"
	fmt.Println(courseGoals, updatedCourseGoals)

	var productsList []Product = []Product{
		{
			"Product1",
			0,
			0.99,
		},
		{
			"Product2",
			1,
			1.99,
		},
	}
	fmt.Println(productsList)
	newProduct := Product{
		"Product3",
		2,
		2.99,
	}
	updatedProductsList := append(productsList, newProduct)

	fmt.Println(updatedProductsList)
}
