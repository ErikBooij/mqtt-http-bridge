package main

import (
	"context"
	"log"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/process"
	"os"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()

	if err != nil {
		log.Printf("Unable to load config. Stopping.\n\n%s\n", err)
		os.Exit(1)
		return
	}

	appStartErr := make(chan error)
	done := make(chan bool)

	go func() {
		process.Start(ctx, cfg, appStartErr)
		done <- true
	}()

	err = <-appStartErr

	if err != nil {
		log.Printf("Error starting process: %s\n", err)
		os.Exit(1)
		return
	}

	log.Printf("Process started successfully.\n")

	<-done
}
