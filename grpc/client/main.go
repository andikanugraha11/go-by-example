package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andikanugraha11/go-by-example/go-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err.Error())
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	req := &greetpb.GreetRequest{
		Greet: &greetpb.Greet{
			Name: "Gopher",
			Age:  17,
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to access: %v", err.Error())
	}

	fmt.Println(res.GetResult())
}
