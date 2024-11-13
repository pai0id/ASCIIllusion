package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/FontParser"
)

func main() {
	fontMap, _ := FontParser.GetFontMap("Ubuntu.ttf", 20, 20, 20, 10)

	fmt.Println(len(fontMap))
}
