package main

import (
	"context"
	dspb "go-lb/servicediscovery"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DiscoveryService struct {
	dspb.UnimplementedDiscoveryServiceServer
}

func (ds *DiscoveryService) HealthCheck(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceHealthCheck, error) {

	return &dspb.ServiceHealthCheck{Id: 1, Status: true, Name: "First Service", Addr: "localhost:50052", Message: "MSG from service one"}, nil
}


func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := dspb.NewDiscoveryServiceClient(conn)

	_, err = client.RegisterService(context.Background(), &dspb.Service{
		Name:      "serviceOne",
		Addr:      "localhost:50052",
		Id:        1,
		Available: true,
		Flag:      nil,
	})

	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	dspb.RegisterDiscoveryServiceServer(grpcServer, &DiscoveryService{})

	log.Printf("Service One listening on :50052")


    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC: %v", err)
    }




}
