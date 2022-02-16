package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/andikanugraha11/go-by-example/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	name := req.GetGreet().GetName()
	age := req.GetGreet().GetAge()

	msg := fmt.Sprintf("Hello %v your age is %v", name, age)

	res := &greetpb.GreetResponse{
		Result: msg,
	}
	return res, nil
}

func main() {
	lstn, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err.Error())
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lstn); err != nil {
		log.Fatalf("Failed to serve: %v", err.Error())
	}
}
