package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/theantichris/events-api/objects"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/theantichris/events-api/handlers"
	"github.com/theantichris/events-api/store"

	"github.com/gorilla/mux"
)

var (
	router    *mux.Router
	flushAll  func(t *testing.T)
	createOne func(t *testing.T, name string) *objects.Event
	//getOne    func(t *testing.T, id string, wantErr bool) *objects.Event
)

func TestMain(t *testing.M) {
	log.Println("Registering...")

	connection := "postgres://user:password@localhost:5432/db?sslmode=disable"
	if c := os.Getenv("DB"); c != "" {
		connection = c
	}

	router = mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	st := store.NewPostgresEventStore(connection)
	handler := handlers.NewEventHandler(st)

	RegisterAllRoutes(router, handler)

	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(connection), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Event{}, "1=1")
	}

	createOne = func(t *testing.T, name string) *objects.Event {
		event := &objects.Event{
			Name:        name,
			Description: "Description of " + name,
			Website:     "https://" + name + ".com",
			TimeSlot: &objects.TimeSlot{
				Start: time.Now().UTC(),
				End:   time.Now().UTC().Add(time.Hour),
			},
		}

		return event
	}

	//getOne = func(t *testing.T, id string, wantErr bool) *objects.Event {
	//	event, err := st.Get(context.TODO(), objects.GetRequest{ID: id})
	//	if err != nil && wantErr {
	//		t.Fatal(err)
	//	}
	//
	//	return event
	//}

	log.Println("Starting...")

	os.Exit(t.Run())
}

func Do(request *http.Request) *httptest.ResponseRecorder {
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	return writer
}

func TestListEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			name: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 0,
		},
		{
			name: "All",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "One")
				_ = createOne(t, "Two")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "Limited",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "Three")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events?limit=2", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "After",
			setup: func(t *testing.T) *http.Request {
				evt := createOne(t, "Four")
				_ = createOne(t, "Five")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events?after="+evt.ID, nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 1,
		},
		{
			name: "Name",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/events?name=e", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.EventResponse{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, len(got.Events), tt.listLen)
		})
	}
}
