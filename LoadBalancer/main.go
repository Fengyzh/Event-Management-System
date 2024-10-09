package main

import (
	"context"
	"encoding/json"
	"fmt"
	espb "go-lb/grpc/eventservice"
	dspb "go-lb/servicediscovery"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoadBalancer struct {
	mu           sync.Mutex
	Services     []*dspb.Service
	CurrentIndex int
}

type EventJSON struct {
	Name         string   `json:"name"`
	Location     string   `json:"location"`
	Ticketamount int32    `json:"ticketamount"`
	Date         []string `json:"date"`
	Seats        []string `json:"seats"`
}

func (lb *LoadBalancer) GrpctoHTTP(grpcRes any) []byte { 
	jsonres, err := json.Marshal(grpcRes)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	
	return jsonres
}


func (lb *LoadBalancer) pickService() (*dspb.Service, error) {

	lb.mu.Lock()
	defer lb.mu.Unlock()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to discovery service")
	}
	defer conn.Close()

	c := dspb.NewDiscoveryServiceClient(conn)
	services, err := c.GetAllService(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to check service health")

	}

	lb.CurrentIndex += 1

	if lb.CurrentIndex >= len(services.Response) {
		lb.CurrentIndex = 0
	}

	return services.Response[lb.CurrentIndex], nil

}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("in servehttp")

	service, err := lb.pickService()
	if err != nil {
		log.Printf("Failed to fetch a service")
	}

	log.Println(service)

	conn, err := grpc.NewClient(service.Grpcport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := espb.NewEventServiceClient(conn)
	w.Header().Set("Content-Type", "application/json")
	var jsonres []byte

	switch method := req.Method; method {
	case "GET":
		response, err := client.GetAllEvents(context.Background(), nil)
		if err != nil {
			log.Fatalf("error while calling gRPC service: %v", err)
		}
		log.Printf("Response from gRPC service: %v", response)
		jsonres = lb.GrpctoHTTP(response)

	case "POST":
		eventGrpcBody := lb.ReflectHTTP(req)
		response, err := client.CreateEvent(context.Background(), eventGrpcBody)
		if err != nil {
			log.Fatalf("error while calling gRPC service: %v", err)
		}
		log.Printf("Response from gRPC service: %v", response)
		jsonres = lb.GrpctoHTTP(response)
	}

	w.Write(jsonres)


}

func (lb *LoadBalancer) ReflectHTTP(req *http.Request) *espb.EventCreateRequest {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Unable to read body")

	}

	var eventJSON EventJSON
	err = json.Unmarshal(body, &eventJSON)
	if err != nil {
		log.Fatalln("Unable to parse body")
	}

	eventGrpc := &espb.EventCreateRequest{
		Location:     eventJSON.Location,
		Name:         eventJSON.Name,
		Seats:        eventJSON.Seats,
		Date:         eventJSON.Date,
		Ticketamount: eventJSON.Ticketamount,
	}

	return eventGrpc
}

func main() {

	lb := &LoadBalancer{
		Services:     []*dspb.Service{},
		CurrentIndex: 0,
	}

	r := mux.NewRouter()
	r.HandleFunc("/event", lb.ServeHTTP)

	http.Handle("/", r)
	fmt.Println("Load balancer listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
