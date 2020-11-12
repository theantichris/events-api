package objects

import (
	"encoding/json"
	"net/http"
)

// MaxListLimit holds the maximum number of listings.
const MaxListLimit = 200

// GetRequest is for retreiving a single Event.
type GetRequest struct {
	ID string `json:"id"`
}

// ListRequest is for getting a list of Events.
type ListRequest struct {
	Limit int    `json:"limit"`
	After string `json:"after"` // for paging
	Name  string `json:"name"`  // optional name matching
}

// CreateRequest is for creating a new Event.
type CreateRequest struct {
	Event *Event `json:"event"`
}

// UpdateRequest is for updating an existing Event.
type UpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone-number"`
}

// CancelRequest is for canceling an existing Event.
type CancelRequest struct {
	ID string `json:"id"`
}

// RescheduleRequest is for rescheduling an existing Event.
type RescheduleRequest struct {
	ID          string    `json:"id"`
	NewTimeSlot *TimeSlot `json:"new-time-slot"`
}

// DeleteRequest is for deleting an existing Event.
type DeleteRequest struct {
	ID string `json:"id"`
}

// EventResponse holds the response to any event request.
type EventResponse struct {
	Event  *Event   `json:"event,omitempty"`
	Events []*Event `json:"events,omitempty"`
	Code   int      `json:"-"`
}

func (e *EventResponse) Json() []byte {
	panic("implement me")
}

// JSON serializes an EventResponse into JSON.
func (e *EventResponse) JSON() []byte {
	if e == nil {
		return []byte("{}")
	}

	res, _ := json.Marshal(e)

	return res
}

// StatusCode returns the HTTP status code of an EventResponse.
func (e *EventResponse) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}

	return e.Code
}
