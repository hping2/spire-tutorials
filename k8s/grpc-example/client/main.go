package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8088", "host:port of the server")
	flag.Parse()

	log.Println("Starting up...")
	log.Println("Server Addr:", addr)

	ctx := context.Background()

	///////////
	//
	//
	client, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
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
