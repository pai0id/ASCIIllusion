package mapping

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

func GetBrightnessMap(cs []Cell) []int { // [0, 100]
	res := make([]int, 0, len(cs))
	maxCnt := 0
	for _, c := range cs {
		cnt := getCellCnt(c)
		res = append(res, cnt)
		if cnt > maxCnt {
			maxCnt = cnt
		}
	}

	if maxCnt != 0 {
		for i := range res {
			res[i] = (MaxBrigtness * res[i]) / maxCnt
		}
	}
	return res
}

func GetBrightness(c Cell) int { // [0, 100]
	return (MaxBrigtness * getCellCnt(c)) / (len(c.GetData()) * len(c.GetData()[0]))
}
