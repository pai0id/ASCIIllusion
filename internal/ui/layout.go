package ui

import (
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
)

type App struct {
	ui     tui.UI
	canvas *Canvas
	output *tui.Label
}

func NewApp(tuictx *Context) *App {
	canvas := NewCanvas(tuictx)

	output := tui.NewLabel("Hello!")

	a := &App{canvas: canvas, output: output}

	entryField := NewEntryField(a.parseEntry)

	root := tui.NewVBox(
		canvas.Label,
		tui.NewSpacer(),
		entryField.Entry,
		output,
	)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+c", func() { ui.Quit() })
	a.ui = ui

	go func() {
		for {
			s := root.Size()
			a.canvas.ctx.v.Resize(s.X, s.Y)
			time.Sleep(time.Millisecond)
		}
	}()

	return a
}

func (a *App) Run() error {
	time.Sleep(time.Millisecond)
	return a.ui.Run()
}
