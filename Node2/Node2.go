package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	pb "handin4/gRPC/handin4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ListOfAllNodes [3]*pb.Node
var HighestSequenceNumber int64 = 0

type Server struct {
	pb.UnimplementedMutualExclusionServer
}

func main() {
	addr := ":50053"

	go ServerSide(addr)

	reader := bufio.NewReader(os.Stdin)
	peerAddr, _ := reader.ReadString('\n')
	peerAddr = peerAddr[:len(peerAddr)-1] // Remove newline character

	log.Println("Enter 'request' to request access to Critical Section")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1] // Remove newline character

	//above 6 lines by help of chatgpt

	ClientSide(addr, input)
}

func ClientSide(addr string, input string) {
	addrconn := flag.String("addr", "localhost:50053", "the address to connect to")
	flag.Parse()
	conn, err := grpc.NewClient(*addrconn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMutualExclusionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	for i := 0; i < 3; i++ {
		thisNode := &pb.Node{
			NodeID:          int64(i),
			SeqeuenceNumber: HighestSequenceNumber,
			Acceptance:      0,
		}
		ListOfAllNodes[i] = thisNode
	}

	if input == "request" {
		for i, node := range ListOfAllNodes { //figure out a breakpoint?
			RequesterNode := ListOfAllNodes[rand.Intn(len(ListOfAllNodes))]
			//RequesterNode := ListOfAllNodes[RequesterNodeIndex]
			RequesterNode.SeqeuenceNumber = HighestSequenceNumber + 1
			log.Println("Requesting: ", RequesterNode.NodeID, " is requesting access to the Critical Section")
			HighestSequenceNumber = max(HighestSequenceNumber, RequesterNode.SeqeuenceNumber)
			//check placement of above statement

			PairofNodes := &pb.PairofNodes{
				Node1: RequesterNode,
				Node2: node,
			}
			req, err := c.RequestAccess(ctx, PairofNodes)
			if err != nil {
				log.Printf("Failure: %v Could not request access", RequesterNode)
			}
			log.Printf("%v has answered for request: %v", node, req.GetAcceptance())
			if RequesterNode.Acceptance == 2 {
				log.Printf("%v has entered the Critical Section and I is %v ", node.NodeID, i)
			}
			RequesterNode.Acceptance = 0
		}
	}
	defer cancel()

}
func ServerSide(addr string) {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Now listening on: 50053")
	}

	grpcServer := grpc.NewServer()
	server := &Server{}
	pb.RegisterMutualExclusionServer(grpcServer, server)

	log.Println("Node is running on : 50053...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *Server) RequestAccess(ctx context.Context, PairofNodes *pb.PairofNodes) (*pb.Reply, error) {
	Node1 := PairofNodes.Node1
	NodeID1 := PairofNodes.Node1.NodeID
	//Node2 := PairofNodes.Node2
	NodeID2 := PairofNodes.Node2.NodeID

	if HighestSequenceNumber > Node1.SeqeuenceNumber || (HighestSequenceNumber == Node1.SeqeuenceNumber && NodeID1 > NodeID2) {
		log.Println("Denied: ", NodeID1, " cannot acccess the Critical Section and will now wait")
		Reply := &pb.Reply{
			Acceptance: false,
		}
		return Reply, nil
	} else {
		Node1.Acceptance++
		Reply := &pb.Reply{
			Acceptance: true,
		}
		return Reply, nil
	}
}
