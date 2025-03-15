package main

import (
	// "fmt"
	"io"
	// "log"
	"net/http"
	"net/rpc"
	"raft/shared"
	// "time"
)

func main() {
	// create a Membership list
	nodes := shared.NewMembership()
	requests := shared.NewRequests()

	shared.SetRequests(&requests)

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
	// go func() {
	http.ListenAndServe("localhost:9005", nil)
	// }()

	// a little clunky putting this here since its also hosting
	// from here, but we need to figure out where to send reqs from
	// server, _ := rpc.DialHTTP("tcp", "localhost:9005")
	// for i := 1; i == 1; i = 1 {
	// 	var input string
	// 	fmt.Scan(&input)	// this is blocking

	// 	var res shared.PutReply
	// 	var args shared.PutArgs
	// 	args.Value = input
	// 	if err := server.Call("KV.PutKV", args, &res); err != nil {
	// 		fmt.Println("Error: KV.PutKV()", err)
	// 	} else {
	// 		fmt.Printf("Success: server workd: %s", args.Value)
	// 	}

	// shared.handleKV(&input);

	// fmt.Println("You entered:", input)
	// shared.TestPrint()
	// }
}

// type PutArgs struct {
// 	Key   string
// 	Value string
// }
// type PutReply struct {
// }
// type GetArgs struct {
// 	Key string
// }
// type GetReply struct {
// 	Value string
// }

// func getkv(server *rpc.Client, key string) string {
// 	args := GetArgs{key}
// 	reply := GetReply{}
// 	err := server.Call("KV.GetKV", &args, &reply)
// 	if err != nil {
// 		log.Fatal("error:", err)
// 	}
// 	// server.Close()
// 	return reply.Value
// }
// func putkv(server *rpc.Client, key string, val string) {
// 	args := PutArgs{key, val}
// 	reply := PutReply{}
// 	err := server.Call("KV.PutKV", &args, &reply)
// 	if err != nil {
// 		log.Fatal("error:", err)
// 	}
// 	// server.Close()
// }
