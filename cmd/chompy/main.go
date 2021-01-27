package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/hebo/chompy/config"
	"github.com/hebo/chompy/server"
	"github.com/spf13/afero"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Failed to parse config", err)
	}

	server := server.New(cfg, afero.NewOsFs())
	server.Serve(cfg.Port)
}
