package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/theantichris/events-api/objects"
)

// EventStore defines the database interactions for storing Events.
type EventStore interface {
	Get(ctx context.Context, request *objects.GetRequest) (*objects.Event, error)
	List(ctx context.Context, request *objects.ListRequest) ([]*objects.Event, error)
	Create(ctx context.Context, request *objects.CreateRequest) error
	Update(ctx context.Context, request *objects.UpdateRequest) error
	Cancel(ctx context.Context, request *objects.CancelRequest) error
	Reschedule(ctx context.Context, request *objects.RescheduleRequest) error
	Delete(ctx context.Context, request *objects.DeleteRequest) error
}

func init() {
	rand.Seed(time.Now().UTC().Unix())
}

// GenerateUniqueID creates a time based sortable unique ID.
func GenerateUniqueID() string {
	word := []byte("0987654321")

	rand.Shuffle(len(word), func(i, j int) {
		word[i], word[j] = word[j], word[i]
	})

	now := time.Now().UTC()

	return fmt.Sprintf("%010v-%010v-%s", now.Unix(), now.Nanosecond(), string(word))
}
