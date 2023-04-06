package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "log"
    "net"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)


func main() {
    // Load the server certificate and its key
    serverCert, err := tls.LoadX509KeyPair("svid.crt", "svid.key")
    if err != nil {
        log.Fatalf("Failed to load server certificate and key. %s.", err)
    }

    // Load the CA certificate
    trustedCert, err := ioutil.ReadFile("root.crt")
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
        Certificates: []tls.Certificate{serverCert},
        RootCAs:      certPool,
        ClientCAs:    certPool,
        MinVersion:   tls.VersionTLS13,
        MaxVersion:   tls.VersionTLS13,
    }

    // Create a new TLS credentials based on the TLS configuration
    cred := credentials.NewTLS(tlsConfig)

    // Create a listener that listens to localhost port 8443
    listener, err := net.Listen("tcp", "localhost:8443")
    if err != nil {
        log.Fatalf("Failed to start listener. %s.", err)
    }
    defer func() {
        err = listener.Close()
        if err != nil {
            log.Printf("Failed to close listener. %s\n", err)
        }
    }()

    // Create a new gRPC server
    server := grpc.NewServer(grpc.Creds(cred))
    helloworld.RegisterGreeterServer(server, greeter{}) // Register the demo service

    // Start the gRPC server
    err = server.Serve(listener)
    if err != nil {
        log.Fatalf("Failed to start gRPC server. %s.", err)
    }
}

type greeter struct {
    helloworld.UnimplementedGreeterServer
}

func (greeter) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
    //////////////
    //
    /////////////
    clientID := "Some-client"

    log.Printf("%s has requested that I say hello to %s", clientID, req.Name)
    return &helloworld.HelloReply{
        Message: fmt.Sprintf("On behalf of %s, hello %s!", clientID, req.Name),
    }, nil
}


