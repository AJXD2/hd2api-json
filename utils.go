package main

import (
	"fmt"
	"os"
)

func must[T any](val T, err error) T {
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	return val
}
