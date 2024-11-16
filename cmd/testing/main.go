package main

func main() {
	// ctx := mapping.NewContext(11, 11, 4, 4, 44)
	// f, err := fontparser.GetFontMap("fonts/IBM.ttf", 44, 44, 20, 144)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// brightnessMap := mapping.GetBrightnessMap(mapping.FontMapToCellSlice(f))

	// shapeMap := mapping.GetShapeMap(ctx, mapping.FontMapToCellSlice(f))

	// canvas := make([][]bool, 44*20)
	// for i := range canvas {
	// 	canvas[i] = make([]bool, 44*100)
	// }
	// for x := range canvas {
	// 	if (x-440)*(x-440)/100 < len(canvas[0]) {
	// 		canvas[x][(x-440)*(x-440)/100] = true
	// 	}
	// }

	// cells := splitMatrix(canvas, 880, 4400, 44, 44)
	// for i := range cells {
	// 	for j := range cells[i] {
	// 		dv := mapping.GetDescriptionVector(ctx, Cell(cells[i][j]))
	// 		mindelt := 10000
	// 		var minid int
	// 		for id, dvf := range shapeMap {
	// 			d, err := mapping.GetVectorDelt(dv, dvf)
	// 			if err != nil {
	// 				fmt.Printf("Error: %v\n", err)
	// 				return
	// 			}
	// 			if d < mindelt {
	// 				mindelt = d
	// 				minid = id
	// 			}
	// 		}
	// 		c, err := fontparser.GetChar(minid)
	// 		if err != nil {
	// 			fmt.Printf("Error: %v\n", err)
	// 			return
	// 		}
	// 		fmt.Printf("%c", c)
	// 	}
	// 	fmt.Println()
	// }
}
