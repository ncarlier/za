package config

import (
	"flag"
	"os"
)

// ConfigFile is the flag used to load the config file
var ConfigFileFlag *string

// InitConfigFlag is the flag used to initialize the config file
var InitConfigFlag = flag.String("init-config", "", "Initialize configuration file")

func init() {
	defaultValue := ""
	if value, ok := os.LookupEnv("ZA_CONFIG"); ok {
		defaultValue = value
	}
	ConfigFileFlag = flag.String("config", defaultValue, "Configuration file to load [ENV: ZA_CONFIG]")
}
