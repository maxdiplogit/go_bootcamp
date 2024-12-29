package main

import (
	"fmt"

	"example.com/bank/fileops"
)

const accountBalanceFileName = "balance.txt"

func main() {

	var accountBalance, err = fileops.GetFloatFromFile(accountBalanceFileName)

	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
	}

	fmt.Println("Welcome to Go Bank!")

	for {
		presentOptions()

		var choice int
		fmt.Print("Enter choice: ")
		fmt.Scan(&choice)

		fmt.Println("Your choice:", choice)

		switch choice {
		case 1:
			fmt.Println("Account Balance:", accountBalance)
		case 2:
			var depositAmount float64
			fmt.Print("Enter Deposit Amount: ")
			fmt.Scan(&depositAmount)

			accountBalance += depositAmount
			fmt.Println("Updated Account Balance:", accountBalance)
			fileops.WriteFloatToFile(accountBalanceFileName, accountBalance)
		case 3:
			var withdrawAmount float64
			fmt.Print("Enter Withdrawal Amount: ")
			fmt.Scan(&withdrawAmount)

			accountBalance -= withdrawAmount
			fmt.Println("Updated Account Balance:", accountBalance)
			fileops.WriteFloatToFile(accountBalanceFileName, accountBalance)
		default:
			return
		}
	}
}
