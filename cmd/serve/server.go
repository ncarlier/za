package serve

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/za/pkg/api"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/server"
)

func startServer(conf *config.Config) {
	slog.Debug("starting ZerØ Analytics server...")

	srv := server.NewServer(conf)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		slog.Debug("server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("could not gracefully shutdown the server", "error", err)
			os.Exit(1)
		}
		close(done)
	}()

	api.Start()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("unable to start the server", "addr", conf.HTTP.ListenAddr, "err", err)
		os.Exit(1)
	}

	<-done
	slog.Debug("server stopped")
}
