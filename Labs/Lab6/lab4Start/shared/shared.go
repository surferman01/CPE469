package shared

import (
	"fmt"
	"math/rand"
	// "sync"
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
	REPLICAS            = 3
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

	// Consistent Hashing
	Hashes map[string]string
	// Mu     sync.Mutex
}

// Generate random crash time from 10-60 seconds
// func (n Node) CrashTime() int {
// 	// rand.Seed(time.Now().UnixNano())
// 	max := 60
// 	min := 10
// 	return rand.Intn(max-min) + min
// }

// func (n Node) InitializeNeighbors(id int) [2]int {
// 	neighbor1 := RandInt()
// 	for neighbor1 == id {
// 		neighbor1 = RandInt()
// 	}
// 	neighbor2 := RandInt()
// 	for neighbor1 == neighbor2 || neighbor2 == id {
// 		neighbor2 = RandInt()
// 	}
// 	return [2]int{neighbor1, neighbor2}
// }

func TestPrint() {
	fmt.Println("hellllllo")
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
	if _, exists := m.Members[payload.ID]; exists {
		// don't overwrite hashtable
		payload.Hashes = m.Members[payload.ID].Hashes
	}

	m.Members[payload.ID] = payload

	*reply = payload
	return nil
}

// Updates entire membership table except hashes
func (m *Membership) Update(payload Membership, reply *Membership) error {
	for id, payloadNode := range payload.Members {
		if existingNode, exists := m.Members[id]; exists {
			existingNode.Hbcounter = payloadNode.Hbcounter
			existingNode.Time = payloadNode.Time
			existingNode.Alive = payloadNode.Alive
			existingNode.Role = payloadNode.Role
			existingNode.LeaderID = payloadNode.LeaderID
			existingNode.Term = payloadNode.Term
			existingNode.ElectionTimer = payloadNode.ElectionTimer
			existingNode.VotedFor = payloadNode.VotedFor
			existingNode.VoteCount = payloadNode.VoteCount

			m.Members[id] = existingNode
		} else {
			m.Members[id] = payloadNode
		}
	}

	fmt.Println("Updating server membership...")
	printMembership(*m)

	return nil
}

// THIS IS UNUSED
// Returns a node with specific ID.
// func (m *Membership) Get(payload int, reply *Node) error {
// 	if _, exists := m.Members[payload]; !exists {
// 		return fmt.Errorf("Node with ID %d doesn't exist", payload)
// 	}

// 	*reply = m.Members[payload]
// 	return nil
// }

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
	RAFTPending   map[int][]ElectionMSG
}

// Reply struct for gossip and RAFT
type Reply struct {
	Table    Membership
	Election []ElectionMSG
}

// Returns a new instance of a Membership (pointer).
func NewRequests() *Requests {
	return &Requests{
		GossipPending: make(map[int]Membership),
		RAFTPending:   make(map[int][]ElectionMSG),
	}
}

// Adds a new message request to the pending list
func (req *Requests) Add(payload Request, reply *bool) error {
	req.GossipPending[payload.ID] = payload.Table

	// Only add election request if request isn't blank
	if payload.Election.MSG != "" {
		if req.RAFTPending[payload.ID] == nil {
			req.RAFTPending[payload.ID] = make([]ElectionMSG, MAX_NODES)
		}
		req.RAFTPending[payload.ID] = append(req.RAFTPending[payload.ID], payload.Election)
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
			if node.Time > existingNode.Time+10 {
				combined.Members[id] = node
			} else if node.Hbcounter > existingNode.Hbcounter {
				combined.Members[id] = node
			}
		} else {
			combined.Members[id] = node
		}
	}

	return combined
}

// // this will return a hash from 1 to MAX_NODES to tell you
// // which node the key will need to go
// // this will NOT be used as the key on the nodes
// // the key for each value will still be a 20 char string
func hashLocation(key string) int {
	hash := 0
	for i := range len(key) {
		// some arbitrary hashing scheme here
		hash = (hash*31 + int(key[i])) % MAX_NODES
	}
	return hash + 1
}

