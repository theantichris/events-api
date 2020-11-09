package objects

import "time"

// EventStatus holds the status of the event.
type EventStatus string

// Default event statuses.
const (
	Original    EventStatus = "original"
	Canceled    EventStatus = "canceled"
	Rescheduled EventStatus = "rescheduled"
)

// TimeSlot holds the start and end times for the event.
type TimeSlot struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

// Event object for the API.
type Event struct {
	ID string `gorm:"primary_key" json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone-number,omitempty"`

	Slot *TimeSlot `gorm:"embedded" json:"slot,omitempty"`

	Status EventStatus `json:"status,omitempty"`

	CreatedAt     time.Time `json:"created-at,omitempty"`
	UpdatedAt     time.Time `json:"updated-at,omitempty"`
	CanceledAt    time.Time `json:"canceled-at,omitempty"`
	RescheduledAt time.Time `json:"reschedueled-at,omitempty"`
}
