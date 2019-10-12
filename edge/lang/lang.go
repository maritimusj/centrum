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
	langMap = map[int]map[StrIndex]string{}
)

func Register(region int, lang map[StrIndex]string, err map[ErrorCode]string) {
	langMap[region] = lang
	errStrMap[region] = err
}

func Str(index StrIndex, params ...interface{}) string {
	if region, ok := langMap[regionIndex]; ok {
		if str, ok := region[index]; ok {
			return fmt.Sprintf(str, params...)
		}
	}
	return fmt.Sprintf("<unknown string index, region: %d, index: %d>", regionIndex, index)
}
