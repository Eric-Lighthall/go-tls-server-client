package main

import (
	"fmt"
	"tlsgo/pkg/client"
)

// Entry point for the client application. This starts the client connection to the server.
func main() {
	fmt.Println("Starting client...")
	client.Run()
}
