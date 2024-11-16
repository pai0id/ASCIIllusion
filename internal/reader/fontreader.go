package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

func ReadCharsTxt(filename string) ([]fontparser.Char, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var chars []fontparser.Char
	for scanner.Scan() {
		text := scanner.Text()
		for _, ch := range text {
			char := fontparser.Char(ch)
			if char >= 32 && char <= 127 {
				chars = append(chars, char)
			} else {
				return nil, fmt.Errorf("error: cannot read file: non-ASCII char")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error: cannot read file: %w", err)
	}

	return chars, nil
}

type JsonData struct {
	Data []string `json:"data"`
}

func ReadCharsJson(filename string) ([]fontparser.Char, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error: cannot open file: %w", err)
	}
	defer file.Close()

	var data JsonData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error: cannot read file: %w", err)
	}

	slice := make([]fontparser.Char, len(data.Data))
	for i, ch := range data.Data {
		slice[i] = fontparser.Char(ch[0])
	}

	return slice, nil
}
