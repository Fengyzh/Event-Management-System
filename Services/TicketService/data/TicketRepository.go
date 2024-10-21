package ticketdata

import (
	tpb "go-lb/grpc/ticketservice"

)



type TicketRepository interface {
	GetAllTickets() (*tpb.TicketList, error)
    GetTicket(*tpb.TicketId) (*tpb.TicketOrder, error)
    UpdateTicket(*tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error)
    DeleteTicket(*tpb.TicketId) (*tpb.TicketOrderResponse, error)
}