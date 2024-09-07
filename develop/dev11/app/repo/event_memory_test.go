package repo

import (
	"context"
	"dev11/app/entity"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNewEventMemory(t *testing.T) {
	want := &eventMemory{events: make(map[string]entity.Event)}
	if got := NewEventMemory(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewEventMemory() = %v, want %v", got, want)
	}
}

func Test_eventMemory_GetByID(t *testing.T) {
	t.Run("EventExists", func(t *testing.T) {
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()
		event, _ := e.Create(ctx, entity.Event{UserID: userID})

		want := entity.Event{ID: event.ID, UserID: userID}
		wantErr := false
		got, err := e.GetByID(ctx, userID, event.ID)
		if (err != nil) != wantErr {
			t.Errorf("eventMemory.GetByID() error = %v, wantErr %v", err, wantErr)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("eventMemory.GetByID() = %v, want %v", got, want)
		}
	})

	t.Run("EventDoesNotExist", func(t *testing.T) {
		id := "18310e71-4df6-42c0-adf4-1a280013dd08"
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()

		wantErr := true
		_, err := e.GetByID(ctx, userID, id)
		if (err != nil) != wantErr {
			t.Errorf("eventMemory.GetByID() error = %v, wantErr %v", err, wantErr)
		}
	})
}

func Test_eventMemory_GetForRange(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
	dateStart := (time.Time{}).AddDate(2010, 5, 2)
	dateEnd := (time.Time{}).AddDate(2010, 5, 3)

	e := NewEventMemory()

	want := make([]entity.Event, 3)
	for i := range want {
		event := entity.Event{
			Title:  fmt.Sprintf("event-%d", i),
			Date:   (time.Time{}).AddDate(2010, 5, i+1),
			UserID: userID,
		}
		event, _ = e.Create(ctx, event)
		want[i] = event
	}
	want = want[1:]
	wantErr := false

	got, err := e.GetForRange(ctx, userID, dateStart, dateEnd)
	slices.SortFunc(got, func(a, b entity.Event) int { return strings.Compare(a.Title, b.Title) })
	if (err != nil) != wantErr {
		t.Errorf("eventMemory.GetForRange() error = %v, wantErr %v", err, wantErr)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("eventMemory.GetForRange() = %v, want %v", got, want)
	}
}

func Test_eventMemory_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := NewEventMemory()

	want := "18310e71-4df6-42c0-adf4-1a280013dd08"
	wantErr := false
	got, err := e.Create(ctx, entity.Event{UserID: want})
	if (err != nil) != wantErr {
		t.Errorf("eventMemory.Create() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got := got.UserID; got != want {
		t.Errorf("eventMemory.Create() = %v, want %v", got, want)
	}
}

func Test_eventMemory_Update(t *testing.T) {
	t.Run("EventExists", func(t *testing.T) {
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()
		event, _ := e.Create(ctx, entity.Event{UserID: userID})

		want := entity.Event{ID: event.ID, UserID: userID, Title: "event"}
		wantErr := false
		got, err := e.Update(ctx, want)
		if (err != nil) != wantErr {
			t.Errorf("eventMemory.Update() error = %v, wantErr %v", err, wantErr)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("eventMemory.Update() = %v, want %v", got, want)
		}
	})

	t.Run("EventDoesNotExist", func(t *testing.T) {
		id := "18310e71-4df6-42c0-adf4-1a280013dd08"
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()

		want := entity.Event{ID: id, UserID: userID, Title: "event"}
		wantErr := true
		_, err := e.Update(ctx, want)
		if (err != nil) != wantErr {
			t.Errorf("eventMemory.Update() error = %v, wantErr %v", err, wantErr)
		}
	})
}

func Test_eventMemory_Delete(t *testing.T) {
	t.Run("EventExists", func(t *testing.T) {
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()
		event, _ := e.Create(ctx, entity.Event{UserID: userID})

		wantErr := false
		if err := e.Delete(ctx, userID, event.ID); (err != nil) != wantErr {
			t.Errorf("eventMemory.Delete() error = %v, wantErr %v", err, wantErr)
			return
		}
		if _, err := e.GetByID(ctx, userID, event.ID); err != ErrNotExist {
			t.Errorf("eventMemory.GetByID() error = %v, want %v", err, ErrNotExist)
		}
	})

	t.Run("EventDoesNotExist", func(t *testing.T) {
		id := "18310e71-4df6-42c0-adf4-1a280013dd08"
		userID := "18310e71-4df6-42c0-adf4-1a280013dd08"
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		e := NewEventMemory()

		wantErr := true
		if err := e.Delete(ctx, userID, id); (err != nil) != wantErr {
			t.Errorf("eventMemory.Delete() error = %v, wantErr %v", err, wantErr)
		}
	})
}
