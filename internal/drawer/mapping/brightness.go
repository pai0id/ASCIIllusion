package mapping

import "github.com/pai0id/CgCourseProject/internal/fontparser"

const MaxBrigtness = 100

func getCellCnt(c Cell) int {
	cnt := 0
	cellData := c.GetData()
	for i := range cellData {
		for _, v := range cellData[i] {
			if v {
				cnt++
			}
		}
	}
	return cnt
}

func GetBrightnessMap(cs map[fontparser.Char]Cell) map[fontparser.Char]int { // [0, 100]
	res := make(map[fontparser.Char]int, len(cs))
	maxCnt := 0
	for ch, c := range cs {
		cnt := getCellCnt(c)
		res[ch] = cnt
		if cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt != 0 {
		for ch := range res {
			res[ch] = (MaxBrigtness * res[ch]) / maxCnt
		}
	}
	return res
}

func GetBrightness(c Cell) int { // [0, 100]
	return (MaxBrigtness * getCellCnt(c)) / (len(c.GetData()) * len(c.GetData()[0]))
}
