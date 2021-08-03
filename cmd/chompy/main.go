package main

import (
	"log"

	"net/http"
	_ "net/http/pprof"

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

	go func() {
		// Pprof
		log.Println("Loading pprof on :6060")
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	server := server.New(cfg, afero.NewOsFs())
	server.Serve(cfg.Port)
}
