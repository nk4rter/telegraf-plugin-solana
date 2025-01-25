package main

import (
	"flag"
	"log"
	"os"
	"time"

	_ "telegraf-plugin-solana/plugins/inputs/solana"

	"github.com/influxdata/telegraf/plugins/common/shim"
)

func main() {
	var pollInterval = flag.Duration("poll_interval", 1*time.Second, "how often to send metrics")
	var configFile = flag.String("config", "", "path to the config file for this plugin")
	var err error

	flag.Parse()

	shim := shim.New()

	err = shim.LoadConfig(configFile)
	if err != nil {
		log.Printf("ERROR: loading config failed: %s\n", err)
		os.Exit(1)
	}

	if err := shim.Run(*pollInterval); err != nil {
		log.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}
