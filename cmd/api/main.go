package main

import "log"

func main() {
	app, cleanup, err := InitializeApp()
	defer cleanup()

	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}
