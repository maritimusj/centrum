package event

import (
	"context"
	"fmt"
	"github.com/kr/pretty"
	"testing"
)

func TestAll(t *testing.T) {
	ctx := context.Background()
	bus := New()
	ch := bus.Sub(ctx, 1, 2)
	go func() {
		for data := range ch {
			if data != nil {
				fmt.Printf("%# v\r\n", pretty.Formatter(data))
			}
		}
	}()

	bus.Register(ctx, func(ctx context.Context, code int, values map[string]interface{}) {
		fmt.Println("code: ", code)
		fmt.Printf("%# v\r\n", pretty.Formatter(values))
	}, 2, 3)

	for i := 1; i < 4; i++ {
		bus.Fire(&Data{
			Code: i,
			Values: map[string]interface{}{
				"one":   "one",
				"two":   2,
				"three": []int{1, 2, 3},
			},
		})
	}
}
