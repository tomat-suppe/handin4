package main

import (
	"flag"
	"log"
	"net"
	"sync"

	pb "handin4/gRPC"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	pb.UnimplementedMutualExclusionServer
	port             int
	participants     []*pb.Node
	mu               sync.Mutex
	LamportTimestamp int64
}

func main() {
	flag.Parse()

	server := &Server{
		port:             *port,
		participants:     []*pb.Node{},
		mu:               sync.Mutex{},
		LamportTimestamp: 0,
	}

	//Sets server to listen for RPCs on its port
	TurnOnServer(server)
}

func TurnOnServer(server *Server) {
	//this method is based on code provided by ChatGPT and also
	//used in handin3

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Now listening on: 50051")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMutualExclusionServer(grpcServer, server)

	log.Println("Node is running on : 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func RequestAccess(Replier *pb.Replier) {
	//broadcastMsg := log.Printf("%s has accepted ----smth----", Replier.RequesterNode.NodeID)
	/*Server.broadcast(
	pb.Reply{
		NodeID: in.GetNodeID(),
		Accepted: in.GetAccepted(),
	})*/
	Reply := (pb.Reply{
		NodeID:   Replier.RequesterNode.GetNodeID(),
		Accepted: true,
	})
	log.Printf("%s has accepted: this is a test", Reply.NodeID)

	//return Reply
}
