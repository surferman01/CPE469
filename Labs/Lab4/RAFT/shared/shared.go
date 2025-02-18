package shared

import (
	"fmt"
	"math/rand"
	// "time"
	// "net/http"
	// "net/rpc"
)

type Role int

const (
	MAX_NODES           = 3
	Follower       Role = 0
	Candidate      Role = 1
	Leader         Role = 2
	START_ELECTION      = "start"
	VOTE                = "vote"
	NEW_LEADER          = "IamUrLeader"
)

// Node struct represents a computing node.
type Node struct {
	ID int

	// Gossip vars
	Hbcounter int
	Time      float64
	Alive     bool

	// Election vars
	Role          Role
	LeaderID      int
	Term          int
	ElectionTimer int
	VotedFor      int
	VoteCount     int
}

// Generate random crash time from 10-60 seconds
// func (n Node) CrashTime() int {
// 	// rand.Seed(time.Now().UnixNano())
// 	max := 60
// 	min := 10
// 	return rand.Intn(max-min) + min
// }

func (n Node) InitializeNeighbors(id int) [2]int {
	neighbor1 := RandInt()
	for neighbor1 == id {
		neighbor1 = RandInt()
	}
	neighbor2 := RandInt()
	for neighbor1 == neighbor2 || neighbor2 == id {
		neighbor2 = RandInt()
	}
	return [2]int{neighbor1, neighbor2}
}

func RandInt() int {
	// rand.Seed(time.Now().UnixNano())
	return rand.Intn(MAX_NODES)
}

/*---------------*/

// Membership struct represents participanting nodes
type Membership struct {
	Members map[int]Node
}

// Returns a new instance of a Membership (pointer).
func NewMembership() *Membership {
	return &Membership{
		Members: make(map[int]Node),
	}
}

// Adds a node to the membership list.
func (m *Membership) Add(payload Node, reply *Node) error {
	// if _, exists := m.Members[payload.ID]; exists {
	// 	return fmt.Errorf("Node with ID %d already exists", payload.ID)
	// }

	m.Members[payload.ID] = payload
	*reply = payload
	return nil
}

// Updates a node in the membership list.
func (m *Membership) Update(payload Node, reply *Node) error {
	// if _, exists := m.Members[payload.ID]; !exists {
	// 	return fmt.Errorf("Node with ID %d doesn't exist", payload.ID)
	// }

	m.Members[payload.ID] = payload
	*reply = payload
	return nil
}

// Returns a node with specific ID.
func (m *Membership) Get(payload int, reply *Node) error {
	if _, exists := m.Members[payload]; !exists {
		return fmt.Errorf("Node with ID %d doesn't exist", payload)
	}

	*reply = m.Members[payload]
	return nil
}

/*---------------*/

type ElectionMSG struct {
	MSG    string
	SRC_ID int
	Term   int
}

// Request struct represents a new message request to a client
type Request struct {
	ID       int
	Table    Membership
	Election ElectionMSG
}

// Requests struct represents pending message requests
type Requests struct {
	GossipPending map[int]Membership
	RAFTPending   map[int]ElectionMSG
}

// Reply struct for gossip and RAFT
type Reply struct {
	Table    Membership
	Election ElectionMSG
}

// Returns a new instance of a Membership (pointer).
func NewRequests() *Requests {
	return &Requests{
		GossipPending: make(map[int]Membership),
		RAFTPending:   make(map[int]ElectionMSG),
	}
}

// Adds a new message request to the pending list
func (req *Requests) Add(payload Request, reply *bool) error {
	req.GossipPending[payload.ID] = payload.Table
	// Only add election request if request isn't blank
	if payload.Election.MSG != "" {
		req.RAFTPending[payload.ID] = payload.Election
	}
	*reply = true
	return nil
}

// Listens to communication from neighboring nodes.
func (req *Requests) Listen(ID int, reply *Reply) error {
	if table, exists := req.GossipPending[ID]; exists {
		reply.Table = table
		delete(req.GossipPending, ID)
	}

	if elect, exists := req.RAFTPending[ID]; exists {
		reply.Election = elect
		delete(req.RAFTPending, ID)
	}

	return nil
}

func CombineTables(table1 *Membership, table2 *Membership) *Membership {
	combined := NewMembership()

	for id, node := range table1.Members {
		combined.Members[id] = node
	}

	for id, node := range table2.Members {
		if existingNode, exists := combined.Members[id]; exists {
			if node.Hbcounter > existingNode.Hbcounter {
				combined.Members[id] = node
			}
		} else {
			combined.Members[id] = node
		}
	}

	return combined
}
