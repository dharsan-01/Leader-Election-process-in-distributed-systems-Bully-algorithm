package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	id                 int
	isLeader           bool
	electionStarted    bool
	electionInProgress bool
	mu                 sync.Mutex
}

// newNode creates a new node with a random ID
func newNode() *Node {
	return &Node{
		id: rand.Intn(1000),
	}
}

// startElection initiates the leader election process
func (n *Node) startElection() {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.electionInProgress {
		fmt.Printf("Node %d started an election\n", n.id)
		n.electionStarted = true
		n.electionInProgress = true
		go n.electLeader()
	}
}

// electLeader runs the leader election algorithm
func (n *Node) electLeader() {
	defer func() {
		n.mu.Lock()
		n.electionInProgress = false
		n.mu.Unlock()
	}()

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) 
	highestID := n.id
	isHighestNode := true

	// broadcast election message to all nodes with higher IDs
	for _, node := range getNodesWithHigherIDs(n.id) {
		response := node.receiveElectionMessage(n.id)
		if response {
			isHighestNode = false
			if node.id > highestID {
				highestID = node.id
			}
		}
	}

	if isHighestNode {
		n.becomeLeader()
	} else {
		n.electionStarted = false
	}
}

// receiveElectionMessage handles an election message from another node
func (n *Node) receiveElectionMessage(senderID int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.id > senderID {
		// I have a higher ID so I'll participate in the election
		if !n.electionStarted {
			fmt.Printf("Node %d is sending an election message to Node %d\n", senderID, n.id)
			go n.startElection()
		}
		return true
	}
	// I have a lower ID so I'll defer to the sender
	return false
}

// becomeLeader makes the node the leader
func (n *Node) becomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Printf("Node %d became the leader\n", n.id)
	n.isLeader = true
	n.electionStarted = false
}

// getNodesWithHigherIDs returns a list of nodes with IDs higher than the given ID
func getNodesWithHigherIDs(id int) []*Node {
	// this is a sample node
	nodes := []*Node{
		{id: 100},
		{id: 200},
		{id: 300},
		{id: 400},
		{id: 500},
	}

	// simulating node failure for node with ID 500
	var higherNodes []*Node
	for _, node := range nodes {
		if node.id > id && node.id != 500 {
			higherNodes = append(higherNodes, node)
		}
	}
	return higherNodes
}

func main() {
	node := newNode()
	fmt.Printf("Created node with ID %d\n", node.id)

	node.startElection()
	time.Sleep(5 * time.Second) 
}
