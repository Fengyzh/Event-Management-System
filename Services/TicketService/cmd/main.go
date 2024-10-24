package main

import (
	"context"
	"database/sql"
	"fmt"
	ticketdata "go-lb/Services/TicketService/data"
	ticket "go-lb/Services/TicketService/internal"
	ticketService "go-lb/Services/TicketService/service"
	tpb "go-lb/grpc/ticketservice"
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
		Addr:      "localhost:50062",
		Id:        1,
		Available: true,
		Flag:      []string{"ticket"},
		Grpcport:  "localhost:50063",
	})

	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	go func() {
		lis, err := net.Listen("tcp", ":50062")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		dspb.RegisterDiscoveryServiceServer(grpcServer, ticket.NewDiscoveryService())

		log.Printf("T1 register service listening on :50062")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	const (
		host     = "127.0.0.1"
		port     = 5433
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
	userRepo := ticketdata.NewPGTicketRepository(db)

	// Create the service and inject the repository
	userService := ticketService.NewTicketService(userRepo)

	go func() {
		lis, err := net.Listen("tcp", ":50063")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		tpb.RegisterTicketServiceServer(grpcServer, ticket.NewTicketService(userService))

		log.Printf("S1 event service listening on :50063")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	select {}
}
