package profile

import (
	"context"
	"log"
	"net"

	// Import the generated protobuf code
	pb "./proto/profile"

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
		Created: true,
		User: req,
	}

	return response,nil
}


// 
func Main() {

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
