package event

import (
	"context"
	ticketservice "go-lb/Services/TicketService/service"
	tpb "go-lb/grpc/ticketservice"
	dspb "go-lb/servicediscovery"
	"log"
)

type DiscoveryService struct {
	dspb.UnimplementedDiscoveryServiceServer
}

type TicketService struct {
	tpb.UnimplementedTicketServiceServer
	service     *ticketservice.TicketService
}

func stringPointer(s string) *string {
	return &s
}

func NewTicketService(db *ticketservice.TicketService) *TicketService {

	return &TicketService{service: db}
}

func NewDiscoveryService() *DiscoveryService {

	return &DiscoveryService{}
}

func (ds *DiscoveryService) HealthCheck(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceHealthCheck, error) {

	return &dspb.ServiceHealthCheck{Id: 1, Status: true, Name: stringPointer("First Service"), Addr: "localhost:50052", Message: stringPointer("MSG from service one")}, nil
}

func (ts *TicketService) CreateTicket(ctx context.Context, req *tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error) {


	_, err := ts.service.CreateTicket(req)
	if err != nil {
		log.Fatalf("Fail to Create event in ticket.go, error: %s", err)
	}

	return &tpb.TicketOrderResponse{Message: "Ticket Created"}, nil
}

func (ts *TicketService) GetAllTickets(ctx context.Context, _ *tpb.Empty) (*tpb.TicketList, error) {

	return ts.service.GetAllTickets()
}

func (ts *TicketService) GetTicket(ctx context.Context, tid *tpb.TicketId) (*tpb.TicketOrder, error) {


	return ts.service.GetTicket(tid)

}

func (ts *TicketService) UpdateTicket(ctx context.Context, event *tpb.TicketOrder) (*tpb.TicketOrderResponse, error) {


	return ts.service.UpdateTicket(event)

}

func (ts *TicketService) DeleteTicket(ctx context.Context, tid *tpb.TicketId) (*tpb.TicketOrderResponse, error) {

	return ts.service.DeleteTicket(tid)

}


