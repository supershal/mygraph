package graph_test

import (
	"testing"

	"github.com/mygraph/graph"
	pb "github.com/mygraph/proto"
	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {
	u := graph.NewUndirectedGraph()
	for i := 0; i < 5; i++ {
		node := u.AddNode(int64(i))
		assert.Equal(t, int64(i), node.Id, "The two IDs should be same")
		same := u.AddNode(int64(i))
		assert.Equal(t, node, same, "The function should return same node if its already added")
	}
}

func TestAddEdge(t *testing.T) {
	u := graph.NewUndirectedGraph()
	node1 := u.AddNode(1)
	node2 := u.AddNode(2)
	err := u.AddEdge(node1, node2)
	assert.Equal(t, err, nil)
	assert.Equal(t, node1.Neighbors[node2.Id], true, "Edge between node1 and node2 should exists")
	assert.Equal(t, node2.Neighbors[node1.Id], true, "Edge between node2 and node1 should exists")

	node1 = u.AddNode(3)
	node2 = &pb.Node{Id: 5}
	err = u.AddEdge(node1, node2)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "destination node does not exists")

}
