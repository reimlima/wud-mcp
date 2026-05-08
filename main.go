package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
	"github.com/reimlima/wud-mcp/tools"
)

func main() {
	c := client.New()

	s := server.NewMCPServer("wud-mcp", "0.1.0")

	tools.RegisterApp(s, c)
	tools.RegisterContainers(s, c)
	tools.RegisterRegistries(s, c)
	tools.RegisterTriggers(s, c)
	tools.RegisterWatchers(s, c)
	tools.RegisterStore(s, c)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
