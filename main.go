package main

import (
	"flag"
	"log"
	"os"

	_ "solana-telegraf/plugins/inputs/solana"

	"github.com/influxdata/telegraf/plugins/common/shim"
)

func main() {
	var pollInterval = flag.Duration("poll_interval", 0, "how often to send metrics")
	var configFile = flag.String("config", "", "path to the config file for this plugin")

	flag.Parse()

	shim := shim.New()

	if err := shim.LoadConfig(configFile); err != nil {
		log.Printf("ERROR: loading config failed: %s\n", err)
		os.Exit(1)
	}

	if err := shim.Run(*pollInterval); err != nil {
		log.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
