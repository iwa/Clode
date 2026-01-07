package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwa/Clode/internal/app"
	"github.com/iwa/Clode/internal/config"
)

func main() {
	config := config.GenerateConfigFromEnv()

	app, err := app.NewApp(config)
	if err != nil {
		panic(err)
	}

	log.Println("Starting Discord Bot...")
	if err := app.DiscordClient.Start(); err != nil {
		log.Fatalf("Failed to start Discord client: %v", err)
	}

	defer app.DiscordClient.Stop()

	// Waiting for sigint or sigterm signal to shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
