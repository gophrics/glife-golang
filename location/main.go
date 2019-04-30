package main

import (
	"context"
	"log"
	"net"
	"strconv"

	// Import the generated protobuf code
	lc "./proto"

	redis "../common/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

type service struct {
}

func (s *service) UpdateMyLocation(ctx context.Context, req *lc.Location) (*lc.Response, error) {

	latitude, _ := strconv.ParseFloat(req.Latitude, 64)
	longitude, _ := strconv.ParseFloat(req.Longitude, 64)
	var profileId string = req.ProfileId

	_, err := redis.Instance.GeoAdd(profileId, &redis.GeoLocation{
		Latitude:  latitude,
		Longitude: longitude,
		Name:      "LastKnown",
	}).Result()

	response := &lc.Response{
		OperationSuccess: err == nil,
	}

	return response, nil
}

func (s *service) NearMe(ctx context.Context, req *lc.Location) (*lc.NearMeResponse, error) {

	return nil, nil
}

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
	lc.RegisterLocationServiceServer(s, &service{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
