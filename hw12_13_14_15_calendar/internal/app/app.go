package app

import (
	"context"

	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	// logger  logger.Logger
	// storage storage.Storage
}

type Logger interface {
	// TODO: Add more levels
	Info(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type Storage interface {
	// TODO: Add more methods
	AddEvent(ctx context.Context, e *storage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
