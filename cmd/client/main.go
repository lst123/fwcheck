package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/lst123/fwcheck/internal/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func connect(host string, port string) (string, error) {
	log.Printf("Trying to open connection to host: %s, port: %s\n", host, port)

	timeout := time.Second * 3
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return "closed", err
	}
	if conn != nil {
		defer conn.Close()
	}
	fmt.Println("Opened", net.JoinHostPort(host, port))
	return "opened", nil
}

func clientCheck() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFWCheckClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := c.CheckTCP(ctx, &pb.ProbRequest{Result: "none"})
	if err != nil {
		log.Fatalf("can't get a task from server: %v", err)
		return
	}

	ip, port := r.GetIp(), r.GetPort()
	res, _ := connect(ip, port)
	if err != nil {
		log.Fatalf("can't connect: %v", err)
		return
	}

	_, err = c.CheckTCP(ctx, &pb.ProbRequest{Ip: ip, Port: port, Result: res})
	if err != nil {
		log.Fatalf("can't send data to the server: %v", err)
		return
	}
	log.Printf("Port is %s", res)
}

func main() {
	flag.Parse()
	for {
		go clientCheck()
		time.Sleep(time.Minute)
	}
}
