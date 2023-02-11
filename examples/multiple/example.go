package main

import (
	"errors"
	"fmt"
	"github.com/garugaru/resilience4go/circuitbreaker"
	"github.com/garugaru/resilience4go/fallback"
	"github.com/garugaru/resilience4go/resilience"
	"github.com/garugaru/resilience4go/retry"
	"time"
)

func main() {
	var executor = resilience.New[string](
		retry.New[string](retry.NewFixedDelay(10*time.Millisecond), 100),
		circuitbreaker.New[string](circuitbreaker.NewCountBasedWindow(10), 0.5, 0.5, 100*time.Millisecond),
		fallback.New[string]("fallback"),
	)

	for i := 0; i < 100; i++ {
		out, err := executor.Execute(func() (string, error) {
			return "", errors.New("error")
		})

		fmt.Println(out)
		fmt.Println(err)
	}

	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 100; i++ {
		out, err := executor.Execute(func() (string, error) {
			return "OK", nil
		})
		time.Sleep(10 * time.Millisecond)

		fmt.Println(out)
		fmt.Println(err)
	}

}
