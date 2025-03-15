package main

import (
	"fmt"
	// "log"
	"net/rpc"
	"raft/shared"
)

// func connect() *rpc.Client {
// 	client, err := rpc.Dial("tcp", ":1234")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}
// 	return client
// }
// func get(key string) string {
// 	client := connect()
// 	args := shared.GetArgs{key}
// 	reply := shared.GetReply{}
// 	err := client.Call("KV.Get", &args, &reply)
// 	if err != nil {
// 		log.Fatal("error:", err)
// 	}
// 	client.Close()
// 	return reply.Value
// }
// func put(key string, val string) {
// 	client := connect()
// 	args := shared.PutArgs{key, val}
// 	reply := shared.PutReply{}
// 	err := client.Call("KV.Put", &args, &reply)
// 	if err != nil {
// 		log.Fatal("error:", err)
// 	}
// 	client.Close()
// }

func main() {
	server, _ := rpc.DialHTTP("tcp", "localhost:9005")
	// args := os.Args[1:]
	// fmt.Scan("Make Requests (get | put):")

	for i := 1; i == 1; i = 1 {
		// var input string
		// fmt.Scan(&input)	// this is blocking

		var reqType, key, value string
		fmt.Scan(&reqType, &key, &value)

		// fmt.Printf("reqType: %s\nkey: %s\nvalue: %s\n", reqType, key, value)

		if reqType == "get" {
			var res shared.GetReply
			var args shared.GetArgs
			args.Key = key
			if err := server.Call("Membership.GetKV", &args, &res); err != nil {
				fmt.Println("Error: Membership.GetKV()", err)
			} else {
				fmt.Println(res.Status, res.Value)
			}
		} else if reqType == "put" {
			var res shared.PutReply
			var args shared.PutArgs
			args.Key = key
			args.Value = value
			if err := server.Call("Membership.PutKV", &args, &res); err != nil {
				fmt.Println("Error: Membership.PutKV()", err)
			} else {
				fmt.Println(res.Status)
				// fmt.Printf("Success: [%s:%s]\n", args.Key, args.Value)
			}
		}
	}
}
