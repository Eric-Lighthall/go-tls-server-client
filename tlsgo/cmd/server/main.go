package main

import (
	"fmt"
	"tlsgo/pkg/server"
)

// Entry point for the server application. This starts the TLS server.
func main() {
	fmt.Println("Starting server...")
	server.Run()
}
