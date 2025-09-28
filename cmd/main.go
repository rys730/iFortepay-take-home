package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/rys730/iFortepay-take-home/cmd/app"
	"github.com/rys730/iFortepay-take-home/internal/common/config"
)

func main() {
	cfg := config.NewConfig()
	app := app.CreateApp(cfg)
	go func() {
		if err := app.Start(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
			log.Fatal().Err(err).Msg("closing server")
		}
	}()
	log.Info().Int("port", cfg.App.Port).Msg("app started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Give outstanding requests a chance to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("error shutting down server")
	}
}
