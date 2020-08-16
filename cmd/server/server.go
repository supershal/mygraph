package main

import (
	"flag"
	"fmt"

	"github.com/supershal/mygraph/testing"
)

var port = flag.Int("port", 8080, "The server port")

func main() {
	flag.Parse()
	serverAddr := fmt.Sprintf("localhost:%d", *port)
	g := testing.GraphServer(serverAddr)
	g.Stop()
}
