package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vjudge/internal/webhook"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read the configuration
	if len(os.Args) > 1 {
		readConfig(os.Args[1])
	} else {
		readConfig("config/config-webhook.json")
	}
	r := gin.Default()
	webhook.Secret = []byte(config.Secret)
	r.POST(config.Endpoint, webhook.Webhook)
	srv := &http.Server{
		Addr:    config.ListenAddress,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	slog.Info("Shutting down the server...")
	_ = srv.Shutdown(context.Background())
}
