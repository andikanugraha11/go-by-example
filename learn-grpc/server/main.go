package main

import (
	"context"
	"log"
	"net"

	userpb "github.com/andikanugraha11/go-by-example/learn-grpc/gen/user/v1"
	"google.golang.org/grpc"
)

type userService struct {
	userpb.UnimplementedUserServiceServer
}

func (u *userService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Uuid:      req.Uuid,
			FirstName: "Andika",
		},
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", "localhost:9897")
	if err != nil {
		log.Fatalf("Failed listen %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userService{})
	grpcServer.Serve(listen)
}
