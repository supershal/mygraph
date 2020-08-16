package graph

import (
	"fmt"
	"sync"

	pb "github.com/supershal/mygraph/proto"
)

type Undirected struct {
	graph *pb.Graph
	m     sync.RWMutex
}

func NewUndirectedGraph() *Undirected {
	return &Undirected{
		graph: &pb.Graph{
			Nodes: make(map[int64]*pb.Node, 0),
		},
	}
}

func MakeUndirectedGraph(g *pb.Graph) *Undirected {
	return &Undirected{
		graph: g,
	}
}

func (u *Undirected) GetGraph() *pb.Graph {
	return u.graph
}
func (u *Undirected) AddNode(id int64) *pb.Node {
	if n, ok := u.graph.Nodes[id]; ok && n != nil {
		return n
	}
	u.m.Lock()
	n := &pb.Node{
		Id:        id,
		Neighbors: make(map[int64]bool),
	}
	u.graph.Nodes[id] = n
	u.m.Unlock()
	return n
}

func (u *Undirected) AddEdge(src, dst *pb.Node) error {
	var ok bool
	if _, ok = u.graph.Nodes[src.Id]; !ok {
		return fmt.Errorf("source node does not exists")
	}
	if _, ok = u.graph.Nodes[dst.Id]; !ok {
		return fmt.Errorf("destination node does not exists")
	}

	u.m.Lock()
	if _, ok := src.Neighbors[dst.Id]; !ok {
		src.Neighbors[dst.Id] = true
	}
	if _, ok := dst.Neighbors[src.Id]; !ok {
		dst.Neighbors[src.Id] = true
	}
	u.m.Unlock()
	return nil
}

func (u *Undirected) String() string {
	graph := u.graph
	graphStr := ""
	for id, node := range graph.Nodes {
		graphStr += fmt.Sprintf("%v->", id)
		for n := range node.Neighbors {
			graphStr += fmt.Sprintf("%v,", n)
		}
		graphStr += fmt.Sprintf("\n")
	}
	return graphStr

}
