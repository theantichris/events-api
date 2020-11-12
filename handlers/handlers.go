package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/theantichris/events-api/errors"
	"github.com/theantichris/events-api/objects"

	"github.com/theantichris/events-api/store"
)

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

type handler struct {
	store store.EventStore
}

// NewEventHandler creates and returns a new EventHandler.
func NewEventHandler(store store.EventStore) EventHandler {
	return &handler{store}
}

func (h handler) Get(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id == "" {
		WriteError(writer, errors.ErrValidEventIDIsRequired)
		return
	}

	event, err := h.store.Get(request.Context(), objects.GetRequest{ID: id})
	if err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{Event: event})
}

func (h handler) List(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	after := values.Get("after")
	name := values.Get("name")
	limit, err := IntFromString(writer, values.Get("limit"))
	if err != nil {
		return
	}

	events, err := h.store.List(request.Context(), objects.ListRequest{
		Limit: limit,
		After: after,
		Name:  name,
	})
	if err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{Events: events})
}

func (h handler) Create(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		WriteError(writer, errors.ErrUnprocessableEntity)
		return
	}

	event := &objects.Event{}

	if Unmarshal(writer, data, event) != nil {
		WriteError(writer, err)
		return
	}

	if err := checkSlot(event.TimeSlot); err != nil {
		WriteError(writer, err)
		return
	}

	if err = h.store.Create(request.Context(), objects.CreateRequest{Event: event}); err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{Event: event})
}

func (h handler) Update(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		WriteError(writer, errors.ErrUnprocessableEntity)
		return
	}

	updateRequest := &objects.UpdateRequest{}
	if Unmarshal(writer, data, updateRequest) != nil {
		return
	}

	if _, err := h.store.Get(request.Context(), objects.GetRequest{ID: updateRequest.ID}); err != nil {
		WriteError(writer, err)
		return
	}

	if err = h.store.Update(request.Context(), *updateRequest); err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{})
}

func (h handler) Cancel(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		WriteError(writer, errors.ErrValidEventIDIsRequired)
		return
	}

	if _, err := h.store.Get(request.Context(), objects.GetRequest{ID: id}); err != nil {
		WriteError(writer, err)
		return
	}

	if err := h.store.Cancel(request.Context(), objects.CancelRequest{ID: id}); err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{})
}

func (h handler) Reschedule(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		WriteError(writer, errors.ErrUnprocessableEntity)
		return
	}

	rescheduleRequest := &objects.RescheduleRequest{}
	if Unmarshal(writer, data, rescheduleRequest) != nil {
		WriteError(writer, err)
		return
	}

	if err := checkSlot(rescheduleRequest.NewTimeSlot); err != nil {
		WriteError(writer, err)
		return
	}

	if _, err := h.store.Get(request.Context(), objects.GetRequest{ID: rescheduleRequest.ID}); err != nil {
		WriteError(writer, err)
		return
	}

	if err := h.store.Reschedule(request.Context(), *rescheduleRequest); err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{})
}

func (h handler) Delete(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		WriteError(writer, errors.ErrValidEventIDIsRequired)
		return
	}

	if _, err := h.store.Get(request.Context(), objects.GetRequest{ID: id}); err != nil {
		WriteError(writer, err)
		return
	}

	if err := h.store.Delete(request.Context(), objects.DeleteRequest{ID: id}); err != nil {
		WriteError(writer, err)
		return
	}

	WriteResponse(writer, &objects.EventResponse{})
}
