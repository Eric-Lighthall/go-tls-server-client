package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
)

// Run initiates a TLS connection to the server and handles user input
func Run() {
	config := &tls.Config{InsecureSkipVerify: true} // This should be avoided in prod. We're skipping the verification of the server's certificate here. This could expose users to MITM attacks.
	conn, err := tls.Dial("tcp", "localhost:8443", config)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// Continuously  read messages from the server
	go readFromServer(conn)

	// Loop for sending user input to the server
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message (or 'quit' to exit): ")
		if !scanner.Scan() {
			break
		}
		message := scanner.Text()
		if message == "quit" {
			break
		}
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}

	fmt.Println("Disconnected from server")
}

// readFromServer continuously reads and displays messages from the server
func readFromServer(conn *tls.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from server: %v", err)
			return
		}
		fmt.Printf("Received from server: %s", buffer[:n])
	}
}
