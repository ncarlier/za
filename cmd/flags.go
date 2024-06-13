package cmd

import (
	"flag"
	"fmt"
	"os"
)

var (
	// ConfigFlag is the flag used to load the config file
	ConfigFlag string
)

func init() {
	defaultValue := ""
	if value, ok := os.LookupEnv("ZA_CONFIG"); ok {
		defaultValue = value
	}
	flag.StringVar(&ConfigFlag, "c", defaultValue, "Configuration file to load [ENV: ZA_CONFIG]")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "ZerØ Analytics is a Google Analytics alternative focused on privacy and simplicity.\n")
		fmt.Fprintf(out, "\nUsage:\n  za [flags] [command]\n")
		fmt.Fprintf(out, "\nAvailable Commands:\n")
		for _, c := range Commands {
			c.Usage()
		}
		fmt.Fprintf(out, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(out, "\nUse \"za [command] --help\" for more information about a command.\n\n")
	}
}
