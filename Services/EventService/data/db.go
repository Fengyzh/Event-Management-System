package eventdata

import (
	"bytes"
	"database/sql"
	"encoding/json"
	espb "go-lb/grpc/eventservice"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type PGEventRepository struct {
	db *sql.DB
}

func NewPGEventRepository(db *sql.DB) EventRepository {
	return &PGEventRepository{db: db}
}



func (er *PGEventRepository) CreateEvent(eventcr *espb.EventCreateRequest) error {
	_, err := er.db.Exec("INSERT INTO Events (name, location, ticket_amount, date, seats) VALUES ($1, $2, $3, $4, $5)", eventcr.Name, eventcr.Location, eventcr.Ticketamount, pq.Array(eventcr.Date), pq.Array(eventcr.Seats))

	if err != nil {
		return err
	}

	log.Println("Successfully created event in DB")

	return err
}

func (er *PGEventRepository) GetAllEvents() (*espb.EventList, error) {
	rows, err := er.db.Query("SELECT * FROM Events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events = &espb.EventList{}

	for rows.Next() {
		var event espb.Event
		var seats []string
		var date []string

		// Scan the row into the Event struct fields
		err := rows.Scan(&event.Eventid, &event.Location, &event.Ticketamount, pq.Array(&date), pq.Array(&seats), &event.Name)
		if err != nil {
			return nil, err
		}

		event.Seats = seats
		event.Date = date
		events.Events = append(events.Events, &event)

	}

	return events, nil
}

func (er *PGEventRepository) GetEventById(eid *espb.EventId) (*espb.Event, error) {
	row := er.db.QueryRow("SELECT * FROM Events WHERE event_id = $1", eid.Eventid)
	var event espb.Event
	var seats []string
	var date []string
	err := row.Scan(&event.Eventid, &event.Location, &event.Ticketamount, pq.Array(&date), pq.Array(&seats), &event.Name)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (er *PGEventRepository) UpdateEvent(event *espb.Event) (*espb.EventResponse, error) {
	tx, _ := er.db.Begin()
	var e = &espb.Event{}

	query := "SELECT event_id, ticket_amount FROM Events WHERE event_id = $1 FOR UPDATE"

	err := tx.QueryRow(query, event.Eventid).Scan(&e.Eventid, &e.Ticketamount)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	log.Printf("Locked row with id: %d\n", event.Eventid)

	if e.Ticketamount == 0 {
		tx.Commit()
		return &espb.EventResponse{Eventid: event.Eventid, Message: "Tickets sold out"}, err
	}

	  _, err = tx.Exec("UPDATE Events SET name = $1, location = $2, ticket_amount = $3, date = $4, seats = $5 WHERE event_id = $6", event.Name, event.Location, event.Ticketamount, pq.Array(event.Date), pq.Array(event.Seats), event.Eventid)


	if err != nil {
		tx.Rollback()
		return &espb.EventResponse{Eventid: event.Eventid, Message: "Failed to update event"}, err
	}

	err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }

	log.Println("Successfully updated event in DB")

	return &espb.EventResponse{Eventid: event.Eventid, Message: "Successfully updated event"}, nil

}

func (er *PGEventRepository) DeleteEvent(eid *espb.EventId) (*espb.EventResponse, error) {

	_, err := er.db.Exec("DELETE FROM Events WHERE event_id = $1", eid.Eventid)
	if err != nil {
		return &espb.EventResponse{Eventid: eid.Eventid, Message: "Failed to delete event"}, err
	}

	log.Println("Successfully deleted event in DB")
	return &espb.EventResponse{Eventid: eid.Eventid, Message: "Successfully deleted event in DB"}, nil
}




func (er *PGEventRepository) OrderTicketEvent(eoq *espb.EventOrderRequest) (*espb.EventResponse, error) {

	eid := eoq.Eventid
	tx, _ := er.db.Begin()
	var e = &espb.Event{}

	query := "SELECT event_id, ticket_amount FROM Events WHERE event_id = $1 FOR UPDATE"

	err := tx.QueryRow(query, eid).Scan(&e.Eventid, &e.Ticketamount)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	log.Printf("Locked row with id: %d\n", eid)

	if e.Ticketamount == 0 {
		tx.Commit()
		return &espb.EventResponse{Eventid: eid, Message: "Tickets sold out"}, err
	}

	_, err = tx.Exec("UPDATE Events SET ticket_amount = ticket_amount - 1 WHERE event_id = $1", eid)

	if err != nil {
		tx.Rollback()
		return &espb.EventResponse{Eventid: eid, Message: "Failed to update ticket amount"}, err
	}
	_, err = tx.Exec("UPDATE Events SET seats = array_remove(seats, $1) WHERE event_id = $2")

	if err != nil {
		tx.Rollback()
		return &espb.EventResponse{Eventid: eid, Message: "Failed to update event seat"}, err
	}


	err = er.SendTicketOrder(eoq)
	if err != nil {
		tx.Rollback()
		return &espb.EventResponse{Eventid: eid, Message: "Failed to create ticket"}, err
	}

	err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }

	log.Println("Successfully updated event in DB")

	return &espb.EventResponse{Eventid: eid, Message: "Successfully updated event"}, nil
}


func (er *PGEventRepository) SendTicketOrder(eoq *espb.EventOrderRequest) (error) {
	jsonData, err := json.Marshal(eoq)
    if err != nil {
        log.Printf("Error marshalling ticket: %v\n", err)
        return err
    }

    url := "http://localhost:8081/ticket"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error creating request: %v\n", err)
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request: %v\n", err)
        return err
    }
    defer resp.Body.Close()


    log.Println(resp)

	return nil
}
