package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Load CA bundle (root + intermediate) to verify server cert
	serverCABytes, err := ioutil.ReadFile("ca-bundle.crt")
	if err != nil {
		log.Fatal("Failed to load CA bundle:", err)
	}

	serverCAs := x509.NewCertPool()
	if !serverCAs.AppendCertsFromPEM(serverCABytes) {
		log.Fatal("Failed to parse server CA cert")
	}

	// Load client certificate
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal("Failed to load client cert:", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      serverCAs,
		MinVersion:   tls.VersionTLS12,
		ServerName:   "server.local", // Match certificate CN/SAN
	}

	// Connect to server
	conn, err := tls.Dial("tcp", "localhost:4443", tlsConfig)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// Get server certificate info
	serverCerts := conn.ConnectionState().PeerCertificates
	if len(serverCerts) > 0 {
		serverCert := serverCerts[0]
		fmt.Printf("[CLIENT] Connected to server: CN=%s\n", serverCert.Subject.CommonName)
	}

	// Send message
	message := "Hello mTLS!"
	fmt.Printf("[CLIENT] Sending: %s\n", message)
	conn.Write([]byte(message))

	// Read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Failed to read:", err)
	}

	response := string(buf[:n])
	fmt.Printf("[CLIENT] Received: %s\n", response)
}
