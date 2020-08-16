

# MyGraph
  A sample grpc server-client application to store graphs and perform operation on a stored graph

## Requirements
- Create service using GRPC to perform graph operations
- Store Graph and return unique id
- Find Shortest path between two nodes in a graph
- Delete a graph
- Service should accept multiple concurrent client connections and prevent data corruptions because of race conditions
- Sample Unit test
- Sample funcional test
- Sample performance test

## Build, Test and install
```shell
make all
```
This will create two binary files `graphserver` and `graphclient` in the root directory.

## Build
```shell
    make build
```

## Test
*All tests: Runs unit, functional ans performance tests*
```shell
make test
```

*sample Unit test:*
```shell
go test -v graph/*.go
```

*functional tests:*
```shell
go test -v testing/functional/*.go
```

*performance test:*
Adds and deletes 1000 graphs concurrently 
```shell
go test -v testing/performance/*.go
```

# Run Server and Client(s)
Following instructions creates binary and runs grpc server and client binary to perform operation on sample hardcoded graph in code.
- Create sample hardcoded graph
- Print Graph
- Store graph on server in memory
- Delete the graph from server

*create server and client binaries*
```shell
make bin
```
Run server in a terminal with default port 8080.
```shell
./graphserver
```
OR

With custom port
```shell
./graphserver 8081
```

Run client in a terminal with default address  localhost:8080.
```shell
./graphclient
```
OR

With custom address
```shell
./graphclient localhost:8081
```

### Sample output
*Server*
```shell
> ./graphserver
Listening on: 127.0.0.1:8080
2020/08/16 10:54:53 A graph with  0 added
2020/08/16 10:54:53 Graph  0 deleted
```

*Client*
```shell
> ./graphclient
Adding Graph:
 1->2,
5->4,6,
6->5,3,
8->
2->3,1,
3->2,4,6,
4->3,5,
7->
9->

Graph added. ID : 0
Shortest path between node 1 and 6 :  3
Graph 0 deleted
```

## Enhancements
- Support TLS
- Support more graph operations

## References:
proto3 specs: https://developers.google.com/protocol-buffers/docs/proto3
Go turorials: https://developers.google.com/protocol-buffers/docs/gotutorial, https://grpc.io/docs/languages/go/basics/