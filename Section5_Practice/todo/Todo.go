package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Todo struct {
	Text string `json:"text`
}

func New(todoText string) (*Todo, error) {
	if todoText == "" {
		return nil, errors.New("Todo text cannot be empty!")
	}

	return &(Todo{
		Text: todoText,
	}), nil
}

func (todo *Todo) Display() {
	fmt.Println(todo.Text)
}

func (todo *Todo) Save() error {
	fileName := "todo.json"

	jsonContent, err := json.Marshal(todo)

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, jsonContent, 0644)
}
