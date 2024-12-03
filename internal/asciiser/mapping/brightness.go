package mapping

import "github.com/pai0id/CgCourseProject/internal/fontparser"

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

func GetBrightnessMap(cs map[fontparser.Char]Cell) map[fontparser.Char]float64 { // [0, 100]
	res := make(map[fontparser.Char]float64, len(cs))
	maxCnt := 0
	for ch, c := range cs {
		cnt := getCellCnt(c)
		res[ch] = float64(cnt)
		if cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt != 0 {
		for ch := range res {
			res[ch] = res[ch] / float64(maxCnt)
		}
	}
	return res
}

func GetBrightness(c Cell) int { // [0, 100]
	return getCellCnt(c) / (len(c.GetData()) * len(c.GetData()[0]))
}
