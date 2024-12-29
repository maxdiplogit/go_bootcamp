package fileops

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func WriteFloatToFile(fileName string, value float64) {
	balanceText := fmt.Sprint(value)
	os.WriteFile(fileName, []byte(balanceText), 0644)
}

func GetFloatFromFile(fileName string) (float64, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return -1, errors.New("Failed to read/find file")
	}

	res, err := strconv.ParseFloat(string(data), 64)

	if err != nil {
		return -1, errors.New("Failed to parse stored value")
	}

	return res, nil
}
