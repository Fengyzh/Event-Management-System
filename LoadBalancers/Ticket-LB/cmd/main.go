package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	loadbalancer "go-lb/LoadBalancers/Ticket-LB/internal"
)





func main() {


	lb := loadbalancer.NewLoadBalancer()

	r := mux.NewRouter()
	r.HandleFunc("/ticket", lb.GetTickets).Methods("GET")
	r.HandleFunc("/ticket/{id}", lb.GetTicketById).Methods("GET")
	r.HandleFunc("/ticket", lb.CreateTicket).Methods("POST")
	r.HandleFunc("/ticket/{id}", lb.UpdateTicket).Methods("POST")
	r.HandleFunc("/ticket/{id}", lb.DeleteTicket).Methods("DELETE")
	//r.HandleFunc("/ticket/order/{id}", lb.OrderTicketEvent).Methods("GET")


	http.Handle("/", r)
	fmt.Println("Load balancer listening on port 8081...")
	http.ListenAndServe(":8081", nil)

}
