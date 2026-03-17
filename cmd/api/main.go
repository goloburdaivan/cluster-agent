package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app.Start()
}
