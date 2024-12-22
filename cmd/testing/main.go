package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pai0id/CgCourseProject/internal/asciiser"
	"github.com/pai0id/CgCourseProject/internal/object"
	"github.com/pai0id/CgCourseProject/internal/reader"
	"github.com/pai0id/CgCourseProject/internal/transformer"
	"github.com/pai0id/CgCourseProject/internal/visualiser"
)

const (
	dir = "test-data"
	w   = 200
	h   = 200
)

var filenames = []string{
	"./data/sphere.obj",
	"./data/cube.obj",
}

var ls = []object.Vec3{
	{X: 0, Y: 0, Z: 0},
	{X: 10, Y: 10, Z: -10},
	{X: -10, Y: -10, Z: 10},
	{X: 0, Y: 0, Z: 10},
	{X: 0, Y: 0, Z: -10},
}

func WriteASCIIMatrixToFile(filename string, matrix asciiser.ASCIImtx) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, row := range matrix {
		for _, char := range row {
			_, err := fmt.Fprintf(file, "%c", char)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(file, "%c", '\n')
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	v, err := visualiser.NewVisualiser("./fonts/IBM_config_xfce.json", "./fonts/slice.json", "./fonts/IBM.ttf")
	if err != nil {
		log.Printf("error occured: %v", err)
	}
	v.Resize(w, h)

	var skeletonize = []bool{true, false}
	for _, skel := range skeletonize {
		for i, filename := range filenames {
			for j, lsrc := range ls {
				obj, err := reader.LoadOBJ(filename)
				if err != nil {
					log.Printf("error loading object file: %v\n", err)
					continue
				}
				if skel {
					obj.Skeletonize()
				}
				id := v.AddObj(obj)
				v.OptimizeCamera()
				v.RotateObj(id, 45, transformer.YAxis)
				v.RotateObj(id, 45, transformer.XAxis)
				lid := v.AddLightSource(lsrc.X, lsrc.Y, lsrc.Z, 1)

				mtx, err := v.Reconvert()
				if err != nil {
					log.Printf("error converting matrix: %v\n", err)
					continue
				}

				suffix := ""
				if skel {
					suffix = "-skeletonized"
				}
				filename := fmt.Sprintf("%s/%d-%d%s", dir, i, j, suffix)
				WriteASCIIMatrixToFile(filename, mtx)

				v.DeleteObj(id)
				v.DeleteLightSource(lid)
			}
		}
	}
}
