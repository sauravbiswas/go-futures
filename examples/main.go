package main

import (
	"fmt"
	"github.com/sauravbiswasiupr/go-futures/futures"
	"time"
)

func main() {
	f := futures.NewFuture(func() (int, error) {
		time.Sleep(1 * time.Second)
		return 42, nil
	})

	f.OnSuccess(func(val int) {
		fmt.Println("Success:", val)
	})

	f.OnFailure(func(err error) {
		fmt.Println("Failed:", err)
	})

	f.Start()

	// Chain another task
	f2 := f.Then(func(i int) (any, error) {
		return 42, nil
	})

	result, err := f2.Result()
	if err != nil {
		fmt.Println("Final error:", err)
	} else {
		fmt.Println("Final result:", result)
	}
}