// KV is one of the aspects of a node

type PutArgs struct {
	Key   string
	Value string
}
type PutReply struct {
	Status string
}

type GetArgs struct {
	Key string
}
type GetReply struct {
	Status string
	Value  string
}

// so whenever the nodes do like a readmessage
// the KV data gets reset
// we need to ensure that the KV data persists
// it may have something to do with th combine
// tables or something else

func checkNode(m *Membership, location int, args *GetArgs) (string, int) {
	for i := 0; i < 3; i++ {
		loc := (location+i+MAX_NODES)%MAX_NODES + 1
		fmt.Println("checking node:", loc)
		fmt.Println("CJECKING MEBRERSHIP")
		printMembership(*m)
		fmt.Println("DONE CHECKING MEMEBRSHIP")
		if node, exists := m.Members[loc]; exists {
			if !node.Alive {
				continue
			}
			if value, ok := node.Hashes[args.Key]; ok {
				fmt.Println("value found @:", loc)
				return value, loc
			}
			return "", loc
		}
	}
	return "", 0
}

var reqCount = 0

func (m *Membership) GetKV(args *GetArgs, reply *GetReply) error {
	reqCount++
	fmt.Println("Request #:", reqCount)
	fmt.Println("Get:", args.Key)
	location := hashLocation(args.Key)

	loc := 0
	reply.Value, loc = checkNode(m, location, args)

	if loc == 0 {
		reply.Status = "no node available"
		return nil
	}
	if reply.Value == "" {
		// fmt.Println("NO VALUE FOUND")
		reply.Status = "not found"
	} else {
		reply.Status = "success"
		// fmt.Println("get:", reply.Value, "location:", loc)
	}

	printMembership(*m)

	return nil
}

func (m *Membership) PutKV(args *PutArgs, reply *PutReply) error {
	reqCount++
	fmt.Println("Request #:", reqCount)
	fmt.Println("Put:", args.Key, args.Value)

	location := hashLocation(args.Key)
	checkArgs := GetArgs{Key: args.Key}
	value := ""
	loc := 0
	value, loc = checkNode(m, location, &checkArgs)
	res := new(bool)

	blankElection := ElectionMSG{MSG: "", SRC_ID: loc, Term: 0}

	req := Request{ID: loc, Table: *m, Election: blankElection}

	if value == "" && loc != 0 {
		for i := 0; i < REPLICAS; i++ {
			idx := (loc+i+MAX_NODES)%MAX_NODES + 1
			m.Members[idx].Hashes[args.Key] = args.Value
		}

		printMembership(*m)
		// send update table to client (loc)
		(*REQUESTS).Add(req, res)
		reply.Status = "success"
	} else if value != "" {
		for i := 0; i < REPLICAS; i++ {
			idx := (loc+i+MAX_NODES)%MAX_NODES + 1
			m.Members[idx].Hashes[args.Key] = args.Value
		}

		printMembership(*m)
		// send update table to client (loc)
		(*REQUESTS).Add(req, res)
		reply.Status = "key exists - successfully overwritten"
	} else {
		reply.Status = "no node available"
	}
	// fmt.Println("put:", args.Key, args.Value, "location:", location)
	return nil
}

func printMembership(m Membership) {
	for _, val := range m.Members {
		status := "is Alive"
		if !val.Alive {
			status = "is Dead"
		}
		fmt.Printf("Node %d has hb %d, time %.1f and %s\n", val.ID, val.Hbcounter, val.Time, status)
		fmt.Println("Hashes:", val.Hashes)
	}
	fmt.Println("")
	// fmt.Println("repeats", repeats)
	// fmt.Println("")

}

// GLOBAL REQUESTS
var REQUESTS **Requests

func SetRequests(r **Requests) {
	REQUESTS = r
}
