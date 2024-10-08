package main

import (
	"context"
	espb "go-lb/grpc/eventservice"
	dspb "go-lb/servicediscovery"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DiscoveryService struct {
	dspb.UnimplementedDiscoveryServiceServer
}

type EventService struct {
	espb.UnimplementedEventServiceServer
	Events []*espb.Event
}

func (ds *DiscoveryService) HealthCheck(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceHealthCheck, error) {

	return &dspb.ServiceHealthCheck{Id: 1, Status: true, Name: "First Service", Addr: "localhost:50052", Message: "MSG from service one"}, nil
}

func (es *EventService) CreateEvent(ctx context.Context, req *espb.EventCreateRequest) (*espb.EventResponse, error) {

	var eid int32 = 1
	newEvent := &espb.Event{Eventid: eid, Ticketamount: req.Ticketamount, Name: req.Name, Location: req.Location, Date: req.Date, Seats: req.Seats}
	es.Events = append(es.Events, newEvent)

	return &espb.EventResponse{Eventid: eid, Message: "Event Created"}, nil
}

func (es *EventService) GetAllEvents(ctx context.Context, _ *espb.Empty) (*espb.EventList, error) {

	return &espb.EventList{Events: es.Events}, nil
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

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		dspb.RegisterDiscoveryServiceServer(grpcServer, &DiscoveryService{})

		log.Printf("S1 register service listening on :50052")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		espb.RegisterEventServiceServer(grpcServer, &EventService{})

		log.Printf("S1 event service listening on :50053")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	select {}
}
