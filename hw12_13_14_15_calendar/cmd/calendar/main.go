package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/app"
	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/server/http"
	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage/memory"
	databasestorage "github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage/sql"
	"go.uber.org/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

// nolint:funlen,cyclop
func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()

		return
	}

	config := NewConfig(configFile)

	l := logger.New(
		config.Logger.LogFile,
		config.Logger.Level,
		config.Logger.MaxSize,
		config.Logger.MaxBackups,
		config.Logger.MaxAge,
		config.Logger.LocalTime,
		config.Logger.Compress,
	)

	logg := l.InitLogger()
	mode := config.Mode.M
	adr := fmt.Sprintf(":%s", config.Server.Port)

	var s storage.Storage
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch mode {
	case "database":
		s = databasestorage.New(logg)
		err := s.Connect(ctx, config.DataBase.DSN)
		if err != nil {
			logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
		}
		defer func() {
			if err := s.Close(); err != nil {
				logg.Fatal("can't close connection to DB: %s\n", zap.String("err", err.Error()))
			}
		}()
	case "memory":
		s = memorystorage.New(logg)
	default:
		logg.Fatal("can't config mode of storage")
	}

	calendar := app.New(logg, s)
	server := internalhttp.NewServer(calendar, logg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx, adr); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
