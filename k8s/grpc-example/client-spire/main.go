package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffegrpc/grpccredentials"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/peer"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8088", "host:port of the server")
	flag.Parse()

	log.Println("Starting up...")
	log.Println("Server Addr:", addr)

	ctx := context.Background()

	log.Println("To call workflow api")
	source, err := workloadapi.NewX509Source(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	serverID := spiffeid.RequireFromString("spiffe://cluster1-demo/greeter-server")
	log.Println("serverID: ", serverID)
	creds := grpccredentials.MTLSServerCredentials(source, source,tlsconfig.AuthorizeID(serverID))

	///////////
	//
	//
	client, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	greeterClient := helloworld.NewGreeterClient(client)

	const interval = time.Second * 10
	log.Println("Issuing requests every %s...", interval)
	for {
		issueRequest(ctx,greeterClient)
		time.Sleep(interval)
	}

}

func issueRequest(ctx context.Context, c helloworld.GreeterClient) {
	p := new(peer.Peer)
	resp, err := c.SayHello(ctx, &helloworld.HelloRequest{
			Name: "Joe",
	}, grpc.Peer(p))

	if err != nil {
		log.Printf("Failed to say hello: %s", err)
		return
	}

	////////////
	//
	///////
	serverID := "SOME-Server"
	if peerID, ok := grpccredentials.PeerIDFromPeer(p); ok {
		serverID = peerID.String()
	}
	log.Printf("%s said %s", serverID, resp.Message) 
}
