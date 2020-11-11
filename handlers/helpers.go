package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/theantichris/events-api/objects"

	"github.com/theantichris/events-api/errors"
)

type Response interface {
	Json() []byte
	StatusCode() int
}

func WriteResponse(writer http.ResponseWriter, response Response) {
	writer.WriteHeader(response.StatusCode())

	_, _ = writer.Write(response.Json())
}

func WriteError(writer http.ResponseWriter, err error) {
	response, ok := err.(*errors.Error)

	if !ok {
		log.Println(err)
		response = errors.ErrInternal
	}

	WriteResponse(writer, response)
}

func IntFromString(writer http.ResponseWriter, v string) (int, error) {
	if v == "" {
		return 0, nil
	}

	response, err := strconv.Atoi(v)
	if err != nil {
		log.Println(err)
		WriteError(writer, errors.ErrInvalidLimit)
	}

	return response, err
}

func Unmarshal(writer http.ResponseWriter, data []byte, v interface{}) error {
	if d := string(data); d == "nul" || d == "" {
		WriteError(writer, errors.ErrObjectIsRequired)

		return errors.ErrObjectIsRequired
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		log.Println(err)
		WriteError(writer, errors.ErrBadRequest)
	}

	return err
}

func checkSlot(slot *objects.TimeSlot) error {
	if slot == nil {
		return errors.ErrEventTimingIsRequired
	}

	if !slot.Start.After(time.Time{}) {
		return errors.ErrInvalidTimeFormat
	}

	if !slot.End.After(time.Time{}) {
		return errors.ErrInvalidTimeFormat
	}

	return nil
}
