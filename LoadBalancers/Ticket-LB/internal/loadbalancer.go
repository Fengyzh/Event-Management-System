package loadbalancer

import (
	"context"
	"encoding/json"
	tpb "go-lb/grpc/ticketservice"
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

type TicketCreateJSON struct {
	Eventid  int32  `json:"eventid"`
	Location string `json:"location"`
	Amount   int32  `json:"amount"`
	Date     string `json:"date"`
	Seats    string `json:"seat"`
}

type TicketJSON struct {
	EventId   int32    `json:"eventid"`
	Location  string   `json:"location"`
	Amount    int32    `json:"amount"`
	Date      string `json:"date"`
	Seats     string   `json:"seat"`
	TicketId  int32    `json:"ticketid"`
	Fulfilled bool     `json:"fulfilled"`
}

func NewLoadBalancer() *LoadBalancer {

	return &LoadBalancer{
		Services:     []*dspb.Service{},
		CurrentIndex: 0,
	}
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
	services, err := c.GetServiceByFlag(context.Background(), &dspb.ServiceFlag{Flag: []string{"ticket"}})
	if err != nil {
		log.Printf("Failed to check service health")

	}

	lb.CurrentIndex += 1

	if lb.CurrentIndex >= len(services.Response) {
		lb.CurrentIndex = 0
	}

	return services.Response[lb.CurrentIndex], nil

}

func (lb *LoadBalancer) GetGrpcClient() (tpb.TicketServiceClient, *grpc.ClientConn) {
	service, err := lb.pickService()
	if err != nil {
		log.Printf("Failed to fetch a service")
	}
	//log.Println(service)
	log.Printf("Picked: %s at %s", service.Name, service.Addr)

	conn, err := grpc.NewClient(service.Grpcport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	client := tpb.NewTicketServiceClient(conn)

	return client, conn
}

func (lb *LoadBalancer) CreateTicket(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	eventGrpcBody := lb.ReflectHTTPCreateRequest(req)
	response, err := client.CreateTicket(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}

func (lb *LoadBalancer) GetTickets(w http.ResponseWriter, r *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	response, err := client.GetAllTickets(context.Background(), nil)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}

func (lb *LoadBalancer) GetTicketById(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	vars := mux.Vars(req)
	id := vars["id"]
	tid, _ := strconv.ParseInt(id, 10, 32)

	response, err := client.GetTicket(context.Background(), &tpb.TicketId{Ticketid: int32(tid)})
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)

}

func (lb *LoadBalancer) UpdateTicket(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	eventGrpcBody := lb.ReflectHTTPEvent(req)
	response, err := client.UpdateTicket(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
}

func (lb *LoadBalancer) DeleteTicket(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	vars := mux.Vars(req)
	id := vars["id"]
	tid, _ := strconv.ParseInt(id, 10, 32)

	eventGrpcBody := &tpb.TicketId{Ticketid: int32(tid)}
	response, err := client.DeleteTicket(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)

}

// TODO
/* func (lb *LoadBalancer) OrderTicketEvent(w http.ResponseWriter, req *http.Request) {
	client, conn := lb.GetGrpcClient()
	defer conn.Close()

	vars := mux.Vars(req)
	id := vars["id"]
	eid, _ := strconv.ParseInt(id, 10, 32)

	eventGrpcBody := &espb.EventId{Eventid: int32(eid)}
	response, err := client.OrderEventTicket(context.Background(), eventGrpcBody)
	if err != nil {
		log.Fatalf("error while calling gRPC service: %v", err)
	}
	log.Printf("Response from gRPC service: %v", response)
	jsonres := lb.GrpctoHTTP(response)
	w.Write(jsonres)
} */

func (lb *LoadBalancer) ReflectHTTPCreateRequest(req *http.Request) *tpb.TicketCreateRequest {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Unable to read body")

	}

	var ticketJSON TicketCreateJSON
	err = json.Unmarshal(body, &ticketJSON)
	if err != nil {
		log.Fatalln("Unable to parse body")
	}

	ticketGrpc := &tpb.TicketCreateRequest{
		Location: ticketJSON.Location,
		Eventid:  ticketJSON.Eventid,
		Seat:     ticketJSON.Seats,
		Date:     ticketJSON.Date,
		Amount:   ticketJSON.Amount,
	}

	return ticketGrpc
}

func (lb *LoadBalancer) ReflectHTTPEvent(req *http.Request) *tpb.TicketOrder {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Unable to read body")

	}

	var ticketJSON TicketJSON
	err = json.Unmarshal(body, &ticketJSON)
	if err != nil {
		log.Fatalln("Unable to parse body")
	}
	vars := mux.Vars(req)
	id := vars["id"]
	eid, _ := strconv.ParseInt(id, 10, 32)

	eventGrpc := &tpb.TicketOrder{
		Eventid:   int32(eid),
		Location:  ticketJSON.Location,
		Seat:      ticketJSON.Seats,
		Date:      ticketJSON.Date,
		Amount:    ticketJSON.Amount,
		Fulfilled: ticketJSON.Fulfilled,
		Ticketid:  ticketJSON.TicketId,
	}

	return eventGrpc
}
