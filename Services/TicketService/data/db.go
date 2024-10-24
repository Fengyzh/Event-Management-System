package ticketdata

import (
	"database/sql"
	tpb "go-lb/grpc/ticketservice"
	"log"

	_ "github.com/lib/pq"
)

type PGTicketRepository struct {
	db *sql.DB
}

func NewPGTicketRepository(db *sql.DB) *PGTicketRepository {
	return &PGTicketRepository{db: db}
}

func (tr *PGTicketRepository) CreateTicket(ticketCreate *tpb.TicketCreateRequest) (*tpb.TicketOrderResponse, error) {
	_, err := tr.db.Exec("INSERT INTO Tickets (event_id, location, amount, date, seat, fulfilled) VALUES ($1, $2, $3, $4, $5, $6)", ticketCreate.Eventid, ticketCreate.Location, ticketCreate.Amount, ticketCreate.Date, ticketCreate.Seat, 0)

	if err != nil {
		return nil, err
	}

	log.Println("Successfully created event in DB")

	return &tpb.TicketOrderResponse{Message: "Ticket Created"}, err
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


		// Scan the row into the Event struct fields
		err := rows.Scan(&ticket.Ticketid, &ticket.Eventid, &ticket.Amount, &ticket.Location, &ticket.Date, &ticket.Seat, &ticket.Fulfilled)
		if err != nil {
			return nil, err
		}


		tickets.Tickets = append(tickets.Tickets, &ticket)

	}

	return tickets, nil

}

func (tr *PGTicketRepository) GetTicket(tid *tpb.TicketId) (*tpb.TicketOrder, error) {
	row := tr.db.QueryRow("SELECT * FROM Tickets WHERE ticket_id = $1", tid.Ticketid)
	var ticket tpb.TicketOrder
	err := row.Scan(&ticket.Ticketid, &ticket.Eventid, &ticket.Amount, &ticket.Location, &ticket.Date, &ticket.Seat, &ticket.Fulfilled)

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (tr *PGTicketRepository) UpdateTicket(ticket *tpb.TicketOrder) (*tpb.TicketOrderResponse, error) {
	_, err := tr.db.Exec("UPDATE Tickets SET event_id = $1 location = $2, ticket_amount = $3, date = $4, seats = $5, fulfilled = $6 WHERE ticket_id = $7", ticket.Eventid, ticket.Location, ticket.Amount, ticket.Date, ticket.Seat, ticket.Fulfilled, ticket.Ticketid)

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
