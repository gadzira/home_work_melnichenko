package memorystorage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type Storage struct {
	mu     sync.RWMutex
	event  map[int]storage.Event
	prevID int
	l      *zap.Logger
}

func New(l *zap.Logger) storage.Storage {
	var s Storage
	s.event = make(map[int]storage.Event)
	s.l = l

	return &s
}

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) AddEvent(_ context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ownerID := e.OwnerID
	startTime := e.StartTime

	err := s.isTimeFree(ownerID, startTime)
	if err != nil {
		return fmt.Errorf("time busy: %w", err)
	}

	newid := s.NewID()
	s.event[newid] = storage.Event{
		ID:          strconv.Itoa(newid),
		Title:       e.Title,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Duration:    e.Duration,
		Description: e.Description,
		OwnerID:     e.OwnerID,
		RemindTime:  e.RemindTime,
	}
	s.prevID = newid

	return nil
}

func (s *Storage) EditEvent(_ context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := strconv.Atoi(e.ID)
	if err != nil {
		return fmt.Errorf("incorrect id: %w", err)
	}

	ee, ok := s.event[id]
	if !ok {
		return fmt.Errorf("unknown id, event does't exist")
	}

	ee.ID = e.ID
	ee.Title = e.Title
	ee.StartTime = e.StartTime
	ee.EndTime = e.EndTime
	ee.Duration = e.Duration
	ee.Description = e.Description
	ee.OwnerID = e.OwnerID
	ee.RemindTime = e.RemindTime
	s.event[id] = ee

	return nil
}

func (s *Storage) RemoveEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("incorrect id: %w", err)
	}

	delete(s.event, idInt)

	return nil
}

func (s *Storage) DayListOfEvents(_ context.Context) ([]storage.Event, error) {
	var listOfEvents []storage.Event
	curDay := time.Now().UTC()
	curDayPlusDay := curDay.AddDate(0, 0, 1).UTC()
	for _, i := range s.event {
		if (i.StartTime.After(curDay) || i.StartTime.Equal(curDay)) && (i.StartTime.Before(curDayPlusDay) || i.StartTime.Equal(curDayPlusDay)) {
			listOfEvents = append(listOfEvents, i)
		}
	}

	return listOfEvents, nil
}

func (s *Storage) WeekListOfEvents(_ context.Context) ([]storage.Event, error) {
	var listOfEvents []storage.Event
	curDay := time.Now().UTC()
	curDayPlusWeek := curDay.AddDate(0, 0, 7).UTC()
	for _, i := range s.event {
		if (i.StartTime.After(curDay) || i.StartTime.Equal(curDay)) && (i.StartTime.Before(curDayPlusWeek) || i.StartTime.Equal(curDayPlusWeek)) {
			listOfEvents = append(listOfEvents, i)
		}
	}
	return listOfEvents, nil
}

func (s *Storage) MonthListOfEvents(_ context.Context) ([]storage.Event, error) {
	var listOfEvents []storage.Event
	curDay := time.Now().UTC()
	curDayPlusMonth := curDay.AddDate(0, 1, 0).UTC()
	for _, i := range s.event {
		if (i.StartTime.After(curDay) || i.StartTime.Equal(curDay)) && (i.StartTime.Before(curDayPlusMonth) || i.StartTime.Equal(curDayPlusMonth)) {
			listOfEvents = append(listOfEvents, i)
		}
	}

	return listOfEvents, nil
}

func (s *Storage) NewID() int {
	newID := s.prevID
	newID++
	return newID
}

func (s *Storage) isTimeFree(ownerID string, timeToCheck time.Time) error {
	// 1. Find all events by ownerId
	var allEventsByOwnerID []storage.Event
	for _, i := range s.event {
		if i.OwnerID == ownerID {
			allEventsByOwnerID = append(allEventsByOwnerID, i)
		}
	}

	// 2. Check the new startTime is not between existing events start and end time
	for _, i := range allEventsByOwnerID {
		if timeToCheck.After(i.StartTime) || timeToCheck.Equal(i.StartTime) && timeToCheck.Before(i.EndTime) || timeToCheck.Equal(i.EndTime) {
			return storage.ErrDateBusy
		}
	}
	return nil
}
