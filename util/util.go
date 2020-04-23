package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"golang.org/x/crypto/bcrypt"
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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FormatFileSize(fileSize uint64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%dB", fileSize)
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func FormatDatetime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

//a的n次方
func Exponent(a, n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}
