package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func runProxy() error {
	port := 8096
	endpoint := "stats:8095"

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := RegisterStatsHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}

	log.Printf("Starting GRPC proxy on :%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func main() {
	if err := runProxy(); err != nil {
		log.Fatalf("Failed to start GRPC proxy: %v", err)
	}
}
