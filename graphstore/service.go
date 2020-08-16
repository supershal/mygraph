package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/supershal/mygraph/proto"
)

// GraphStore stores graphs
type GraphStore struct {
	graphs map[int64]*pb.Graph
	m      sync.RWMutex
	lastID int64
}

// New returns new instance of server
func New() *GraphStore {
	return &GraphStore{
		graphs: make(map[int64]*pb.Graph, 0),
		lastID: 0,
	}
}

//AddGraph add a cluster of nodes in form of graph
func (s *GraphStore) AddGraph(ctx context.Context, req *pb.AddGraphRequest) (*pb.AddGraphResponse, error) {
	//log.Println("Add Graph request received")
	if req.Graph == nil {
		return nil, fmt.Errorf("Invalid request. Request Does not contain graph")
	}
	var graphID int64
	s.m.Lock()
	graphID = s.lastID
	s.lastID++
	s.graphs[graphID] = req.Graph
	s.m.Unlock()
	log.Println("A graph with ", graphID, "added")
	return &pb.AddGraphResponse{GraphId: graphID}, nil
}

// DeleteGraph deltes graph
func (s *GraphStore) DeleteGraph(ctx context.Context, req *pb.DeleteGraphRequest) (*empty.Empty, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.graphs[req.GraphId]; ok {
		delete(s.graphs, req.GraphId)
		log.Println("Graph ", req.GraphId, "deleted")
		return &empty.Empty{}, nil
	}
	return nil, fmt.Errorf("Graph with id %v does not exists", req.GraphId)
}

// ShortestPath returns shortest path between two nodes if it exists
func (s *GraphStore) ShortestPath(ctx context.Context, req *pb.ShortestPathRequest) (*pb.ShortestPathResponse, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	graphID := req.GraphId
	graph, ok := s.graphs[graphID]
	if !ok {
		return nil, fmt.Errorf("Graph %v does not exists", graphID)
	}
	srcNode, ok := graph.Nodes[req.SourceNodeId]
	if !ok {
		return nil, fmt.Errorf("Source node %v does not exists", req.SourceNodeId)
	}
	dstNode, ok := graph.Nodes[req.DestNodeId]
	if !ok {
		return nil, fmt.Errorf("Destination node %v does not exists", req.DestNodeId)
	}

	queue := []*pb.Node{srcNode}
	parents := map[*pb.Node]*pb.Node{}
	visited := map[*pb.Node]bool{}
	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]
		visited[top] = true
		if top == dstNode {
			break
		}
		for n := range top.Neighbors {
			neighbor := graph.Nodes[n]
			if !visited[neighbor] {
				parents[neighbor] = top
				queue = append(queue, neighbor)
			}
		}
	}

	var path int64
	p := parents[dstNode]
	for p != nil {
		p = parents[p]
		path++
	}

	return &pb.ShortestPathResponse{Distance: path}, nil
}
