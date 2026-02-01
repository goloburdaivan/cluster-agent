package main

import (
	"log"
)

func main() {
	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app.Start()
}
