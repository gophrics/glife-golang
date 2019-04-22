package main

import (
	"context"
	"log"
	"net"

	// Import the generated protobuf code
	lc "./proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

type service struct {
}

func (s *service) GetUser(ctx context.Context, req *lc.User) (*lc.Response, error) {

	response := &lc.Response{
		OperationSuccess: true,
		User: req,
	}

	return response,nil
}

func (s *service) RegisterUser(ctx context.Context, req *lc.User) (*lc.Response, error) {

	response := &lc.Response{
		OperationSuccess: true,
		User: req,
	}

	return response, nil
}

// 
func main() {

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	location.RegisterProfileServiceServer(s, &service{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
