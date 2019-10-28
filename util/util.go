package util

import (
	"github.com/maritimusj/centrum/gate/lang"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func If(cond bool, yes interface{}, no interface{}) interface{} {
	var result interface{}
	if cond {
		result = yes
	} else {
		result = no
	}
	switch v := result.(type) {
	case func():
		v()
		return nil
	case func() error:
		return v()
	case func() interface{}:
		return v()
	default:
		return v
	}
}

func HashPassword(password []byte) ([]byte, error) {
	data, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, lang.InternalError(err)
	}
	return data, nil
}

func ComparePassword(hashedPassword []byte, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}

const (
	RandNum   = iota // 纯数字
	RandLower        // 小写字母
	RandUpper        // 大写字母
	RandAll          // 数字、大小写字母
)

// 随机字符串
func RandStr(size int, kind int) string {
	rand.Seed(time.Now().UnixNano())

	kinds, result := [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0

	for i := 0; i < size; i++ {
		if isAll { // random kind
			kind = rand.Intn(3)
		}
		scope, base := kinds[kind][0], kinds[kind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
