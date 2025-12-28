package main

import "log"

func main() {
	app, err := InitializeApp()

	if err != nil {
		log.Fatal(err)
	}

	app.Start()
}
