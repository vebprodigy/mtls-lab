package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	// Load CA bundle (root + intermediate) to verify client certs
	clientCABytes, err := ioutil.ReadFile("ca-bundle.crt")
	if err != nil {
		log.Fatal("Failed to load CA bundle:", err)
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(clientCABytes) {
		log.Fatal("Failed to parse client CA cert")
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		// Load server certificate and key
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{}, // Will be loaded below
			},
		},
		// Require and verify client certificate
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  clientCAs,
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	// Load server certificate
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal("Failed to load server cert:", err)
	}
	tlsConfig.Certificates[0] = cert

	// Create listener
	listener, err := tls.Listen("tcp", "localhost:4443", tlsConfig)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	defer listener.Close()

	fmt.Println("[SERVER] mTLS server listening on localhost:4443")

	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	tlsConn := conn.(*tls.Conn)

	// Get client certificate
	clientCerts := tlsConn.ConnectionState().PeerCertificates
	if len(clientCerts) > 0 {
		clientCert := clientCerts[0]
		fmt.Printf("[SERVER] Client connected: CN=%s\n", clientCert.Subject.CommonName)
	}

	// Read request
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading:", err)
		return
	}

	message := string(buf[:n])
	fmt.Printf("[SERVER] Received: %s\n", message)

	// Send response
	response := fmt.Sprintf("Echo: %s", message)
	conn.Write([]byte(response))
	fmt.Printf("[SERVER] Sent: %s\n", response)
}
