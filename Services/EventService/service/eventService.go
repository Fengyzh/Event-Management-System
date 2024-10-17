package eventservice

import (
	eventdata "go-lb/Services/EventService/data"
	espb "go-lb/grpc/eventservice"
)

// This is dup with event.go

type EventService struct {
	repo eventdata.EventRepository
}

func NewEventService(repo eventdata.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAllEvents() (*espb.EventList, error) {
	return s.repo.GetAllEvents()
}

func (s *EventService) GetEventById(eid *espb.EventId) (*espb.Event, error) {
	return s.repo.GetEventById(eid)
}

func (s *EventService) CreateEvent(event *espb.EventCreateRequest) error {
	return s.repo.CreateEvent(event)
}

func (s *EventService) UpdateEvent(event *espb.Event) (*espb.EventResponse, error) {
	return s.repo.UpdateEvent(event)
}

func (s *EventService) DeleteEvent(eid *espb.EventId) (*espb.EventResponse, error) {
	return s.repo.DeleteEvent(eid)
}

func (s *EventService) OrderEventTicket(eid *espb.EventId) (*espb.EventResponse, error) {
	return s.repo.OrderTicketEvent(eid)
}