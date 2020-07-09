package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/trackr/pkg/api"
	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/server"
	"github.com/ncarlier/trackr/pkg/version"
	configflag "github.com/ncarlier/webhookd/pkg/config/flag"
)

func main() {
	conf := &config.Config{}
	configflag.Bind(conf, "TRACKR")

	flag.Parse()

	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	level := "info"
	if conf.Debug {
		level = "debug"
	}
	logger.Init(level)

	logger.Debug.Println("starting trackr server...")

	srv := server.NewServer(conf)

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Debug.Println("server is shutting down...")
		api.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error.Fatalf("could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	addr := conf.ListenAddr
	logger.Info.Println("server is ready to handle requests at", addr)
	api.Start()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error.Fatalf("could not listen on %s : %v\n", addr, err)
	}

	<-done
	logger.Debug.Println("server stopped")
}
