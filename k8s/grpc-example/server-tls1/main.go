package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract client certificate from the request
		cert := r.TLS.PeerCertificates[0]

		// Extract information from the client certificate
		//commonName := cert.Subject.CommonName
		//organization := cert.Subject.Organization
		spiffeID := cert.URIs[0].String()

		// Use the extracted information
		fmt.Println("Connection from Client with spiffeID:", spiffeID)
		res := "From Server: how are you?"
		fmt.Println("Send response: :", res)
		//fmt.Println("Client Certificate Common Name (CN):", commonName)
		//fmt.Println("Client Certificate Organization (O):", organization)
		fmt.Fprint(w, res)
	})

	// Load the server certificate and key
	serverCert, err := tls.LoadX509KeyPair("/run/certs/svid.crt", "/run/certs/svid.key")
	if err != nil {
		fmt.Println("Error loading server certificate and key:", err)
		return
	}

	// Load the CA certificate
	caCert, err := ioutil.ReadFile("/run/certs/root.crt")
	if err != nil {
		fmt.Println("Error loading CA certificate:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure the TLS settings for the server
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}

	// Create the HTTPS server
	server := &http.Server{
		Addr:      ":8111",
		Handler:   nil,
		TLSConfig: tlsConfig,
	}

	// Start the server
	fmt.Println("Starting server at port 8111")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

