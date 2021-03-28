package memorystorage

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/logger"
	"github.com/gadzira/home_work_melnichenko/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	l := logger.New(
		"logfile-for-test",
		"Info",
		1024,
		1,
		1,
		true,
	)
	logg := l.InitLogger()
	s := New(logg)

	t.Run("Add new event", func(t *testing.T) {
		timeOfEvent := "2021-04-02T10:00:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}

		m, _ := time.ParseDuration("45m")
		d := m.Minutes()
		et := st.Add(m)

		shiftTo, _ := time.ParseDuration("10m")
		r := st.Add(time.Duration(-shiftTo))

		ne := &storage.Event{
			ID:          "1",
			Title:       "121",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "How about grow my salary?",
			OwnerID:     "LordOfSummer",
			RemindTime:  r,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.AddEvent(ctx, ne)
		require.Nil(t, err)
	})

	t.Run("Edit existing event", func(t *testing.T) {
		timeOfEvent := "2021-04-02T23:30:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}
		m, _ := time.ParseDuration("1h")
		d := m.Minutes()
		et := st.Add(m)
		shiftTo, _ := time.ParseDuration("11m")
		r := st.Add(time.Duration(-shiftTo))

		ee := &storage.Event{
			ID:          "1",
			Title:       "121",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "Where is my money?",
			OwnerID:     "LordOfSummer",
			RemindTime:  r,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.EditEvent(ctx, ee)
		require.Nil(t, err)

	})

	t.Run("Edit un-existing event", func(t *testing.T) {
		timeOfEvent := "2021-03-31T10:30:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}
		m, _ := time.ParseDuration("15m")
		d := m.Minutes()
		et := st.Add(m)
		shiftTo, _ := time.ParseDuration("15m")
		r := st.Add(time.Duration(-shiftTo))

		ee := &storage.Event{
			ID:          "256",
			Title:       "Happy Programmer's day",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "Let's drink and coding! I know - it's cool!",
			OwnerID:     "JustDev",
			RemindTime:  r,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.EditEvent(ctx, ee)
		require.Error(t, err)
	})

	t.Run("Add another event", func(t *testing.T) {
		timeOfEvent := "2021-04-09T11:30:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}
		m, _ := time.ParseDuration("1h10m")
		d := m.Minutes()
		et := st.Add(m)
		shiftTo, _ := time.ParseDuration("49m")
		r := st.Add(time.Duration(-shiftTo))
		ee := &storage.Event{
			ID:          "2",
			Title:       "review",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "New features of the platform module",
			OwnerID:     "SystemError",
			RemindTime:  r,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.AddEvent(ctx, ee)
		require.Nil(t, err)

	})

	t.Run("Add another event again", func(t *testing.T) {
		timeOfEvent := "2021-05-07T11:30:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}
		m, _ := time.ParseDuration("1h45m")
		d := m.Minutes()
		et := st.Add(m)
		shiftTo, _ := time.ParseDuration("10m")
		r := st.Add(time.Duration(-shiftTo))

		ee := &storage.Event{
			ID:          "3",
			Title:       "All Hands",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "Plans on the year",
			OwnerID:     "ChiefOfAll",
			RemindTime:  r,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.AddEvent(ctx, ee)
		require.Nil(t, err)
	})

	t.Run("Add an event to the same time ", func(t *testing.T) {
		timeOfEvent := "2021-05-07T11:30:00Z"
		st, err := time.Parse(time.RFC3339, timeOfEvent)
		if err != nil {
			log.Fatal("Can't parse time:", timeOfEvent)
		}
		m, _ := time.ParseDuration("1h45m")
		d := m.Minutes()
		et := st.Add(m)
		shiftTo, _ := time.ParseDuration("10m")
		r := st.Add(time.Duration(-shiftTo))

		ee := &storage.Event{
			ID:          "4",
			Title:       "Party hard",
			StartTime:   st,
			EndTime:     et,
			Duration:    d,
			Description: "Absolutely silent party, without music and conversation, but with hard work on Friday till midnight!",
			OwnerID:     "ChiefOfAll",
			RemindTime:  r,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = s.AddEvent(ctx, ee)
		require.Error(t, err)
	})
}
