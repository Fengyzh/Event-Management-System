package ticketdata

import (
	"database/sql"
	tpb "go-lb/grpc/ticketservice"

	"github.com/lib/pq"
)

type PGTicketRepository struct {
	db *sql.DB
}

func NewPGTicketRepository(db *sql.DB) (*PGTicketRepository) {
	return &PGTicketRepository{db: db}
}

func (tr *PGTicketRepository) GetAllTickets() (*tpb.TicketList, error) {
	rows, err := tr.db.Query("SELECT * FROM Tickets")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tickets = &tpb.TicketList{}
	//var tempEvents []*espb.Event

	for rows.Next() {
		var ticket tpb.TicketOrder
		var seats string
		var date []string

		// Scan the row into the Event struct fields
		err := rows.Scan(&ticket.Eventid, &ticket.Ticketid ,&ticket.Location, &ticket.Amount, pq.Array(&date), pq.Array(&seats), &ticket.Fulfilled)
		if err != nil {
			return nil, err
		}

		ticket.Seat = seats
		ticket.Date = date
		tickets.Tickets = append(tickets.Tickets, &ticket)

	}

	return tickets, nil

}

func (tr *PGTicketRepository) GetTicket(tid *tpb.TicketId) (*tpb.TicketOrder, error) {
	row := tr.db.QueryRow("SELECT * FROM Tickets WHERE ticket_id = $1", tid.Ticketid)
	var ticket tpb.TicketOrder
	var seats string
	var date []string
	err := row.Scan(&ticket.Eventid, &ticket.Ticketid ,&ticket.Location, &ticket.Amount, pq.Array(&date), pq.Array(&seats), &ticket.Fulfilled)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (tr *PGTicketRepository) UpdateTicket(ticket *tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error) {
	_, err := tr.db.Exec("UPDATE Tickets SET Event_id = $1 location = $2, ticket_amount = $3, date = $4, seats = $5, fulfilled = $6 WHERE event_id = $7", ticket.Eventid, ticket.Location, ticket.Amount, pq.Array(ticket.Date), pq.Array(ticket.Seat), ticket.Eventid)

	if err != nil {
		return nil, err
	}

	return &tpb.TicketOrderResponse{Message: "Successfully Updated Ticket"}, err
}

func (tr *PGTicketRepository) DeleteTicket(tid *tpb.TicketId) (*tpb.TicketOrderResponse, error) {
	_, err := tr.db.Exec("DELETE FROM Tickets WHERE ticket_id = $1", tid.Ticketid)

	if err != nil {
		return nil, err
	}

	return &tpb.TicketOrderResponse{Message: "Successfully Deleted Ticket"}, nil

}




