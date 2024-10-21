package ticketservice

import (
	ticketdata "go-lb/Services/TicketService/data"
	tpb "go-lb/grpc/ticketservice"

)

type TicketService struct {
	repo ticketdata.TicketRepository 
}

func NewTicketService(db ticketdata.TicketRepository) (*TicketService) {
	return &TicketService{repo: db}
}

func (ts *TicketService) GetAllTickets() (*tpb.TicketList, error) {
	return ts.repo.GetAllTickets()
}

func (ts *TicketService) GetTicket(tid *tpb.TicketId) (*tpb.TicketOrder, error) {
	return ts.repo.GetTicket(tid)
}

func (ts *TicketService) UpdateTicket(ticketC *tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error) {
	return ts.repo.UpdateTicket(ticketC)
}

func (ts *TicketService) DeleteTicket(tid *tpb.TicketId) (*tpb.TicketOrderResponse, error) {
	return ts.repo.DeleteTicket(tid)
}