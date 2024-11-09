//any node can request access at any time, but only one can gain it
//the ones that do not gain access need to wait and eventually gain access

//use gRPC calls for message passing between nodes

//start system with at least 3 nodes

//make logs showing state of program

//Be mindful of the last 'Note' about service discovery. Either start the program with 3 nodes
//or hardcode the nodes.

/*package main

import (
	//"log"
	"log"

	pb "handin4/gRPC"

	"google.golang.org/grpc"
)
var N int = 3 //number of Nodes, could potentially take an input and thus vary number of Nodes
var SequenceNumber int = 0 //the number of sequences that have been started or completed
var ListOfAllNodes List
var Availability bool = true //whether Critical Section is available


func main() {
	for i := 0; i < 3; i++ {
		node := *pb.Node{
			nodeID: i,
		}
		ListOfAllNodes.add(node)
	}
	for node, *pb.Node := range ListOfAllNodes{ //figure out a breakpoint?
		go Request(node)
	}
}

func (Requester *pb.Requester) InitializeNode(node *pb.Node) {
	List<*pb.Node> RecipientsList
	for node2, *pb.Node := range RecipientsList{
		if node2.nodeID != node.nodeID{
			RecipientsList.add(node2)
		}
	}
	Requester = *pb.Requester{
		Node: node,
		Recipients: RecipientsList,
		Timestamp: SequenceNumber + 1,
	}
	return Requester
}

func Request(node *pb.Node) {
	stream = pb.RequestAccess(InitializeNode(node))
	log.Printf("Requesting: ", node.nodeID, " is requesting access to the Critical Section")
	Acceptance = 0
	for index, Reply := range stream {
		if !Reply.Accepted{
			log.Printf("Denied: ", node.nodeID, " cannot acccess the Critical Section and will now wait")
			await *pb.Reply.Accepted
		} else {
			Acceptance++
		}
		if Acceptance == 2 {
			lock
			Availability = false
			log.Printf("Accessing: ", node.nodeID, " has acccessed the Critical Section")
			Availability = true
			unlock
		}
	}
}

func (*pb.Reply) Reply(node *pb.Node) {
	if Requester.Timestamp > someother timestamp {
		defer reply
	} else {
		reply affirmative
	}

	return *pb.Reply {
		Reply.nodeID: node.nodeID,
		Reply.Accepted: Availability,
		Reply.Recipients: InitializeNode(node),
	}
}*/