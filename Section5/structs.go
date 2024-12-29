package main

import (
	"fmt"

	"example.com/structs/user"
)

type str string

func main() {
	userFirstName := getUserData("Enter First Name: ")
	userLastName := getUserData("Enter Last Name: ")
	userBirthDate := getUserData(("Enter Birth Date: "))

	var appUser *user.User
	appUser, err1 := user.New(userFirstName, userLastName, userBirthDate)

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	var appAdmin *user.Admin
	appAdmin, err2 := user.NewAdmin("admin@gmail.com", "1234")

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	appAdmin.OutputUserData()
	appAdmin.ClearUserName()
	appAdmin.OutputUserData()

	appUser.OutputUserData()
	appUser.ClearUserName()
	appUser.OutputUserData()

	var name str = "Max"

	name.log()
}

func getUserData(promptText string) string {
	fmt.Print(promptText)
	var value string
	fmt.Scanln(&value)
	return value
}

func (text str) log() {
	fmt.Println(text)
}
