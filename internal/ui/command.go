package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

func (a *App) parseEntry(s string) {
	parts := strings.Split(s, " ")
	cmd := parts[0]
	switch cmd {
	case "l":
		if len(parts) == 2 {
			a.loadObject(parts[1])
		}
	case "r":
		if len(parts) == 4 {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			angle, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing angle: %v\n", err))
				return
			}
			switch parts[3] {
			case "x":
				a.canvas.ctx.v.RotateObj(int(id), angle, transformer.XAxis)
			case "y":
				a.canvas.ctx.v.RotateObj(int(id), angle, transformer.YAxis)
			case "z":
				a.canvas.ctx.v.RotateObj(int(id), angle, transformer.ZAxis)
			default:
				a.canvas.SetText(fmt.Sprintf("invalid rotation axis: %s\n", parts[2]))
				return
			}
			a.reload()
		}
	case "t":
		if len(parts) == 5 {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			tx, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing tx: %v\n", err))
				return
			}
			ty, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing ty: %v\n", err))
				return
			}
			tz, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing tz: %v\n", err))
				return
			}
			a.canvas.ctx.v.TranslateObj(int(id), tx, ty, tz)
			a.reload()
		}
	case "s":
		if len(parts) == 5 {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			sx, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing sx: %v\n", err))
				return
			}
			sy, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing sy: %v\n", err))
				return
			}
			sz, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing sz: %v\n", err))
				return
			}
			a.canvas.ctx.v.ScaleObj(int(id), sx, sy, sz)
			a.reload()
		}
	case "ls":
		if len(parts) == 4 {
			x, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing x: %v\n", err))
				return
			}
			y, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing y: %v\n", err))
				return
			}
			z, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.canvas.SetText(fmt.Sprintf("error parsing z: %v\n", err))
				return
			}
			a.canvas.ctx.v.AddLightSource(x, y, z)
			a.reload()
		}
	default:
		a.reload()
	}
}

func (a *App) loadObject(filepath string) {
	obj, err := reader.LoadOBJ(filepath)
	if err != nil {
		a.canvas.SetText(fmt.Sprintf("error loading object file: %v\n", err))
		return
	}
	a.canvas.ctx.v.AddObj(obj)

	a.canvas.ctx.v.OptimizeCamera()
	a.reload()
}

func (a *App) reload() {
	mtx, err := a.canvas.ctx.v.Reconvert()
	if err != nil {
		a.canvas.SetText(fmt.Sprintf("error converting matrix: %v\n", err))
		return
	}
	a.canvas.ctx.mtx = mtx

	a.canvas.RenderMatrix()
}
