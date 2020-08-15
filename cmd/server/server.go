package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	gs "github.com/mygraph/graphstore"
	pb "github.com/mygraph/proto"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8080, "The server port")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNetworkServiceServer(grpcServer, gs.New())
	fmt.Println("Listening on:", lis.Addr().String())
	grpcServer.Serve(lis)
}
