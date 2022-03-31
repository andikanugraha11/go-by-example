package main

import (
	"context"
	"fmt"
	"log"

	userpb "github.com/andikanugraha11/go-by-example/learn-grpc/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()
	client := userpb.NewUserServiceClient(conn)

	res, err := client.GetUser(context.Background(), &userpb.GetUserRequest{
		Uuid: "Test",
	})
	if err != nil {
		log.Fatalf("fail to get user: %v", err)
	}

	fmt.Printf("%+v\n", res)
}
