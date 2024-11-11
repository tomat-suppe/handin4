package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "handin4/gRPC/handin4/handin4/gRPC/handin4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ListOfAllNodes []Node
var HighestSequenceNumber int64 = 0
var mu sync.Mutex

type Server struct {
	pb.UnimplementedMutualExclusionServer
}

type Node struct {
	pb.UnimplementedMutualExclusionServer
	NodeID     int
	Addr       string
	nextNode   string
	client     pb.MutualExclusionClient
	server     *grpc.Server
	listenAddr string
	SequenceNo int
}

func (node *Node) RequestAccess(ctx context.Context, req *pb.Requester) (*pb.Reply, error) {
	if node.SequenceNo > int(HighestSequenceNumber) || (HighestSequenceNumber == int64(node.SequenceNo)) {
		log.Printf("Node %v has NOT gotten the token and NOT gained access to the Critical Section", node.NodeID)
		return &pb.Reply{NodeID: int64(node.NodeID)}, nil
	} else {
		mu.Lock()
		log.Printf("Node %v has gotten the token and gained access to the Critical Section", node.NodeID)

		time.Sleep(10000)

		log.Printf("Node %v has left the Critical Section, and is passing the token", node.NodeID)
		mu.Unlock()
	}
	return &pb.Reply{NodeID: int64(node.NodeID)}, nil
}

func main() {
	nodeaddr := []string{
		"localhost:50051", // Node 1
		"localhost:50052", // Node 2
		"localhost:50053", // Node 3
	}

	var ListOfAllNodes []Node
	for i := 0; i < len(nodeaddr); i++ {
		nextNode := nodeaddr[(i+1)%len(nodeaddr)] // The next node in the ring
		ListOfAllNodes = append(ListOfAllNodes, Node{
			NodeID:     i,
			Addr:       nodeaddr[i],
			nextNode:   nextNode,
			listenAddr: nodeaddr[i],
		})
	}

	for _, node := range ListOfAllNodes {
		go node.ServerSide()
		go node.ClientSide()

	}
	for {
		for _, node := range ListOfAllNodes {
			go node.NextNode()
		}
	}
}

func (node *Node) ClientSide() {
	conn, err := grpc.NewClient(node.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//node.client = pb.NewMutualExclusionClient(conn)
}

func (node *Node) ServerSide() {
	listener, err := net.Listen("tcp", node.listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Node %v now listening on: %s", node.NodeID, node.listenAddr)
	}
	node.server = grpc.NewServer()
	//server := Server{}
	pb.RegisterMutualExclusionServer(node.server, node)

	log.Printf("Node %v is running on : %s ...", node.NodeID, node.Addr)
	if err := node.server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (node Node) NextNode() {
	conn, err := grpc.NewClient(node.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMutualExclusionClient(conn)
	node.SequenceNo = int(HighestSequenceNumber) + 1
	HighestSequenceNumber = int64(node.SequenceNo)
	client.RequestAccess(context.Background(), &pb.Requester{NodeID: int64(node.NodeID)})
}
