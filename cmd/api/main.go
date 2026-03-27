package main

import (
	"api/internal/config"
	"api/internal/server"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic("Cannot load config: " + err.Error())
	}

	// 2. Setup and Run Server
	srv := server.New(cfg)
	srv.Run()
}
