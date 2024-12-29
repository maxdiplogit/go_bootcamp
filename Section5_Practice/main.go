package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/notes/note"
	"example.com/notes/todo"
)

type saver interface {
	Save() error
}

type outputtable interface {
	saver
	Display()
}

func main() {
	x := add(1, 2)

	fmt.Println(x)

	title, content := getNoteData()
	todoText := getUserInput("Todo Text: ")

	var newNote *note.Note
	newNote, err := note.New(title, content)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = outputData(newNote)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Note saved successfully!")

	newTodo, err := todo.New(todoText)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = outputData(newTodo)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("New Todo saved successfully!")
}

func getNoteData() (string, string) {
	title := getUserInput("Enter Title of Note: ")

	content := getUserInput("Enter Note Content: ")

	return title, content
}

func printSomething(value any) {
	intVal, ok := value.(int)

	if ok {
		intVal = intVal + 1
		fmt.Println(intVal)
	}
}

func add[T int | float64 | string](a, b T) T {
	return a + b
}

func getUserInput(promptText string) string {
	fmt.Print(promptText)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return ""
	}

	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")

	return text
}

func outputData(data outputtable) error {
	data.Display()
	err := data.Save()

	if err != nil {
		return err
	}

	return nil
}

func saveData(data saver) error {
	err := data.Save()

	if err != nil {
		fmt.Println("Unable to save ", data)
		return err
	}

	return nil
}
