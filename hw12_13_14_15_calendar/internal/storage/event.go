package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDateBusy = errors.New("date is busy")
)

type Storage interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	EventStore
}

type EventStore interface {
	AddEvent(ctx context.Context, e *Event) error
	EditEvent(ctx context.Context, e *Event) error
	RemoveEvent(ctx context.Context, id string) error
	DayListOfEvents(ctx context.Context) ([]Event, error)
	WeekListOfEvents(ctx context.Context) ([]Event, error)
	MonthListOfEvents(ctx context.Context) ([]Event, error)
}

type Event struct {
	ID          string
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Duration    float64
	Description string
	OwnerID     string
	RemindTime  time.Time
}
