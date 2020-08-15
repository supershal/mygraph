package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mygraph/graph"
	pb "github.com/mygraph/proto"
	"google.golang.org/grpc"
)

var serverAddr = flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")

func addGraph(g *pb.Graph) (int64, error) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewNetworkServiceClient(conn)
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

func deleteGraph(id int64) error {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewNetworkServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := &pb.DeleteGraphRequest{
		GraphId: id,
	}
	_, err = client.DeleteGraph(ctx, req)
	return err
}

func shortestPath(id, srcid, dstid int64) error {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewNetworkServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := &pb.ShortestPathRequest{
		GraphId:      id,
		SourceNodeId: srcid,
		DestNodeId:   dstid,
	}
	res, err := client.ShortestPath(ctx, req)
	if err != nil {
		return err
	}
	fmt.Println("Shortest path between node", srcid, "and", dstid, "in graph ", id, "is", res.Distance)
	return nil
}

func main() {
	flag.Parse()
	g := sampleGraph()
	var wg sync.WaitGroup

	sendGraph := func() int64 {
		defer wg.Done()
		var id int64
		var err error
		//fmt.Println("Adding Graph with ID", i)
		if id, err = addGraph(g.GetGraph()); err != nil {
			log.Fatalln(err)
		}
		log.Println("graph generated:", id)
		return id
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go sendGraph()
	}

	wg.Wait()
	err := deleteGraph(3)
	if err != nil {
		log.Fatalln(err)
	}
	g.String()
	err = shortestPath(1, 1, 6)
	if err != nil {
		log.Fatalln(err)
	}
}

func sampleGraph() *graph.Undirected {
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
