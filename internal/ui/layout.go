package ui

import (
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
)

type App struct {
	ui     tui.UI
	canvas *Canvas
}

func NewApp(tuictx *Context) *App {
	canvas := NewCanvas(tuictx)
	a := &App{canvas: canvas}

	entryField := NewEntryField(a.parseEntry)

	root := tui.NewVBox(
		canvas.Label,
		tui.NewSpacer(),
		entryField.Entry,
	)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("q", func() { ui.Quit() })
	a.ui = ui

	go func() {
		for {
			time.Sleep(time.Second)
			s := root.Size()
			a.canvas.ctx.v.Resize(s.X, s.Y)
		}
	}()

	return a
}

func (a *App) Run() error {
	return a.ui.Run()
}
