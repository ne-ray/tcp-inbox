package main

import (
	"log"

	"github.com/ne-ray/tcp-inbox/config"
	"github.com/ne-ray/tcp-inbox/internal/app"
)

const CONFIG_PATH = "./config/config.yaml"

func main() {
	// Configuration
	cfg, err := config.NewConfig(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
