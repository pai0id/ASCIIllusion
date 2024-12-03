package ui

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

type Canvas struct {
	*tui.Label
	ctx *Context
}

func NewCanvas(ctx *Context) *Canvas {
	label := tui.NewLabel("")
	label.SetSizePolicy(tui.Maximum, tui.Maximum)
	return &Canvas{
		Label: label,
		ctx:   ctx,
	}
}

func (c *Canvas) RenderMatrix() {
	res := ""
	for i := range c.ctx.mtx {
		for _, ch := range c.ctx.mtx[i] {
			res += fmt.Sprintf("%c", ch)
		}
		res += "\n"
	}

	c.SetText(res)
}
