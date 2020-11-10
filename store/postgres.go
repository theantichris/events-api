package store

import (
	"context"
	"log"
	"os"

	"github.com/theantichris/events-api/errors"
	"github.com/theantichris/events-api/objects"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pg struct {
	db *gorm.DB
}

// NewPostgresEventStore creates and returns a Postgres implementation of an
// EventStore.
func NewPostgresEventStore(conn string) EventStore {
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},
			),
		})

	if err != nil {
		panic("Unable to connect to the database: " + err.Error())
	}

	if err := db.AutoMigrate(&objects.Event{}); err != nil {
		panic("Unable to migrate database: " + err.Error())
	}

	return &pg{db}
}

func (p *pg) Get(ctx context.Context, request *objects.GetRequest) (*objects.Event, error) {
	event := &objects.Event{}

	err := p.db.WithContext(ctx).Take(event, "id = ?", request.ID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.ErrEventNotFound
	}

	return event, err
}

func (p *pg) List(ctx context.Context, request *objects.ListRequest) ([]*objects.Event, error) {
	if request.Limit == 0 || request.Limit > objects.MaxListLimit {
		request.Limit = objects.MaxListLimit
	}

	query := p.db.WithContext(ctx).Limit(request.Limit)

	if request.After != "" {
		query = query.Where("id > ?", request.After)
	}

	if request.Name != "" {
		query = query.Where("name ilike ?", "%"+request.Name+"%")
	}

	list := make([]*objects.Event, 0, request.Limit)

	err := query.Order("id").Find(&list).Error

	return list, err
}

func (p *pg) Create(ctx context.Context, request *objects.CreateRequest) error {
	if request.Event == nil {
		return errors.ErrObjectIsRequired
	}

	request.Event.ID = GenerateUniqueID()
	request.Event.Status = objects.Original
	request.Event.CreatedAt = p.db.NowFunc()

	return p.db.WithContext(ctx).Create(request.Event).Error
}

func (p *pg) Update(ctx context.Context, request *objects.UpdateRequest) error {
	event := &objects.Event{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
		Website:     request.Website,
		Address:     request.Address,
		PhoneNumber: request.PhoneNumber,
		UpdatedAt:   p.db.NowFunc(),
	}

	return p.db.WithContext(ctx).Model(event).Select(
		"name",
		"description",
		"website",
		"address",
		"phone_number",
		"updated_at",
	).Updates(event).Error
}

func (p *pg) Cancel(ctx context.Context, request *objects.CancelRequest) error {
	event := &objects.Event{
		ID:         request.ID,
		Status:     objects.Canceled,
		CanceledAt: p.db.NowFunc(),
	}

	return p.db.WithContext(ctx).Model(event).Select("status", "cancelled_at").Updates(event).Error
}

func (p *pg) Reschedule(ctx context.Context, request *objects.RescheduleRequest) error {
	event := &objects.Event{
		ID:            request.ID,
		TimeSlot:      request.NewTimeSlot,
		Status:        objects.Rescheduled,
		RescheduledAt: p.db.NowFunc(),
	}

	return p.db.WithContext(ctx).Model(event).Select(
		"status",
		"start_time",
		"end_time",
		"rescheduled_at",
	).Updates(event).Error
}

func (p *pg) Delete(ctx context.Context, request *objects.DeleteRequest) error {
	event := &objects.Event{ID: request.ID}

	return p.db.WithContext(ctx).Model(event).Delete(event).Error
}
