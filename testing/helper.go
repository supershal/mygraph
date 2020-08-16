package testing

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/supershal/mygraph/graph"
	gs "github.com/supershal/mygraph/graphstore"
	pb "github.com/supershal/mygraph/proto"
	"google.golang.org/grpc"
)

// GraphServer test graphsever
func GraphServer(serverAddr string) *grpc.Server {
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g := grpc.NewServer()
	pb.RegisterNetworkServiceServer(g, gs.New())
	fmt.Println("Listening on:", lis.Addr().String())
	if err := g.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
	return g
}

// NewClientConnection new client connection
func NewClientConnection(serverAddr string) (*grpc.ClientConn, pb.NetworkServiceClient) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalln("Unable to connect to server", err)
	}
	return conn, pb.NewNetworkServiceClient(conn)
}

// SampleGraph returns smaple hardcoded graph
func SampleGraph() *graph.Undirected {
	g := graph.NewUndirectedGraph()
	nodes := make([]*pb.Node, 0)
	for i := 1; i < 10; i++ {
		nodes = append(nodes, g.AddNode(int64(i)))
	}
	// create edge between 1 and other nodes
	g.AddEdge(nodes[0], nodes[1])
	g.AddEdge(nodes[1], nodes[2])
	g.AddEdge(nodes[2], nodes[3])
	g.AddEdge(nodes[3], nodes[4])
	g.AddEdge(nodes[4], nodes[5])
	g.AddEdge(nodes[2], nodes[5]) // shortest between 1 to 6 = 3. longest = 5
	return g
}

// AddGraph sends grpc server request to store the graph
func AddGraph(client pb.NetworkServiceClient, g *pb.Graph) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.AddGraphRequest{
		Graph: g,
	}
	res, err := client.AddGraph(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.GetGraphId(), nil
}

// DeleteGraph deltes graph from server
func DeleteGraph(client pb.NetworkServiceClient, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.DeleteGraphRequest{
		GraphId: id,
	}
	_, err := client.DeleteGraph(ctx, req)
	return err
}

// ShortestPath returns shortest path between two nodes in a graph
func ShortestPath(client pb.NetworkServiceClient, id, srcid, dstid int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	req := &pb.ShortestPathRequest{
		GraphId:      id,
		SourceNodeId: srcid,
		DestNodeId:   dstid,
	}

	res, err := client.ShortestPath(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.Distance, nil
}
