package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
)

var objId int64 = 1
var lsId int64 = 1

func (a *App) parseEntry(s string) {
	parts := strings.Split(s, " ")
	cmd := parts[0]
	switch cmd {
	case "l":
		if !(len(parts) == 2 || (len(parts) == 3 && parts[2] == "$")) {
			a.output.SetText("Load command: l FILEPATH [$]")
		} else {
			obj, err := reader.LoadOBJ(parts[1])
			if err != nil {
				a.output.SetText(fmt.Sprintf("error loading object file: %v\n", err))
				return
			}
			if len(parts) == 3 {
				obj.Skeletonize = true
			}
			a.canvas.ctx.v.AddObj(obj, objId)

			a.canvas.ctx.v.OptimizeCamera()
			a.reload()
			a.output.SetText(fmt.Sprintf("Load complete: ID = %d", objId))
			objId++
		}
	case "rm":
		if len(parts) != 2 {
			a.output.SetText("Remove command: rm ID")
		} else {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			err = a.canvas.ctx.v.DeleteObj(id)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error deleting object: %v\n", err))
				return
			}
			a.reload()
			a.output.SetText(fmt.Sprintf("Removed object with ID = %d", id))
		}
	case "r":
		if len(parts) != 4 {
			a.output.SetText("Rotate command: r ID ANGLE AXIS(x/y/z)")
		} else {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			angle, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing angle: %v\n", err))
				return
			}
			switch parts[3] {
			case "x":
				a.canvas.ctx.v.RotateObj(id, angle, transformer.XAxis)
			case "y":
				a.canvas.ctx.v.RotateObj(id, angle, transformer.YAxis)
			case "z":
				a.canvas.ctx.v.RotateObj(id, angle, transformer.ZAxis)
			default:
				a.output.SetText(fmt.Sprintf("invalid rotation axis: %s\n", parts[2]))
				return
			}
			a.reload()
			a.output.SetText(fmt.Sprintf("Rotated object with ID = %d", id))
		}
	case "t":
		if len(parts) != 5 {
			a.output.SetText("Translate command: t ID tX tY tZ")
		} else {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			tx, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing tx: %v\n", err))
				return
			}
			ty, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing ty: %v\n", err))
				return
			}
			tz, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing tz: %v\n", err))
				return
			}
			a.canvas.ctx.v.TranslateObj(id, tx, ty, tz)
			a.reload()
			a.output.SetText(fmt.Sprintf("Translated object with ID = %d", id))
		}
	case "s":
		if len(parts) != 5 {
			a.output.SetText("Scale command: s ID sX sY sZ")
		} else {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing object id: %v\n", err))
				return
			}
			sx, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing sx: %v\n", err))
				return
			}
			sy, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing sy: %v\n", err))
				return
			}
			sz, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing sz: %v\n", err))
				return
			}
			a.canvas.ctx.v.ScaleObj(id, sx, sy, sz)
			a.reload()
			a.output.SetText(fmt.Sprintf("Scaled object with ID = %d", id))
		}
	case "ls":
		if len(parts) != 4 {
			a.output.SetText("Add light source command: ls X Y Z")
		} else {
			x, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing x: %v\n", err))
				return
			}
			y, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing y: %v\n", err))
				return
			}
			z, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing z: %v\n", err))
				return
			}
			a.canvas.ctx.v.AddLightSource(x, y, z, lsId)
			lsId++
			a.reload()
		}
	case "rmls":
		if len(parts) != 2 {
			a.output.SetText("Remove light source command: rmls ID")
		} else {
			id, err := strconv.ParseInt(parts[1], 10, 0)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error parsing light source id: %v\n", err))
				return
			}
			err = a.canvas.ctx.v.DeleteLightSource(id)
			if err != nil {
				a.output.SetText(fmt.Sprintf("error deleting light source: %v\n", err))
				return
			}
			a.reload()
		}
	case "h":
		helpMsg := "Help:\n"
		helpMsg += "Load command: l FILEPATH\n"
		helpMsg += "Remove command: rm ID\n"
		helpMsg += "Rotate command: r ID ANGLE AXIS(x/y/z)\n"
		helpMsg += "Translate command: t ID tX tY tZ\n"
		helpMsg += "Scale command: s ID sX sY sZ\n"
		helpMsg += "Add light source command: ls X Y Z\n"
		helpMsg += "Remove light source command: rmls ID\n"
		helpMsg += "Quit command: q\n"
		a.canvas.SetText(helpMsg)
	case "q":
	default:
		a.reload()
	}
}

func (a *App) reload() {
	mtx, err := a.canvas.ctx.v.Reconvert()
	if err != nil {
		a.output.SetText(fmt.Sprintf("error converting matrix: %v\n", err))
		return
	}
	a.canvas.ctx.mtx = mtx

	a.canvas.RenderMatrix()
}
