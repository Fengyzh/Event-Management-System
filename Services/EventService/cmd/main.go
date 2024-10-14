package main

import (
	"context"
	"database/sql"
	"fmt"
	eventdata "go-lb/Services/EventService/data"
	event "go-lb/Services/EventService/internal"
	eventservice "go-lb/Services/EventService/service"
	espb "go-lb/grpc/eventservice"
	dspb "go-lb/servicediscovery"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
		Grpcport:  "localhost:50053",
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
		dspb.RegisterDiscoveryServiceServer(grpcServer, event.NewDiscoveryService())

		log.Printf("S1 register service listening on :50052")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	const (
		host     = "127.0.0.1"
		port     = 5432
		user     = "postgres"
		password = "password"
		dbname   = "mydatabase"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the repository
	userRepo := eventdata.NewPGEventRepository(db)

	// Create the service and inject the repository
	userService := eventservice.NewEventService(userRepo)

	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		espb.RegisterEventServiceServer(grpcServer, event.NewEventService(userService))

		log.Printf("S1 event service listening on :50053")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	select {}
}
