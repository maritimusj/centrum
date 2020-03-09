package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maritimusj/centrum/register"
)

func main() {
	owner := flag.String("o", "", "owner name")
	fingerprints := flag.String("p", "", "hardware fingerprints")

	flag.Parse()

	if *owner == "" || *fingerprints == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("owner: %s, code: %s\r\n", *owner, register.Code(*owner, *fingerprints))

}
