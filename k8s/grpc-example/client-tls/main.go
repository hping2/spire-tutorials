package main

import (
    "context"
	"time"
    "crypto/tls"
    "crypto/x509"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "log"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	ctx := context.Background()

    // Load the client certificate and its key
    clientCert, err := tls.LoadX509KeyPair("/run/certs/svid.crt", "/run/certs/svid.key")
    if err != nil {
        log.Fatalf("Failed to load client certificate and key. %s.", err)
    }

    // Load the CA certificate
    trustedCert, err := ioutil.ReadFile("/run/certs/root.crt")
    if err != nil {
        log.Fatalf("Failed to load trusted certificate. %s.", err)
    }

    // Put the CA certificate to certificate pool
    certPool := x509.NewCertPool()
    if !certPool.AppendCertsFromPEM(trustedCert) {
        log.Fatalf("Failed to append trusted certificate to certificate pool. %s.", err)
    }

    // Create the TLS configuration
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{clientCert},
        RootCAs:      certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Code copy/pasted and adapted from
			// https://github.com/golang/go/blob/81555cb4f3521b53f9de4ce15f64b77cc9df61b9/src/crypto/tls/handshake_client.go#L327-L344, but adapted to skip the hostname verification.
			// See https://github.com/golang/go/issues/21971#issuecomment-412836078.

			// If this is the first handshake on a connection, process and
			// (optionally) verify the server's certificates.
			certs := make([]*x509.Certificate, len(rawCerts))
			for i, asn1Data := range rawCerts {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					return err
				}
				certs[i] = cert
			}

			opts := x509.VerifyOptions{
				Roots:         certPool,
				CurrentTime:   time.Now(),
				DNSName:       "", // <- skip hostname verification
				Intermediates: x509.NewCertPool(),
			}

			for i, cert := range certs {
				if i == 0 {
					continue
				}
				opts.Intermediates.AddCert(cert)
			}
			_, err := certs[0].Verify(opts)
			return err
		},
    }

    // Create a new TLS credentials based on the TLS configuration
    cred := credentials.NewTLS(tlsConfig)

    // Dial the gRPC server with the given credentials
    conn, err := grpc.Dial("localhost:5443", grpc.WithTransportCredentials(cred))
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        err = conn.Close()
        if err != nil {
            log.Printf("Unable to close gRPC channel. %s.", err)
        }
    }()

    greeterClient := helloworld.NewGreeterClient(conn)
	const interval = time.Second * 10
    log.Printf("Issuing requests every 10 seconds ...")
    for {
        issueRequest(ctx,greeterClient)
        time.Sleep(interval)
    }
}

func issueRequest(ctx context.Context, c helloworld.GreeterClient) {
    resp, err := c.SayHello(ctx, &helloworld.HelloRequest{
            Name: "Joe",
    })

    if err != nil {
        log.Printf("Failed to say hello: %s", err)
        return
    }

    serverID := "SOME-Server"
    log.Printf("%s said %s", serverID, resp.Message)
}
