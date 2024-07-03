package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"time"
)

// generateCert creates a self-signed X.509 certificate using ECDSA with the P256 curve.
// The certificate is valid for 30 days and is intended for server authentication in this enviroment for testing only.
// In a production environment, you should use a signed certificate from a trusted CA.
func generateCert() (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 30),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})

	cert, err := tls.X509KeyPair(certPEM, privateKeyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

// Run starts a TLS server on port 8443
// It generates a self-signed certificate from generateCert()
func Run() {
	cert, err := generateCert()
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	// Create TLS configuration with certificate
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server is running on :8443")

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection manages individual client connections
// It echoes any recieved messages back to the client
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New client connected: %s\n", conn.RemoteAddr())
	_, err := conn.Write([]byte("Welcome to the TLS server!\n"))
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	buffer := make([]byte, 1024)
	for {
		// Read incoming message
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from client: %v", err)
			}
			break
		}
		fmt.Printf("Received from client: %s", buffer[:n])

		// Echo message back
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Printf("Error writing to client: %v", err)
			break
		}
	}

	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}
