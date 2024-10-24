package ticketdata

import (
	tpb "go-lb/grpc/ticketservice"

)



type TicketRepository interface {
    CreateTicket(*tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error)
	GetAllTickets() (*tpb.TicketList, error)
    GetTicket(*tpb.TicketId) (*tpb.TicketOrder, error)
    UpdateTicket(*tpb.TicketOrder) (*tpb.TicketOrderResponse, error)
    DeleteTicket(*tpb.TicketId) (*tpb.TicketOrderResponse, error)
}