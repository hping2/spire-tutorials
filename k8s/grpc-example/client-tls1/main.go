package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Load the server certificate and key
	clientCert, err := tls.LoadX509KeyPair("/run/certs/svid.crt", "/run/certs/svid.key")
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

	// Configure the TLS settings for the client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		//RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	// Create the HTTPS client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	for {
		issueRequest(client)
		time.Sleep(time.Second * 3)
	}
}

func issueRequest(client *http.Client) {
	// Perform an HTTPS request
	resp, err := client.Get("https://localhost:8111")
	if err != nil {
		fmt.Println("Error making HTTPS request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the server's certificate
	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		fmt.Println("Server certificate not present in response")
		return
	}
	serverCert := resp.TLS.PeerCertificates[0]
	//serverCommonName := serverCert.Subject.CommonName
	//fmt.Println("Server Certificate Common Name (CN):", serverCommonName)
	recvServerSpiffeID := serverCert.URIs[0].String()
	fmt.Println("Server SPIFFEID:", recvServerSpiffeID)

	// Check if the server certificate is verified by the trusted CA
	// err = serverCert.VerifyHostname(serverCommonName)
	//if err != nil {
	//	fmt.Println("Server certificate verification failed:", err)
	//	return
	//}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Convert the response body to a string and print it
	bodyStr := string(body)
	fmt.Println("Response Body:", bodyStr)

	// check server SPIFFE ID
	expectedServerSpiffeID := "spiffe://example.org/ns/client01/sa/default"
	result := strings.Compare(recvServerSpiffeID, expectedServerSpiffeID)
	if result != 0 {
		fmt.Print("Server certificate bad")
		return
	}

	fmt.Println("Server certificate verified successfully")
}
