package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/asciiser/mapping"
	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/reader"
)

func main() {
	mctx := mapping.NewContext(8, 8, 2, 1, 1)
	chars, err := reader.ReadCharsTxt("fonts/slice.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	f, err := fontparser.GetFontMap("fonts/IBM.ttf", 8, 8, 6, 72, chars)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	dctx := asciiser.NewDrawContext(mctx, f)

	file, err := os.Open("data/output.csv")
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	var matrix [][]float64

	for _, row := range records {
		var floatRow []float64
		for _, value := range row {
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return
			}
			floatRow = append(floatRow, floatValue)
		}
		matrix = append(matrix, floatRow)
	}

	canvas := asciiser.NewImage(len(matrix[0]), len(matrix))
	for i := range matrix {
		for j := range matrix[i] {
			canvas[j][i].Brightness = matrix[j][i]
			canvas[j][i].IsPolygon = true
		}
	}

	cells, err := asciiser.SplitToCells(canvas, 8, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	mtx, err := cells.ConvertToASCII(dctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i := range mtx {
		for j := range mtx[i] {
			fmt.Printf("%c", mtx[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}
