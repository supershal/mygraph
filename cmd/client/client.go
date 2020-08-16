package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/supershal/mygraph/testing"
)

var serverAddr = flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")

func main() {
	flag.Parse()

	conn, client := testing.NewClientConnection(*serverAddr)
	defer conn.Close()

	g := testing.SampleGraph()
	fmt.Println("Adding Graph:\n", g.String())

	id, err := testing.AddGraph(client, g.GetGraph())
	if err != nil {
		log.Fatalln("Unable to add Graph", err)
	}
	fmt.Println("Graph added. ID :", id)

	//hardcoded ids just as examples
	if path, err := testing.ShortestPath(client, id, 1, 6); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("Shortest path between node 1 and 6 : ", path)
	}

	if err := testing.DeleteGraph(client, id); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Graph", id, "deleted")
}
