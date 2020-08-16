package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var funcServerAddr string = "localhost:8081"

// TestMain intitalizes grpc server
func TestMain(m *testing.M) {
	go func() {
		g := GraphServer(funcServerAddr)
		defer g.Stop()
	}()
	os.Exit(m.Run())
}

func TestAddGraph(t *testing.T) {
	conn, client := NewClientConnection(funcServerAddr)
	defer conn.Close()

	g := SampleGraph()
	cases := []struct {
		id  int64
		err error
	}{
		{0, nil},
		{1, nil},
		{2, nil},
		{3, nil},
	}

	for _, c := range cases {
		id, err := AddGraph(client, g.GetGraph())
		assert.Equal(t, c.err, err, "Error while adding graph")
		assert.Equal(t, c.id, id, "Graph id sequence should match")
		t.Log("Graph added. ID:", id)
	}

}

func TestShortestPath(t *testing.T) {
	conn, client := NewClientConnection(funcServerAddr)
	defer conn.Close()

	g := SampleGraph()
	cases := []struct {
		srcNodeId int64
		dstNodeId int64
		path      int64
		err       error
	}{
		{1, 6, 3, nil},
		{1, 4, 3, nil},
	}

	id, err := AddGraph(client, g.GetGraph())
	assert.Equal(t, err, nil, "Error while adding graph")

	for _, c := range cases {
		path, err := ShortestPath(client, id, c.srcNodeId, c.dstNodeId)
		assert.Equal(t, c.err, err, "Error while finding shortest path ")
		assert.Equal(t, c.path, path, "Incorrect shortest path")
		t.Log("Shortest path between", c.srcNodeId, "and", c.dstNodeId, ":", path)
	}
}

func TestDeleteGraph(t *testing.T) {
	conn, client := NewClientConnection(funcServerAddr)
	defer conn.Close()

	g := SampleGraph()
	cases := []struct {
		id  int64
		err error
	}{
		{0, nil},
		{0, fmt.Errorf("Graph with id %v does not exists", 0)},
	}

	_, err := AddGraph(client, g.GetGraph())
	assert.Equal(t, err, nil, "Error while adding graph")

	for _, c := range cases {
		err := DeleteGraph(client, c.id)
		if c.err != nil {
			assert.NotNil(t, err, "Error while deleting graph does not match with server")
		} else {
			assert.Equal(t, c.err, err, "Error while deleting graph")
		}
	}
}
