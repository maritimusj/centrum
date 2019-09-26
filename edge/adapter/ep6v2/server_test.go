package ep6v2

import (
	"context"
	"log"
	"testing"
)

func TestAll(t *testing.T) {

	server := NewInverseServer("", 10502)
	err := server.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	server.Wait()
}
