package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MichaelGenchev/smart-home-system/internal/app/gateway"

	"github.com/MichaelGenchev/smart-home-system/internal/common/config"
)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ",err)
	}

	app, err := gatewayapp.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to create app: ",err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal("Failed to run app: ",err)
		}
	}()

	<-ctx.Done()
	app.Shutdown()
}