package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	lc "./proto"

	redis "../common/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052" // gRPC
)

type service struct {
}

/*
	gRPC API
	TODO: Add Validations
*/
func (s *service) UpdateMyLocation(ctx context.Context, req *lc.Location) (*lc.Response, error) {

	latitude := req.Latitude
	profileId := req.ProfileId
	longitude := req.Longitude

	_, err := redis.Instance.GeoAdd("LastKnown", &redis.GeoLocation{
		Latitude:  latitude,
		Longitude: longitude,
		Name:      profileId,
	}).Result()

	response := &lc.Response{
		OperationSuccess: err == nil,
	}

	return response, nil
}

/*
	gRPC API
*/
func (s *service) NearMe(req *lc.Location, stream lc.LocationService_NearMeServer) error {

	fmt.Printf("NearMe HIT")
	latitude := req.Latitude
	longitude := req.Longitude

	res, err := redis.Instance.GeoRadius("LastKnown", latitude, longitude, &redis.GeoRadiusQuery{
		// Update Radius! WTF
		Radius:    10000,
		WithCoord: true,
	}).Result()

	var profileArray []*lc.Location

	if err != nil {
		panic(err.Error())
	}

	for _, element := range res {
		profile := &lc.Location{
			ProfileId: element.Name,
			Latitude:  element.Latitude,
			Longitude: element.Longitude,
		}
		stream.Send(profile)
		profileArray = append(profileArray, profile)
	}

	return nil
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
