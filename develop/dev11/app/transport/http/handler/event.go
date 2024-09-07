package handler

import (
	"dev11/app/entity"
	"dev11/app/service"
	"net/http"
	"time"
)

// ParseFormEvent парсит Event, переданный в виде www-url-form-encoded,
// возвращает ошибку, если данные нелья распарсить.
func ParseFormEvent(r *http.Request) (entity.Event, error) {
	if err := r.ParseForm(); err != nil {
		return entity.EmptyEvent, err
	}

	dateValue := r.FormValue("date")
	date, err := time.Parse("2006-01-02T15:04:05Z", dateValue)
	if err != nil {
		return entity.EmptyEvent, err
	}

	return entity.Event{
		ID:          r.FormValue("id"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Date:        date,
		UserID:      r.FormValue("user_id"),
	}, nil
}

// Структура HTTP-обработчика для метода /create_event.
type EventCreate struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := ParseFormEvent(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	event, err = h.Service.Create(r.Context(), event)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusCreated, event)
}

// Структура HTTP-обработчика для метода /update_event.
type EventUpdate struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := ParseFormEvent(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	event, err = h.Service.Update(r.Context(), event)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusOK, event)
}

// Структура HTTP-обработчика для метода /delete_event.
type EventDelete struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventDelete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.Service.Delete(r.Context(), r.FormValue("user_id"), r.FormValue("id"))
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusNoContent, nil)
}

// Структура HTTP-обработчика для метода /events_for_day.
type EventGetForDay struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventGetForDay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	dayValue := query.Get("day")
	day, err := time.Parse("2006-01-02", dayValue)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	events, err := h.Service.GetForDay(r.Context(), query.Get("user_id"), day)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusOK, events)
}

// Структура HTTP-обработчика для метода /events_for_week.
type EventGetForWeek struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventGetForWeek) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	weekValue := query.Get("week")
	week, err := time.Parse("2006-01-02", weekValue)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	events, err := h.Service.GetForWeek(r.Context(), query.Get("user_id"), week)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusOK, events)
}

// Структура HTTP-обработчика для метода /events_for_month.
type EventGetForMonth struct {
	Service service.Event
}

// ServeHTTP обрабатывает запрос r и записывает ответ в w.
func (h EventGetForMonth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	monthValue := query.Get("month")
	month, err := time.Parse("2006-01", monthValue)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	events, err := h.Service.GetForMonth(r.Context(), query.Get("user_id"), month)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	WriteResult(w, http.StatusOK, events)
}
