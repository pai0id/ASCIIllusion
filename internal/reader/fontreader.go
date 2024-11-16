package reader

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pai0id/CgCourseProject/internal/fontparser"
)

func ReadChars(filename string) ([]fontparser.Char, error) {
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
			chars = append(chars, fontparser.Char(ch))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return chars, nil
}
