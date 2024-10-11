package main

import (
	dspb "go-lb/servicediscovery"
	"log"
	"net"
	discoveryservice "go-lb/DiscoverService/internal"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	discoveryServer := discoveryservice.NewDiscoveryService()

	dspb.RegisterDiscoveryServiceServer(s, discoveryServer)


	go discoveryServer.CheckServiceHealth()

	log.Printf("Server is listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
