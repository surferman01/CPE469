package main

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"netRPCstart/shared"
	"os"
	"strconv"

	// "strings"
	"sync"
	"time"
)

const (
	MAX_NODES  = 8
	X_TIME     = 1
	Y_TIME     = 2
	Z_TIME_MAX = 120
	Z_TIME_MIN = 30
)

var repeats = [MAX_NODES + 1]int{}

var self_node shared.Node

// Send the current membership table to a neighboring node with the provided ID
func sendMessage(server *rpc.Client, id int, membership shared.Membership) {
	req := shared.Request{ID: id, Table: membership}
	var reply *bool

	err := server.Call("Requests.Add", req, &reply)
	if err != nil {
		fmt.Println("Error in sendMessage:", err)
		// } else {
		// 	fmt.Println("add request:", reply)
	}
}

// Read incoming messages from other nodes
func readMessages(server *rpc.Client, id int, membership shared.Membership) *shared.Membership {
	var reply *shared.Membership

	err := server.Call("Requests.Listen", id, &reply)
	if err != nil {
		fmt.Println("Error in readMessages:", err)
		// } else {
		// 	fmt.Println("REPLY", *reply)
	}

	// combine memberships here
	// local + reply

	return shared.CombineTables(&membership, reply)
}

func calcTime() float64 {
	return float64(time.Now().Unix())
}

var wg = &sync.WaitGroup{}

func main() {
	// rand.Seed(time.Now().UnixNano())
	Z_TIME := rand.Intn(Z_TIME_MAX-Z_TIME_MIN) + Z_TIME_MIN

	// Connect to RPC server
	server, _ := rpc.DialHTTP("tcp", "localhost:9005")

	args := os.Args[1:]

	// Get ID from command line argument
	if len(args) == 0 {
		fmt.Println("No args given")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Found Error", err)
	}

	fmt.Println("Node", id, "will fail after", Z_TIME, "seconds")

	currTime := calcTime()
	// Construct self
	self_node = shared.Node{ID: id, Hbcounter: 0, Time: currTime, Alive: true}
	var self_node_response shared.Node // Allocate space for a response to overwrite this

	// Add node with input ID
	if err := server.Call("Membership.Add", self_node, &self_node_response); err != nil {
		fmt.Println("Error:2 Membership.Add()", err)
	} else {
		fmt.Printf("Success: Node created with id= %d\n", id)
	}

	neighbors := self_node.InitializeNeighbors(id)
	fmt.Println("Neighbors:", neighbors)

	membership := shared.NewMembership()
	membership.Add(self_node, &self_node)

	sendMessage(server, neighbors[0], *membership)

	// crashTime := self_node.CrashTime()

	time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, &self_node, &membership, id) })
	time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, neighbors, &membership, id) })
	time.AfterFunc(time.Second*time.Duration(Z_TIME), func() { runAfterZ(server, id) })

	wg.Add(1)
	wg.Wait()
}

func runAfterX(server *rpc.Client, node *shared.Node, membership **shared.Membership, id int) {
	// fmt.Println("runAfterX NOW 1")

	// Increment the heartbeat counter
	node.Hbcounter++

	// Update the node's time
	node.Time = calcTime()
	// (*membership).Update(*node, nil)
	(*membership).Members[id] = *node

	// Send the updated node information to the server membership table
	if err := server.Call("Membership.Update", *node, nil); err != nil {
		fmt.Println("Error: Membership.Update()", err)
	} else {
		fmt.Printf("Success: Node %d updated\n", id)
	}

	// Print the current membership table
	printMembership(**membership)

	new_membership := readMessages(server, id, **membership)
	if new_membership != nil {
		for _, n := range (*membership).Members {
			if n.ID == id {
				continue
			}
			tempNode := new_membership.Members[n.ID]

			if (*membership).Members[n.ID].Hbcounter == new_membership.Members[n.ID].Hbcounter {
				repeats[n.ID]++
				// if the heartbeat is the same for 20 checks,
				// then assume it has died
				if repeats[n.ID] >= 20 {
					tempNode.Alive = false
					new_membership.Members[n.ID] = tempNode
				}
			} else {
				repeats[n.ID] = 0
				tempNode.Alive = true
				new_membership.Members[n.ID] = tempNode
			}
		}
		*membership = new_membership
	}

	// temp is now updated membership
	// now update membership

	// below is another method of checking for a node being dead
	// if the timestamp has not updated on the node in x time
	// then assume it is dead
	// this could be used in conjunction with the checking of
	// equivalent heartbeats for further status checking

	// cur_time := calcTime()
	// for id, n := range (**membership).Members {
	// 	if cur_time-n.Time > 20 {
	// 		t := (**membership).Members[id]
	// 		t.Alive = false
	// 		(**membership).Members[id] = t
	// 	}
	// }


	// Schedule the next runAfterX call
	time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, node, membership, id) })
}

func runAfterY(server *rpc.Client, neighbors [2]int, membership **shared.Membership, id int) {
	//TODO
	sel := rand.Intn(2)
	// fmt.Println("runAfterY NOW 2 with sel:", sel)
	// send a heartbeat to a randomly selected neighbor of yours
	sendMessage(server, neighbors[sel], **membership)

	time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, neighbors, membership, id) })
}

func runAfterZ(server *rpc.Client, id int) {
	// this will listening for others
	// readMessages(*server, id, )
	//TODO
	fmt.Println("Node", id, "crashed!")
	os.Exit(0)
}

func printMembership(m shared.Membership) {
	for _, val := range m.Members {
		status := "is Alive"
		if !val.Alive {
			status = "is Dead"
		}
		fmt.Printf("Node %d has hb %d, time %.1f and %s\n", val.ID, val.Hbcounter, val.Time, status)
	}
	fmt.Println("")
	fmt.Println("repeats", repeats)
	fmt.Println("")

}
