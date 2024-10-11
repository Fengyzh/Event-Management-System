package discoveryservice

import (
	"context"
	"errors"
	dspb "go-lb/servicediscovery"
	"log"
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



func (s *DiscoveryService) ContainsService(service *dspb.Service) (int) {
	for idx, serv := range s.Services {
		if serv.Addr == service.Addr {
			return idx
		}
	}
	return -1
}


func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{Services: make([]*dspb.Service, 0)}
}

func (s *DiscoveryService) RegisterService(ctx context.Context, service *dspb.Service) (*dspb.ServiceResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ServiceIdx := s.ContainsService(service)
	if ServiceIdx >= 0 {
		s.Services[ServiceIdx] = service
		log.Printf("Service at address %s re-registered", service.Addr)
	} else {
		s.Services = append(s.Services, service)
		log.Printf("Service %s registered at %s", service.Name, service.Addr)
	}


	return &dspb.ServiceResponse{Message: "Added Service into list", Status: 200}, nil
}

func (s *DiscoveryService) GetAllService(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceResponseList, error) {

	log.Printf("Service get request")
	var availableServiceList []*dspb.Service

	for _, service := range(s.Services) {
		if service.Available {
			availableServiceList = append(availableServiceList, service)
		}
	}

	var err error
	if len(availableServiceList) == 0 {
		err = errors.New("unable to find available service")
	}

	return &dspb.ServiceResponseList{Response: availableServiceList}, err
}

func (s *DiscoveryService) CheckServiceHealth() {

	//var AllServiceStatus []*dspb.ServiceHealthCheck

	for {
		s.mu.Lock()

		for index, service := range s.Services {
			conn, err := grpc.NewClient(service.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("Failed to check health for %s, at %s", service.Name, service.Addr)
			}

			client := dspb.NewDiscoveryServiceClient(conn)
			res, err := client.HealthCheck(context.Background(), nil)
			//AllServiceStatus = append(AllServiceStatus, res)

			if err != nil || !res.Status {
				log.Printf("Service %s at %s is unhealthy", service.Name, service.Addr)
				s.Services[index].Available = false
			} else {
				log.Printf("Service %s at %s is healthy", *res.Name, res.Addr)
			}
			conn.Close()

		}
		s.mu.Unlock()
		time.Sleep(5 * time.Second)

	}
	//return AllServiceStatus, nil

}
