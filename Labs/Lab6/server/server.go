package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"log"
	"raft/shared"
)

func main() {
	// create a Membership list
	nodes := shared.NewMembership()
	requests := shared.NewRequests()

	// register nodes with `rpc.DefaultServer`
	rpc.Register(nodes)
	rpc.Register(requests)

	// register an HTTP handler for RPC communication
	rpc.HandleHTTP()

	// sample test endpoint
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "RPC SERVER LIVE!")
	})

	// listen and serve default HTTP server
	// make it run in a goroutine so we can provide
	// user inputs
	go func() {
		http.ListenAndServe("localhost:9005", nil)
	} ()

	// a little clunky putting this here since its also hosting
	// from here, but we need to figure out where to send reqs from
	server, _ := rpc.DialHTTP("tcp", "localhost:9005")

	var input string

	fmt.Scan(&input)
	fmt.Println("You entered:", input)

}

type PutArgs struct {
	Key   string
	Value string
}
type PutReply struct {
}
type GetArgs struct {
	Key string
}
type GetReply struct {
	Value string
}

func getkv(server *rpc.Client, key string) string {
	args := GetArgs{key}
	reply := GetReply{}
	err := server.Call("KV.Get", &args, &reply)
	if err != nil {
		log.Fatal("error:", err)
	}
	// server.Close()
	return reply.Value
}
func putkv(server *rpc.Client, key string, val string) {
	args := PutArgs{key, val}
	reply := PutReply{}
	err := server.Call("KV.Put", &args, &reply)
	if err != nil {
		log.Fatal("error:", err)
	}
	// server.Close()
}

