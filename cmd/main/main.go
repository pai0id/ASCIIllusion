package main

import (
	"flag"
	"log"

	"github.com/pai0id/CgCourseProject/internal/ui"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

func main() {
	cfgFileName := flag.String("font-config", "", "Font configuration file")
	sliceFileName := flag.String("slice", "", "ASCII slice configuration file")
	fontFileName := flag.String("font-file", "", "Font file")

	flag.Parse()

	if *cfgFileName == "" || *sliceFileName == "" || *fontFileName == "" {
		log.Fatal("Please provide all required flags: -font-config, -slice, -font-file")
	}

	v, err := visualiser.NewVisualiser(*cfgFileName, *sliceFileName, *fontFileName)
	if err != nil {
		log.Printf("error occured: %v", err)
	}
	tuictx := ui.NewContext(nil, v)
	app := ui.NewApp(tuictx)

	// Run the UI application
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
