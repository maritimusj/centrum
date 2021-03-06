package lang

import (
	"fmt"

	"github.com/maritimusj/centrum/synchronized"
)

const (
	_ = iota
	ZhCN
	EnUS
	ZhTW
)

var (
	regionIndex = EnUS
)

var (
	langMap = map[int]map[StrIndex]string{}
)

func Register(region int, lang map[StrIndex]string, err map[ErrIndex]string) {
	langMap[region] = lang
	errStrMap[region] = err
}

func Active(r int) {
	regionIndex = r
}

//func Lang() map[string]int {
//	return map[string]int{
//		"zhCN": ZhCN,
//		"enUS": EnUS,
//	}
//}

func Str(index StrIndex, params ...interface{}) string {
	str := <-synchronized.Do("lang.str", func() interface{} {
		if region, ok := langMap[regionIndex]; ok {
			if str, ok := region[index]; ok {
				return fmt.Sprintf(str, params...)
			}
		}
		return fmt.Sprintf("<unknown string index, region: %d, index: %d>", regionIndex, index)
	})
	return str.(string)
}
