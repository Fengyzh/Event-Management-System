package main

import (
	"context"
	dspb "go-lb/servicediscovery"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DiscoveryService struct {
	dspb.UnimplementedDiscoveryServiceServer
	Services []*dspb.Service
	mu       sync.Mutex
}

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{Services: make([]*dspb.Service, 0)}
}

func (s *DiscoveryService) RegisterService(ctx context.Context, service *dspb.Service) (*dspb.ServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Services = append(s.Services, service)
	log.Printf("Service %s registered at %s", service.Name, service.Addr)

	return &dspb.ServiceResponse{Message: "Added Service into list", Status: 200}, nil
}

func (s *DiscoveryService) GetAllService(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceResponseList, error) {

	log.Printf("Service get request")
	return &dspb.ServiceResponseList{Response: s.Services}, nil
}

func (s *DiscoveryService) checkServiceHealth() {

	//var AllServiceStatus []*dspb.ServiceHealthCheck

	for {
		s.mu.Lock()

		for _, service := range s.Services {
			conn, err := grpc.NewClient(service.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("Failed to check health for %s, at %s", service.Name, service.Addr)
			}

			client := dspb.NewDiscoveryServiceClient(conn)
			res, err := client.HealthCheck(context.Background(), nil)
			//AllServiceStatus = append(AllServiceStatus, res)

			if err != nil || !res.Status {
				log.Printf("Service %s at %s is unhealthy", res.Name, res.Addr)
			} else {
				log.Printf("Service %s at %s is healthy", res.Name, res.Addr)
			}
			conn.Close()

		}
		s.mu.Unlock()
		time.Sleep(5 * time.Second)

	}
	//return AllServiceStatus, nil

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	discoveryServer := NewDiscoveryService()

	dspb.RegisterDiscoveryServiceServer(s, discoveryServer)

	go discoveryServer.checkServiceHealth()

	log.Printf("Server is listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
