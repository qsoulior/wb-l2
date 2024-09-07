// Пакет entity предоставляет структуры сущностей, с которыми работает бизнес-логика,

// а также методы сериализации этих сущностей.

package entity

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEvent_Encode(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2010-05-20")
	e := Event{Title: "event", Date: date}

	var w bytes.Buffer
	wantW := "{\"id\":\"\",\"title\":\"event\",\"description\":\"\",\"date\":\"2010-05-20T00:00:00Z\",\"user_id\":\"\"}\n"

	e.Encode(&w)
	if gotW := w.String(); gotW != wantW {
		t.Errorf("Event.Encode() = %v, want %v", gotW, wantW)
	}
}

func TestEvent_Decode(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2010-05-20")
	r := strings.NewReader(`{"title":"event","date":"2010-05-20T00:00:00Z"}`)

	var e Event
	wantE := Event{Title: "event", Date: date}

	e.Decode(r)
	if gotE := e; !reflect.DeepEqual(gotE, wantE) {
		t.Errorf("Event.Decode() = %v, want %v", gotE, wantE)
	}
}

func TestEvent_ValidateCreate(t *testing.T) {
	id := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	tests := []struct {
		name    string
		e       *Event
		wantErr bool
	}{
		{"ValidEvent", &Event{Title: "event", UserID: id}, false},
		{"InvalidTitle", &Event{UserID: id}, true},
		{"InvalidUserID", &Event{Title: "event"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.ValidateCreate(); (err != nil) != tt.wantErr {
				t.Errorf("Event.ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvent_ValidateUpdate(t *testing.T) {
	id := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	tests := []struct {
		name    string
		e       *Event
		wantErr bool
	}{
		{"ValidEvent", &Event{ID: id, Title: "event", UserID: id}, false},
		{"InvalidID", &Event{Title: "event", UserID: id}, true},
		{"InvalidTitle", &Event{ID: id, UserID: id}, true},
		{"InvalidUserID", &Event{ID: id, Title: "event"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.ValidateUpdate(); (err != nil) != tt.wantErr {
				t.Errorf("Event.ValidateUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
