package utils

import (
	"fmt"
	"time"
)

func PrintDay() {
	switch time.Now().Weekday() {
		case time.Saturday, time.Sunday:
			fmt.Println("It's weekend!")
		default:
			fmt.Println("It's a weekday...")
	}
}
