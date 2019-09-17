package lang

import "fmt"

const (
	_ = iota
	ZhCN
)

var (
	regionIndex = ZhCN
)

var (
	langMap = map[int]map[int]string{}
)

func Register(region int, lang map[int]string, err map[ErrorCode]string) {
	langMap[region] = lang
	errStrMap[region] = err
}

func Active(r int) {
	regionIndex = r
}

func Lang() map[string]int {
	return map[string]int{
		"zhCN": ZhCN,
	}
}

func Str(index int, params ...interface{}) string {
	if region, ok := langMap[regionIndex]; ok {
		if str, ok := region[index]; ok {
			return fmt.Sprintf(str, params...)
		}
	}
	return fmt.Sprintf("<unknown string index, region: %d, index: %d>", regionIndex, index)
}
