package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	loadbalancer "go-lb/LoadBalancer/internal"
)





func main() {


	lb := loadbalancer.NewLoadBalancer()

	r := mux.NewRouter()
	r.HandleFunc("/event", lb.GetEvents).Methods("GET")
	r.HandleFunc("/event/{id}", lb.GetEventById).Methods("GET")
	r.HandleFunc("/event", lb.CreateEvent).Methods("POST")
	r.HandleFunc("/event/{id}", lb.UpdateEvent).Methods("POST")
	r.HandleFunc("/event/{id}", lb.DeleteEvent).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Load balancer listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
