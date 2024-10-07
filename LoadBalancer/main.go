package main

import (
	"context"
	"fmt"
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

	proxyReq, err := http.NewRequest(req.Method, service.Addr, req.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	log.Println(proxyReq)

	if err != nil {
		http.Error(w, "Failed to reach backend server", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {

	lb := &LoadBalancer{
		Services:     []*dspb.Service{},
		CurrentIndex: 0,
	}

	r := mux.NewRouter()
	r.HandleFunc("/lb", lb.ServeHTTP)

	http.Handle("/", r)
	fmt.Println("Load balancer listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
