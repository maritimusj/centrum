package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/maritimusj/centrum/util"
)

func main() {
	owner := flag.String("o", "", "owner name")
	flag.Parse()

	*owner = strings.TrimSpace(*owner)
	if *owner == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	code := strings.ToLower(util.RandStr(4, util.RandAll))
	hash := hmac.New(sha1.New, []byte(code))
	v := fmt.Sprintf("%x", hash.Sum([]byte(*owner)))
	fmt.Printf("register code for '%s' is: %s-%s-%s\r\n", *owner, code, v[:4], v[len(v)-4:])
}
