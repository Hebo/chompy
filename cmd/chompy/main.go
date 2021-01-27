package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/hebo/chompy/config"
	"github.com/hebo/chompy/server"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	server := server.New(cfg)
	server.Serve(cfg.Port)
}
