package event

import (
	"context"
	eventservice "go-lb/Services/EventService/service"
	espb "go-lb/grpc/eventservice"
	dspb "go-lb/servicediscovery"
	"log"
	"math/rand"
)

type DiscoveryService struct {
	dspb.UnimplementedDiscoveryServiceServer
}

type EventService struct {
	espb.UnimplementedEventServiceServer
	Events []*espb.Event
	db     *eventservice.EventService
}

func stringPointer(s string) *string {
	return &s
}

func NewEventService(db *eventservice.EventService) *EventService {

	return &EventService{Events: []*espb.Event{}, db: db}
}

func NewDiscoveryService() *DiscoveryService {

	return &DiscoveryService{}
}

func (ds *DiscoveryService) HealthCheck(ctx context.Context, _ *dspb.Empty) (*dspb.ServiceHealthCheck, error) {

	return &dspb.ServiceHealthCheck{Id: 1, Status: true, Name: stringPointer("First Service"), Addr: "localhost:50052", Message: stringPointer("MSG from service one")}, nil
}

func (es *EventService) CreateEvent(ctx context.Context, req *espb.EventCreateRequest) (*espb.EventResponse, error) {

	var eid int32 = rand.Int31()
	newEvent := &espb.Event{Eventid: eid, Ticketamount: req.Ticketamount, Name: req.Name, Location: req.Location, Date: req.Date, Seats: req.Seats}
	es.Events = append(es.Events, newEvent)

	err := es.db.CreateEvent(req)
	if err != nil {
		log.Fatalf("Fail to Create event in event.go, error: %s", err)
	}

	return &espb.EventResponse{Eventid: eid, Message: "Event Created"}, nil
}

func (es *EventService) GetAllEvents(ctx context.Context, _ *espb.Empty) (*espb.EventList, error) {

	return es.db.GetAllEvents()
	//return &espb.EventList{Events: es.Events}, nil
}

func (es *EventService) GetEvent(ctx context.Context, eid *espb.EventId) (*espb.Event, error) {


	return es.db.GetEventById(eid)


	/* for _, e := range es.Events {
		if eid.Eventid == e.Eventid {
			return e, nil
		}
	} */
}

func (es *EventService) UpdateEvent(ctx context.Context, event *espb.Event) (*espb.EventResponse, error) {

	/* eid := event.Eventid
	for idx, e := range es.Events {
		if e.Eventid == eid {
			es.Events[idx] = event
			break
		}
	} */

	return es.db.UpdateEvent(event)

	//return &espb.EventResponse{Eventid: event.Eventid, Message: "Update Complete"}, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, eid *espb.EventId) (*espb.EventResponse, error) {

	return es.db.DeleteEvent(eid)

}


func (es *EventService) OrderEventTicket(ctx context.Context, eid *espb.EventId) (*espb.EventResponse, error) {
	
	return es.db.OrderEventTicket(eid)
}

