package main

import (
	"fmt"

	"github.com/pai0id/CgCourseProject/internal/fontparser"
	"github.com/pai0id/CgCourseProject/internal/shape"
)

// Pacifico.ttf
// Ubuntu.ttf
// DroidSansMono.ttf
func main() {
	charMat, err := fontparser.GetFontMap("fonts/DroidSansMono.ttf", 16, 25, 20, 72)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	minD := 0
	var cc1, cc2 int
	for i, c1 := range charMat {
		for j, c2 := range charMat {
			if i == j {
				continue
			}
			dv1 := shape.GetDescriptionVector(shape.NewContext(16, 25, 8, 12, 12.0), shape.Cell(c1))
			dv2 := shape.GetDescriptionVector(shape.NewContext(16, 25, 8, 12, 12.0), shape.Cell(c2))
			d, err := shape.GetVectorDelt(dv1, dv2)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if d < minD {
				minD = d
				cc1 = i
				cc2 = j
			}
		}
	}

	for i := range charMat[cc1] {
		for _, v2 := range charMat[cc1][i] {
			if v2 {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}

	for i := range charMat[cc2] {
		for _, v2 := range charMat[cc2][i] {
			if v2 {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}

	fmt.Printf("Vector delt: %d\n", minD)
}
