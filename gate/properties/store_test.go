package properties

import (
	"fmt"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestStoreImp(t *testing.T) {
	err := Open("test111.db")
	if err != nil {
		t.Fatal(err)
	}

	defer Close()

	err = Write(time.Now().Format("2006-01-02 15:04:05"), "clientId.js.001")
	if err != nil {
		t.Fatal(err)
	}

	err = Write(time.Now().Format("2006-01-02 15:04:05"), "clientId.js.002")
	if err != nil {
		t.Fatal(err)
	}

	err = Write(time.Now().Format("2006-01-02 15:04:05"), "clientId.js.003")
	if err != nil {
		t.Fatal(err)
	}

	data := LoadAllString("clientId.js")
	for name := range data {
		fmt.Println(name)
	}

	fmt.Printf("%#v", pretty.Formatter(data))
}
