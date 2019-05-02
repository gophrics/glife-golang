package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	// Import the generated protobuf code
	lc "./proto"

	redis "../common/redis"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port   = ":50052" // gRPC
	wsport = ":60061" // Websocket
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
		time.Sleep(time.Second)
		profileArray = append(profileArray, profile)
	}

	/*
		response := &lc.NearMeResponse{
			ProfileArray: profileArray,
		}
	*/
	return nil
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

/*
	Websocket API
*/
func NearMeWS(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Attempting connection to websocket")

	// upgrade this connection to a WebSocket
	// connection
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
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

	http.HandleFunc("/location/v1/nearmews", NearMeWS)
	http.ListenAndServe(wsport, nil)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
