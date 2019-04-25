package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	dbserver "../common"
	pb "./proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type service struct {
}

func (s *service) GetUser(ctx context.Context, req *pb.User) (*pb.Response, error) {

	response := &pb.Response{
		OperationSuccess: true,
		User:             req,
	}

	return response, nil
}

func (s *service) RegisterUser(ctx context.Context, req *pb.User) (*pb.Response, error) {

	insert, err := dbserver.Instance.Query(fmt.Sprintf("INSERT INTO User VALUES (%s,%s,%s,%s)", req.Id, req.Name, req.Phone, req.Country))

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	var user *pb.User
	err = dbserver.Instance.QueryRow("SELECT * from Users").Scan(&user.Id, &user.Name, &user.Phone, &user.Country)

	response := &pb.Response{
		OperationSuccess: true,
		User:             user,
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
	pb.RegisterProfileServiceServer(s, &service{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
