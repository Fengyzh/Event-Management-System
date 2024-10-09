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
	"strconv"
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

type EventCreateJSON struct {
	Name         string   `json:"name"`
	Location     string   `json:"location"`
	Ticketamount int32    `json:"ticketamount"`
	Date         []string `json:"date"`
	Seats        []string `json:"seats"`
}

type EventJSON struct {
	EventId int32 `json:"eventid"`
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
		eventGrpcBody := lb.ReflectHTTPCreateRequest(req)
		response, err := client.CreateEvent(context.Background(), eventGrpcBody)
		if err != nil {
			log.Fatalf("error while calling gRPC service: %v", err)
		}
		log.Printf("Response from gRPC service: %v", response)
		jsonres = lb.GrpctoHTTP(response)
	}

	w.Write(jsonres)

}

func (lb *LoadBalancer) GetGrpcClient() (espb.EventServiceClient, *grpc.ClientConn) {
	service, err := lb.pickService()
	if err != nil {
		log.Printf("Failed to fetch a service")
	}
	log.Println(service)

	conn, err := grpc.NewClient(service.Grpcport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	client := espb.NewEventServiceClient(conn)

	return client, conn
}

func (lb *LoadBalancer) CreateEvent(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	eventGrpcBody := lb.ReflectHTTPCreateRequest(req)
	response, err := client.CreateEvent(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}

func (lb *LoadBalancer) GetEvents(w http.ResponseWriter, r *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	response, err := client.GetAllEvents(context.Background(), nil)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}

func (lb *LoadBalancer) GetEventById(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	vars := mux.Vars(req)
	id := vars["id"]
	eid, _ := strconv.ParseInt(id, 10, 32)

	response, err := client.GetEvent(context.Background(), &espb.EventId{Eventid:int32(eid)})
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)

}


func (lb *LoadBalancer) UpdateEvent(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	eventGrpcBody := lb.ReflectHTTPEvent(req)
	response, err := client.UpdateEvent(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}




func (lb *LoadBalancer) ReflectHTTPCreateRequest(req *http.Request) *espb.EventCreateRequest {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Unable to read body")

	}

	var eventJSON EventCreateJSON
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

func (lb *LoadBalancer) ReflectHTTPEvent(req *http.Request) *espb.Event {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Unable to read body")

	}

	var eventJSON EventJSON
	err = json.Unmarshal(body, &eventJSON)
	if err != nil {
		log.Fatalln("Unable to parse body")
	}
	vars := mux.Vars(req)
	id := vars["id"]
	eid, _ := strconv.ParseInt(id, 10, 32)

	eventGrpc := &espb.Event{
		Eventid:	  int32(eid),
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
	r.HandleFunc("/event", lb.GetEvents).Methods("GET")
	r.HandleFunc("/event/{id}", lb.GetEventById).Methods("GET")
	r.HandleFunc("/event", lb.CreateEvent).Methods("POST")
	r.HandleFunc("/event/{id}", lb.UpdateEvent).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Load balancer listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
