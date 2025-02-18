package main

import (
	"fmt"
	// "math/rand"
	"net/rpc"
	"os"
	"strconv"

	"raft/shared"

	// "strings"
	"sync"
	"time"
)

const (
	X_TIME        = 1
	Y_TIME        = 2
	Z_TIME_MAX    = 120
	Z_TIME_MIN    = 30
	ElectionCount = 10
)

var repeats = [shared.MAX_NODES + 1]int{}

var self_node shared.Node

// Send the current membership table to a neighboring node with the provided ID
func sendMessage(server *rpc.Client, id int, membership shared.Membership, elect shared.ElectionMSG) {
	req := shared.Request{ID: id, Table: membership, Election: elect}
	var reply *bool

	err := server.Call("Requests.Add", req, &reply)
	if err != nil {
		fmt.Println("Error in sendMessage:", err)
	}
}

// Read incoming messages from other nodes
func readMessages(server *rpc.Client, node *shared.Node, membership shared.Membership) *shared.Membership {
	var reply *shared.Reply

	err := server.Call("Requests.Listen", node.ID, &reply)
	if err != nil {
		fmt.Println("Error in readMessages:", err)
	}

	if reply.Election.MSG == shared.START_ELECTION {
		VoteRequest(server, node, reply.Election.SRC_ID, membership)
	} else if reply.Election.MSG == shared.VOTE {
		CountVote(server, node, &membership)
	} else if reply.Election.MSG == shared.NEW_LEADER {
		fmt.Printf("Node %d is the leader for term %d\n", reply.Election.SRC_ID, node.Term)
		node.LeaderID = reply.Election.SRC_ID
		node.ElectionTimer = ElectionCount
	}

	return shared.CombineTables(&membership, &reply.Table)
}

func calcTime() float64 {
	return float64(time.Now().Unix())
}

var wg = &sync.WaitGroup{}

func main() {
	// rand.Seed(time.Now().UnixNano())
	// Z_TIME := rand.Intn(Z_TIME_MAX-Z_TIME_MIN) + Z_TIME_MIN

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

	// Construct self
	currTime := calcTime()
	self_node = shared.Node{ID: id, Hbcounter: 0, Time: currTime, Alive: true, Role: shared.Follower, LeaderID: 0, Term: 0, ElectionTimer: ElectionCount, VotedFor: 0, VoteCount: 0}
	var self_node_response shared.Node // Allocate space for a response to overwrite this

	// Add node with input ID
	if err := server.Call("Membership.Add", self_node, &self_node_response); err != nil {
		fmt.Println("Error:2 Membership.Add()", err)
	} else {
		fmt.Printf("Success: Node created with id= %d\n", id)
	}

	membership := shared.NewMembership()
	membership.Add(self_node, &self_node)

	blankElection := shared.ElectionMSG{MSG: "", SRC_ID: id}

	sendMessage(server, shared.RandInt(), *membership, blankElection)
	sendMessage(server, shared.RandInt(), *membership, blankElection)

	time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, &self_node, &membership, id) })
	time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, shared.RandInt(), shared.RandInt(), &membership, id) })
	// time.AfterFunc(time.Second*time.Duration(Z_TIME), func() { runAfterZ(id) })

	wg.Add(1)
	wg.Wait()
}

func runAfterX(server *rpc.Client, node *shared.Node, membership **shared.Membership, id int) {
	// fmt.Println("runAfterX NOW 1")

	// Increment the heartbeat counter
	node.Hbcounter++

	// Update the node's time
	node.Time = calcTime()

	// Decrement election counter
	if node.Role != shared.Leader {
		node.ElectionTimer--
		if node.ElectionTimer == 0 {
			StartElection(server, node, *membership)
		}
	} else {
		// Tell other nodes you are the leader
		electionOver := shared.ElectionMSG{MSG: shared.NEW_LEADER, SRC_ID: node.ID}
		for id := range (*membership).Members {
			if id != node.ID {
				sendMessage(server, id, **membership, electionOver)
			}
		}
	}

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

	new_membership := readMessages(server, node, **membership)
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

	// Schedule the next runAfterX call
	time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, node, membership, id) })
}

func runAfterY(server *rpc.Client, neighbor1 int, neighbor2 int, membership **shared.Membership, id int) {
	fmt.Println("neightbors:", neighbor1, neighbor2)
	// send a heartbeat to a randomly selected neighbor of yours
	blankElection := shared.ElectionMSG{MSG: "", SRC_ID: id}
	sendMessage(server, neighbor1, **membership, blankElection)
	sendMessage(server, neighbor2, **membership, blankElection)

	time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, shared.RandInt(), shared.RandInt(), membership, id) })
}

// func runAfterZ(id int) {
// 	fmt.Println("Node", id, "crashed!")
// 	os.Exit(0)
// }

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

func StartElection(server *rpc.Client, n *shared.Node, membership *shared.Membership) {
	n.Role = shared.Candidate
	n.Term++
	n.VotedFor = n.ID
	n.VoteCount = 1
	n.ElectionTimer = ElectionCount

	fmt.Printf("Node %d is starting an election for term %d\n", n.ID, n.Term)

	// Request votes from other nodes
	fmt.Println("Starting Election...")
	startElection := shared.ElectionMSG{MSG: shared.START_ELECTION, SRC_ID: n.ID}
	for id := range membership.Members {
		if id != n.ID {
			// send election request
			sendMessage(server, id, *membership, startElection)
		}
	}
}

func VoteRequest(server *rpc.Client, n *shared.Node, src_ID int, membership shared.Membership) {
	n.VotedFor = src_ID

	elect := shared.ElectionMSG{MSG: shared.VOTE, SRC_ID: n.ID}

	fmt.Println("Sending Vote...")
	sendMessage(server, src_ID, membership, elect)
}

func CountVote(server *rpc.Client, n *shared.Node, membership *shared.Membership) {
	n.VoteCount++

	numNodes := len(membership.Members)
	majority := numNodes/2 + 1

	fmt.Println("Received Vote...")

	if n.VoteCount >= majority {
		n.Role = shared.Leader
		n.LeaderID = n.ID
		fmt.Printf("Node %d has won the election and is now the leader for term %d\n", n.ID, n.Term)

		// Tell other nodes you are the leader
		electionOver := shared.ElectionMSG{MSG: shared.NEW_LEADER, SRC_ID: n.ID}
		for id := range membership.Members {
			if id != n.ID {
				sendMessage(server, id, *membership, electionOver)
			}
		}
	}
}
