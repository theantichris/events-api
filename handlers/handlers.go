package handlers

import "net/http"

// EventHandler defines the contract for all the handlers.
type EventHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Cancel(w http.ResponseWriter, r *http.Request)
	Reschedule(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct{}

// NewEventHandler creates and returns a new EventHandler.
func NewEventHandler() EventHandler {
	return &handler{}
}

func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) List(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) Cancel(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h handler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
