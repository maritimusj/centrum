package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	pid := fmt.Sprintf("%d", os.Getpid())
	err := ioutil.WriteFile("./edge.pid", []byte(pid), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	p, err := os.FindProcess(4396)
	if err != nil {
		log.Fatal(err)
	}

	st, err := p.Wait()
	if err != nil {
		log.Fatal(err)
	}

	println(st.Exited())
}
