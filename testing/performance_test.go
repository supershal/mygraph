package testing

import (
	"log"
	"sync"

	"github.com/supershal/mygraph/graph"
)

var perfServerAddr string = "localhost:8082"

func sendGraphParallel(g *graph.Undirected, times int) {
	var wg sync.WaitGroup

	sendGraph := func() int64 {
		defer wg.Done()
		// Create new connection for each operation
		conn, client := NewClientConnection(perfServerAddr)
		defer conn.Close()
		var id int64
		var err error
		//fmt.Println("Adding Graph with ID", i)
		if id, err = AddGraph(client, g.GetGraph()); err != nil {
			log.Fatalln(err)
		}
		log.Println("graph generated:", id)
		return id
	}
	for i := 0; i < times; i++ {
		wg.Add(1)
		go sendGraph()
	}
	wg.Wait()
}

func deleteGraphParallel(times int) {
	var wg sync.WaitGroup

	deleteGraph := func(id int64) {
		defer wg.Done()

		// Create new connection for each operation
		conn, client := NewClientConnection(perfServerAddr)
		defer conn.Close()
		var err error
		if err = DeleteGraph(client, id); err != nil {
			log.Fatalln(err)
		}
		log.Println("graph Deleted:", id)
	}
	for i := 0; i < times; i++ {
		wg.Add(1)
		go deleteGraph(int64(i))
	}
	wg.Wait()
}
