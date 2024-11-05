# Simple Event Management System

This is a microservice based simple event management system written in Go
    - This is mainly a learning project, for simplicity all services are under the same module

## Services

- Discovery Service: Allow clients to register themselves and perform periodic health checks on the registered services
- Load Balancers (One for each service): Load balance using the round-robin strategy
- Event Service: Handles the event details such as location, total ticket amount and dates
- Ticket Service: Handles the ticket information for the events
- Database (One for each service): Postgres DB containing information related to the service