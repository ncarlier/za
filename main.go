package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ncarlier/za/cmd"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/logger"

	_ "github.com/ncarlier/za/cmd/all"
	_ "github.com/ncarlier/za/pkg/outputs/all"
)

func main() {
	// parse command line
	flag.Parse()

	// load configuration
	conf := config.NewConfig()
	if cmd.ConfigFlag != "" {
		if err := conf.LoadFile(cmd.ConfigFlag); err != nil {
			log.Fatalf("unable to load configuration file: %v", err)
		}
	}

	// configure the logger
	logger.Configure(conf.Log.Format, conf.Log.Level)

	args := flag.Args()
	commandLabel, idx := cmd.GetFirstCommand(args)

	if command, ok := cmd.Commands[commandLabel]; ok {
		if err := command.Exec(args[idx+1:], conf); err != nil {
			log.Fatalf("error during command execution: %v", err)
		}
	} else {
		if commandLabel != "" {
			fmt.Fprintf(os.Stderr, "undefined command: %s\n", commandLabel)
		}
		flag.Usage()
		os.Exit(0)
	}
}
