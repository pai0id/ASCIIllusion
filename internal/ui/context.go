package ui

import (
	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

type Context struct {
	mtx asciiser.ASCIImtx
	v   *visualiser.Visualiser
}

func NewContext(mtx asciiser.ASCIImtx, v *visualiser.Visualiser) *Context {
	return &Context{mtx: mtx, v: v}
}
