package main

import (
	"fmt"

	"github.com/drekunov/gc/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
