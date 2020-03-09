package register

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	fingerprints := Fingerprints()
	fmt.Println("fingerprints:", fingerprints)

	owner := "js"
	code := Code(owner, fingerprints)
	fmt.Println("code:", code)

	if Verify("js", code) {
		fmt.Println("verified")
	}
}
