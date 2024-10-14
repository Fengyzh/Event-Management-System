package eventdata

import (
	espb "go-lb/grpc/eventservice"

)




type EventRepository interface {
	GetAllEvents() (*espb.EventList, error)
	GetEventById(*espb.EventId) (*espb.Event, error)
	CreateEvent(*espb.EventCreateRequest) (error)
	UpdateEvent(*espb.Event) (*espb.EventResponse, error)
	DeleteEvent(*espb.EventId) (*espb.EventResponse, error)
}

