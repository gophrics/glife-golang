package main

import (
	"flag"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pf "./profile/proto"
	lc "./location/proto"
)

var (
	profileEndpoint = flag.String("profile", "localhost:50051", "endpoint of profile service")
	locationEndpoint = flag.String("location", "localhost:50052", "endpoint of location service")
)


func run() error {
	ctx := context.Background();
	ctx, cancel := context.WithCancel(ctx);
	defer cancel();
	 
	mux := runtime.NewServeMux();
	opts := []grpc.DialOption{grpc.WithInsecure()};
	err := pf.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, *profileEndpoint, opts);
	if err != nil {
	  return err
	}

	err = lc.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, *locationEndpoint, opts);
	if err != nil {
		return err
	}
	 
	return http.ListenAndServe(":8080", mux);
}

func main() {
	flag.Parse();
	run();
}
