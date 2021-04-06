package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage"

	// import psql driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
	l  *zap.Logger
}

func New(log *zap.Logger) storage.Storage {
	// nolint:exhaustivestruct
	return &Storage{
		l: log,
	}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	var err error
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	s.l.Info("Successfully connected!")

	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(ctx context.Context, e *storage.Event) error {
	err := s.isTimeFree(ctx, e.OwnerID, e.StartTime)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx, `INSERT INTO events (id, title, start_time, end_time, duration, description, owner_id, remind_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		e.ID, e.Title, e.StartTime, e.EndTime, e.Duration, e.Description, e.OwnerID, e.RemindTime)
	if err != nil {
		return fmt.Errorf("cannot insert: %w", err)
	}

	return nil
}

func (s *Storage) EditEvent(ctx context.Context, e *storage.Event) error {
	_, err := s.db.ExecContext(
		ctx, `UPDATE  events SET title = $2, start_time = $3, end_time = $4, duration = $5, description = $6, owner_id = $7, remind_time = $8 WHERE id = $1`,
		e.ID,
		e.Title,
		e.StartTime,
		e.EndTime,
		e.Duration,
		e.Description,
		e.OwnerID,
		e.RemindTime)
	if err != nil {
		return fmt.Errorf("cannot insert: %w", err)
	}

	return nil
}

func (s *Storage) RemoveEvent(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete: %w", err)
	}

	return nil
}

func (s *Storage) DayListOfEvents(ctx context.Context) ([]storage.Event, error) {
	eventListForDay, err := s.getEventsForPeriod(ctx, 0, 0, 1)
	if err != nil {
		return nil, err
	}

	return eventListForDay, nil
}

func (s *Storage) WeekListOfEvents(ctx context.Context) ([]storage.Event, error) {
	eventListForWeek, err := s.getEventsForPeriod(ctx, 0, 0, 8)
	if err != nil {
		return nil, err
	}

	return eventListForWeek, nil
}

func (s *Storage) MonthListOfEvents(ctx context.Context) ([]storage.Event, error) {
	eventListForMonth, err := s.getEventsForPeriod(ctx, 0, 1, 1)
	if err != nil {
		return nil, err
	}

	return eventListForMonth, nil
}

func (s *Storage) getEventsForPeriod(ctx context.Context, y, m, d int) ([]storage.Event, error) {
	// nolint:gofumpt
	var cdf = time.Now().Format("2006-01-02")
	var cd = time.Now()
	var nextDay = cd.AddDate(y, m, d)
	var ndf = nextDay.Format("2006-01-02")

	rows, err := s.db.QueryContext(ctx, `SELECT title, start_time, end_time, duration, description, owner_id, remind_time FROM events WHERE start_time>=$1 AND start_time<$2 ORDER BY start_time`, cdf, ndf)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dbrows err: %w", err)
	}
	defer rows.Close()

	var listOfEvents []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(&e.Title, &e.StartTime, &e.EndTime, &e.Duration, &e.Description, &e.OwnerID, &e.RemindTime); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		listOfEvents = append(listOfEvents, e)
	}

	return listOfEvents, nil
}

func (s *Storage) isTimeFree(ctx context.Context, ownerID string, startTime time.Time) error {
	query := `SELECT title, start_time, end_time, duration, description, owner_id, remind_time FROM events WHERE owner_id = $1 ORDER BY start_time`
	rows, err := s.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return fmt.Errorf("cannot select: %w", err)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("dbrows err: %w", err)
	}
	defer rows.Close()

	var listOfEvents []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(&e.Title, &e.StartTime, &e.EndTime, &e.Duration, &e.Description, &e.OwnerID, &e.RemindTime); err != nil {
			return fmt.Errorf("cannot scan: %w", err)
		}
		listOfEvents = append(listOfEvents, e)
	}

	for _, i := range listOfEvents {
		if inTimeSpan(i.StartTime, i.EndTime, startTime) {
			return storage.ErrDateBusy
		}
	}

	return nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) || check.Equal(start) && check.Before(end) || check.Equal(end)
}
