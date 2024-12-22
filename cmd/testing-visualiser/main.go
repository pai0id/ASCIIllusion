package main

import (
	"log"

	"github.com/pai0id/CgCourseProject/internal/ui"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

func main() {
	v, err := visualiser.NewVisualiser("./fonts/IBM_config_xfce.json", "./fonts/slice.json", "./fonts/IBM.ttf")
	if err != nil {
		log.Printf("error occured: %v", err)
	}
	tuictx := ui.NewContext(nil, v)
	app := ui.NewApp(tuictx)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
