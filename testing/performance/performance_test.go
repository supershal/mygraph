package performance

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supershal/mygraph/graph"
	helper "github.com/supershal/mygraph/testing"
)

var perfServerAddr string = "localhost:8081"

// TestMain intitalizes grpc server
func TestMain(m *testing.M) {
	go func() {
		g := helper.GraphServer(perfServerAddr)
		defer g.Stop()
	}()
	os.Exit(m.Run())
}
func TestAddGraph_1000(t *testing.T) {
	g := helper.SampleGraph()
	ids := sendGraphParallel(g, 1000)
	assert.Equal(t, len(ids), 1000, "Graph service should return 1000 graphs")
	graphs := map[int64]bool{}
	for _, id := range ids {
		graphs[id] = true
	}
	for i := 0; i < 1000; i++ {
		_, ok := graphs[int64(i)]
		assert.Equal(t, ok, true, fmt.Sprintf("Graph with id %v should be returned", i))
	}
}

func TestDeleteGraph_1000(t *testing.T) {
	g := helper.SampleGraph()
	ids := sendGraphParallel(g, 1000)
	assert.Equal(t, len(ids), 1000, "Graph service should return 1000 graphs")
	deleteGraphParallel(ids)
}

func sendGraphParallel(g *graph.Undirected, times int) []int64 {
	var wg sync.WaitGroup
	var m sync.Mutex
	graphIds := make([]int64, 0)

	sendGraph := func() int64 {
		defer wg.Done()
		// Create new connection for each operation
		conn, client := helper.NewClientConnection(perfServerAddr)
		defer conn.Close()
		var id int64
		var err error
		//fmt.Println("Adding Graph with ID", i)
		if id, err = helper.AddGraph(client, g.GetGraph()); err != nil {
			log.Fatalln(err)
		}
		log.Println("graph generated:", id)
		m.Lock()
		graphIds = append(graphIds, id)
		m.Unlock()
		return id
	}
	for i := 0; i < times; i++ {
		wg.Add(1)
		go sendGraph()
	}
	wg.Wait()
	return graphIds
}

func deleteGraphParallel(ids []int64) {
	var wg sync.WaitGroup

	deleteGraph := func(id int64) {
		defer wg.Done()

		// Create new connection for each operation
		conn, client := helper.NewClientConnection(perfServerAddr)
		defer conn.Close()
		var err error
		if err = helper.DeleteGraph(client, id); err != nil {
			log.Fatalln(err)
		}
		log.Println("graph Deleted:", id)
	}
	for _, id := range ids {
		wg.Add(1)
		go deleteGraph(id)
	}
	wg.Wait()
}
