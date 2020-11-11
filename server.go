package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theantichris/events-api/handlers"
)

// Args holds the arguments used to run the server.
type Args struct {
	// Postgres connection string e.g. "postgres://user:password@localhost:5432/database?sslmode=disable"
	conn string

	// Port for the server e.g. ":8080
	port string
}

// Run runs the server based on the given args.
func Run(args Args) error {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()

	//st := store.NewPostgresEventStore(args.conn)
	handler := handlers.NewEventHandler( /*st*/ )

	RegisterAllRoutes(router, handler)

	log.Println("Starting server at port:", args.port)

	return http.ListenAndServe(args.port, router)
}

func RegisterAllRoutes(router *mux.Router, handler handlers.EventHandler) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(writer, request)
		})
	})

	router.HandleFunc("/event", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/event", handler.Create).Methods(http.MethodPost)
	router.HandleFunc("/event", handler.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/event/cancel", handler.Cancel).Methods(http.MethodPatch)
	router.HandleFunc("/event/details", handler.Update).Methods(http.MethodPut)
	router.HandleFunc("/event/reschedule", handler.Reschedule).Methods(http.MethodPatch)
	router.HandleFunc("/events", handler.List).Methods(http.MethodGet)
}
