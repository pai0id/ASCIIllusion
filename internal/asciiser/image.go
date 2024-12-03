package asciiser

type Pixel struct {
	Brightness float64
	IsLine     bool
}

type Image [][]Pixel

func NewImage(width, height int) Image {
	res := make([][]Pixel, height)
	for i := range res {
		res[i] = make([]Pixel, width)
	}
	return res
}
