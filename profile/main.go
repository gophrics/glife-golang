package main

import (
	"context"
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

func (s *service) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.Response, error) {

	// If struct not initialzed, inner variables don't exist
	var user *pb.User = &pb.User{}

	rows, err := dbserver.Instance.Query("SELECT * from User where Id=?", req.Id) //.Scan(&user.Id, &user.Name, &user.Country, &user.Phone)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Country, &user.Phone)
	}
	response := &pb.Response{
		OperationSuccess: true,
		User:             user,
	}

	return response, nil
}

func (s *service) RegisterUser(ctx context.Context, req *pb.User) (*pb.Response, error) {

	prep, _ := dbserver.Instance.Prepare("INSERT INTO User VALUES (?,?,?,?)")
	prep.Exec(req.Id, req.Name, req.Phone, req.Country)

	var user *pb.User
	dbserver.Instance.Query("SELECT * from User")

	// if rows.Next() {
	// 	rows.Scan(&user.Id, &user.Name, &user.Country, &user.Phone)
	// 	fmt.Printf("%s %s %s %s", user.Id, user.Name, user.Country, user.Phone)
	// }

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
