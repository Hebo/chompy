package main

import (
	"flag"

	"github.com/hebo/chompy/server"
)

const (
	defaultDownloadsDir = "./downloads"
	defaultPort         = 8000
)

func main() {
	downloadsDir := flag.String("downloads-dir", defaultDownloadsDir, "Directory for video downloads")
	port := flag.Int("port", defaultPort, "Port to listen on")
	flag.Parse()

	server := server.New(*downloadsDir)
	server.Serve(*port)
}
