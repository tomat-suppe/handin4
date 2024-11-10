//any node can request access at any time, but only one can gain it
//the ones that do not gain access need to wait and eventually gain access

//use gRPC calls for message passing between nodes

//start system with at least 3 nodes

//make logs showing state of program

//Be mindful of the last 'Note' about service discovery. Either start the program with 3 nodes
//or hardcode the nodes.

package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "handin4/gRPC"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var N int = 3              //number of Nodes, could potentially take an input and thus vary number of Nodes
var SequenceNumber int = 0 //the number of sequences that have been started or completed
var HighestSequenceNumber int64 = 0
var Acceptance int = 0
var ListOfAllNodes [3]*pb.Node
var Availability bool = true //whether Critical Section is available

var thisNode *pb.Node

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMutualExclusionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		thisNode := &pb.Node{
			NodeID:          int64(i),
			SeqeuenceNumber: HighestSequenceNumber,
		}
		ListOfAllNodes[i] = thisNode
	}
	for i, node := range ListOfAllNodes { //figure out a breakpoint?
		if node.NodeID == thisNode.NodeID {
			break
		}
		node.SeqeuenceNumber = HighestSequenceNumber + 1
		//go Request(node, c)
		log.Printf("Requesting: ", node.NodeID, " is requesting access to the Critical Section")

		HighestSequenceNumber = max(HighestSequenceNumber, node.SeqeuenceNumber)
		if HighestSequenceNumber > node.SeqeuenceNumber || (HighestSequenceNumber == node.SeqeuenceNumber && node.NodeID > thisNode.NodeID) {
			log.Printf("Denied: ", node.NodeID, " cannot acccess the Critical Section and will now wait")
			for {

			}
		} else {
			Acceptance++
			Replier := &pb.Replier{
				RequesterNode:  thisNode,
				SequenceNumber: thisNode.SeqeuenceNumber,
			}
			c.RequestAccess(ctx, Replier) //this returns reply for all repliers
		}

		if Acceptance == 2 {
			log.Printf("%v has entered the Critical Section and I is %v ", node.NodeID, i)
		}
	}
}

/*func (Requester *pb.Requester) InitializeNode(node *pb.Node) {
	List<*pb.Node> RecipientsList
	for node2, *pb.Node := range RecipientsList{
		if node2.nodeID != node.nodeID{
			RecipientsList.add(node2)
		}
	}
	Requester = *pb.Requester{
		Node: node,
		Recipients: RecipientsList,
		//Timestamp: SequenceNumber + 1,
	}
	return Requester
}*/

/*func Request(node *pb.Node, c Connection) {
	log.Printf("Requesting: ", node.NodeID, " is requesting access to the Critical Section")

	HighestSequenceNumber = max(HighestSequenceNumber, node.SeqeuenceNumber)
	if HighestSequenceNumber > node.SeqeuenceNumber || (HighestSequenceNumber == node.SeqeuenceNumber && node.NodeID > thisNode.NodeID) {
		log.Printf("Denied: ", node.NodeID, " cannot acccess the Critical Section and will now wait")
		defer Request(node, c)
	} else {
		Acceptance++
		c.RequestAccess(thisNode) //this returns reply for all repliers
	}
}*/
