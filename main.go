package main

import (
	"flag"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gw "./profile/proto/profile"
)

var (
	profileEndpoint = flag.String("profile", "localhost:50051", "endpoint of profile service")
)


func run() error {
	ctx := context.Background();
	ctx, cancel := context.WithCancel(ctx);
	defer cancel();
	 
	mux := runtime.NewServeMux();
	opts := []grpc.DialOption{grpc.WithInsecure()};
	err := gw.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, *profileEndpoint, opts);
	if err != nil {
	  return err
	}
	 
	return http.ListenAndServe(":8080", mux);
}

func main() {
	flag.Parse();
	run();
}
