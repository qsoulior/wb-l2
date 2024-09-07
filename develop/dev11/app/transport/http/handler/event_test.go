package handler

import (
	"dev11/app/entity"
	"dev11/app/service"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestParseFormEvent(t *testing.T) {
	const layout = "2006-01-02T15:04:05Z"
	date, _ := time.Parse(layout, "2010-05-20T16:00:00Z")

	tests := []struct {
		name    string
		r       func() *http.Request
		want    entity.Event
		wantErr bool
	}{
		{"ValidForm", func() *http.Request {
			data := url.Values{"id": {"0"}, "title": {"0"}, "description": {"0"}, "date": {date.Format(layout)}, "user_id": {"0"}}
			r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			return r
		}, entity.Event{ID: "0", Title: "0", Description: "0", Date: date, UserID: "0"}, false},
		{"InvalidForm", func() *http.Request {
			r := httptest.NewRequest("POST", "/", nil)
			r.Header.Add("Content-Type", "\n")
			return r
		}, entity.EmptyEvent, true},
		{"InvalidDate", func() *http.Request {
			r := httptest.NewRequest("POST", "/", nil)
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			return r
		}, entity.EmptyEvent, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormEvent(tt.r())
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFormEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventCreate_ServeHTTP(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("POST", "/", nil)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entity.EmptyEvent, &service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "title": {"0"}, "description": {"0"}, "date": {"2010-05-20T16:00:00Z"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entity.EmptyEvent, nil)
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "title": {"0"}, "description": {"0"}, "date": {"2010-05-20T16:00:00Z"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventCreate{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventCreate.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventUpdate_ServeHTTP(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("POST", "/", nil)
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.EmptyEvent, &service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "title": {"0"}, "description": {"0"}, "date": {"2010-05-20T16:00:00Z"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.EmptyEvent, nil)
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "title": {"0"}, "description": {"0"}, "date": {"2010-05-20T16:00:00Z"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventUpdate{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventUpdate.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventDelete_ServeHTTP(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("POST", "/", nil)
				r.Header.Add("Content-Type", "\n")
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().Delete(gomock.Any(), gomock.Eq("0"), gomock.Eq("0")).Return(&service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().Delete(gomock.Any(), gomock.Eq("0"), gomock.Eq("0")).Return(nil)
			},
			func() *http.Request {
				data := url.Values{"id": {"0"}, "user_id": {"0"}}
				r := httptest.NewRequest("POST", "/", strings.NewReader(data.Encode()))
				r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return r
			},
			http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventDelete{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventDelete.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetForDay_ServeHTTP(t *testing.T) {
	const layout = "2006-01-02"
	day, _ := time.Parse(layout, "2010-05-20")

	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().GetForDay(gomock.Any(), gomock.Any(), gomock.Eq(day)).Return(nil, &service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"day": {day.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().GetForDay(gomock.Any(), gomock.Any(), gomock.Eq(day)).Return([]entity.Event{}, nil)
			},
			func() *http.Request {
				data := url.Values{"day": {day.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventGetForDay{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventGetForDay.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetForWeek_ServeHTTP(t *testing.T) {
	const layout = "2006-01-02"
	week, _ := time.Parse(layout, "2010-05-20")

	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().GetForWeek(gomock.Any(), gomock.Any(), gomock.Eq(week)).Return(nil, &service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"week": {week.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().GetForWeek(gomock.Any(), gomock.Any(), gomock.Eq(week)).Return([]entity.Event{}, nil)
			},
			func() *http.Request {
				data := url.Values{"week": {week.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventGetForWeek{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventGetForWeek.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetForMonth_ServeHTTP(t *testing.T) {
	const layout = "2006-01"
	month, _ := time.Parse(layout, "2010-05")

	tests := []struct {
		name    string
		prepare func(s *service.MockEvent)
		r       func() *http.Request
		want    int
	}{
		{
			"InvalidForm",
			func(s *service.MockEvent) {},
			func() *http.Request {
				r := httptest.NewRequest("GET", "/", nil)
				return r
			},
			http.StatusBadRequest,
		},
		{
			"ServiceError",
			func(s *service.MockEvent) {
				s.EXPECT().GetForMonth(gomock.Any(), gomock.Any(), gomock.Eq(month)).Return(nil, &service.ExternalError{})
			},
			func() *http.Request {
				data := url.Values{"month": {month.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusServiceUnavailable,
		},
		{
			"ValidForm",
			func(s *service.MockEvent) {
				s.EXPECT().GetForMonth(gomock.Any(), gomock.Any(), gomock.Eq(month)).Return([]entity.Event{}, nil)
			},
			func() *http.Request {
				data := url.Values{"month": {month.Format(layout)}}
				target, _ := url.Parse("/")
				target.RawQuery = data.Encode()
				r := httptest.NewRequest("GET", target.String(), nil)
				return r
			},
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := service.NewMockEvent(ctrl)
			tt.prepare(service)
			h := EventGetForMonth{service}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.r())

			if got := w.Code; got != tt.want {
				t.Errorf("EventGetForMonth.ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
