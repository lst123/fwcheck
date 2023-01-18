package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/lst123/fwcheck/internal/protobuf"
	"google.golang.org/grpc"
)

var (
	port    = flag.Int("port", 50051, "The server port")
	minPort = 80
	maxPort = 65534
)

type server struct {
	pb.UnimplementedFWCheckServer
}

func genCheck() *pb.ProbReply {
	buf := make([]byte, 4)
	ip := rand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)
	host := fmt.Sprint(net.IP(buf))

	rand.Seed(time.Now().UnixNano())
	port := fmt.Sprint(rand.Intn(maxPort-minPort+1) + minPort)

	return &pb.ProbReply{Ip: host, Port: port}
}

func (s *server) CheckTCP(ctx context.Context, in *pb.ProbRequest) (*pb.ProbReply, error) {
	if in.Result == "none" {
		return genCheck(), nil
	} else {
		log.Printf("Received: %v", in)
		return &pb.ProbReply{}, nil
	}
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
