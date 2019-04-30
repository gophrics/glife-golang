package main

import (
	"context"
	"log"
	"net"

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

/*
	TODO: Add Validations
*/
func (s *service) UpdateMyLocation(ctx context.Context, req *lc.Location) (*lc.Response, error) {

	latitude := req.Latitude
	longitude := req.Longitude
	var profileId string = req.ProfileId

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

func (s *service) NearMe(ctx context.Context, req *lc.Location) (*lc.NearMeResponse, error) {

	latitude := req.Latitude
	longitude := req.Longitude
	profileId := req.ProfileId

	res, err := redis.Instance.GeoRadius(profileId, latitude, longitude, &redis.GeoRadiusQuery{
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

		profileArray = append(profileArray, profile)
	}

	response := &lc.NearMeResponse{
		ProfileArray: profileArray,
	}

	return response, nil
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
