package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/lst123/fwcheck/src/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedFWCheckServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) CheckPort(ctx context.Context, in *pb.ProbRequest) (*pb.ProbReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.ProbReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFWCheckServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
