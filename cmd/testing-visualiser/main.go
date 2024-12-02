package main

import (
	"log"

	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

func main() {
	err := visualiser.NewVisualiser("./fonts/IBM_config.json", "./fonts/slice.json", "./fonts/IBM.ttf")
	if err != nil {
		log.Println("Error occured: %v", err)
	}
	// obj, err := reader.LoadOBJ("data/tetra.obj")
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }
	// visualiser.AddObj(obj)
	// visualiser.
}
