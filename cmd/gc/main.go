package main

import (
	"log"

	"github.com/drekunov/gc/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
