package main

import (
	"log"

	"github.com/pai0id/CgCourseProject/internal/ui"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

func main() {
	v, err := visualiser.NewVisualiser("./fonts/IBM_config.json", "./fonts/slice.json", "./fonts/IBM.ttf")
	if err != nil {
		log.Printf("error occured: %v", err)
	}
	tuictx := ui.NewContext(nil, v)
	// p := tea.NewProgram(initialModel, tea.WithAltScreen())
	// initialModel.AppendCommand("q", tui.Quit)
	// initialModel.AppendCommand("l", tui.LoadObject)

	// if _, err := p.Run(); err != nil {
	// 	log.Printf("error: %v", err)
	// 	os.Exit(1)
	// }
	app := ui.NewApp(tuictx)

	// Run the UI application
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
