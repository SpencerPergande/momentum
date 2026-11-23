// Command server serves the Momentum web app.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SpencerPergande/momentum/pkg/config"
	"github.com/SpencerPergande/momentum/pkg/server"
	"github.com/SpencerPergande/momentum/services/web/internal/static"
)

func main() {
	log := server.NewLogger("web")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := server.Run(ctx, server.Options{
		Name:    "web",
		Addr:    config.String("MOMENTUM_ADDR", ":8080"),
		Handler: static.Handler(static.Dist()),
		Logger:  log,
	})
	if err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
