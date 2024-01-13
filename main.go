package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncarlier/za/pkg/api"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/logger"
	_ "github.com/ncarlier/za/pkg/outputs/all"
	"github.com/ncarlier/za/pkg/server"
	"github.com/ncarlier/za/pkg/version"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: za OPTIONS\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func main() {
	// parse command line
	flag.Parse()

	// show version if asked
	if *version.ShowVersion {
		version.Print()
		os.Exit(0)
	}

	// init configuration file
	if config.InitConfigFlag != nil && *config.InitConfigFlag != "" {
		if err := config.WriteDefaultConfigFile(*config.InitConfigFlag); err != nil {
			log.Fatalf("unable to init configuration file: %v", err)
		}
		os.Exit(0)
	}

	// load configuration file
	conf := config.NewConfig()
	if config.ConfigFileFlag != nil && *config.ConfigFileFlag != "" {
		if err := conf.LoadConfigFromfile(*config.ConfigFileFlag); err != nil {
			log.Fatalf("unable to load configuration file: %v", err)
		}
	}

	logger.Configure(conf.Log.Format, conf.Log.Level)

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
