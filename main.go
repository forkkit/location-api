package main

import (
	log "github.com/golang/glog"
	"github.com/myodc/geo-api/handler"
	"github.com/myodc/go-micro/cmd"
	"github.com/myodc/go-micro/server"
)

func main() {
	// optionally setup command line usage
	cmd.Init()

	server.Name = "go.micro.api.geo"

	// Initialise Server
	server.Init()

	// Register Handlers
	server.Register(
		server.NewReceiver(new(handler.Location)),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
