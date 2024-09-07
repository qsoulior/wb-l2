package service

import (
	"context"
	"dev11/app/entity"
	"dev11/app/repo"
	"errors"
	"time"
)

// Структура сервиса (бизнес-логики) для сущности "событие",
// представляющая первую версию реализации интерфейса.
type eventV1 struct {
	repo repo.Event
}

// NewEventV1 возвращает сервис v1, реализующий интерфейс.
func NewEventV1(repo repo.Event) Event {
	if repo != nil {
		return eventV1{repo}
	}
	return nil
}

// GetForRange возвращает []Event по его userID и диапазону дат, валидируя входные данные.
func (e eventV1) GetForRange(ctx context.Context, userID string, dateStart, dateEnd time.Time) ([]entity.Event, error) {
	if dateEnd.Before(dateStart) {
		return nil, ErrInvalidRange
	}

	events, err := e.repo.GetForRange(ctx, userID, dateStart, dateEnd)
	if err != nil {
		return nil, &InternalError{err}
	}

	return events, nil
}

// GetForDay возвращает []Event по его userID и дню day, валидируя входные данные.
func (e eventV1) GetForDay(ctx context.Context, userID string, day time.Time) ([]entity.Event, error) {
	d := 24 * time.Hour
	dateStart := day.Truncate(d)
	dateEnd := dateStart.Add(d - time.Nanosecond)

	events, err := e.repo.GetForRange(ctx, userID, dateStart, dateEnd)
	if err != nil {
		return nil, &InternalError{err}
	}

	return events, nil
}

// GetForWeek возвращает []Event по его userID и неделе week, валидируя входные данные.
func (e eventV1) GetForWeek(ctx context.Context, userID string, week time.Time) ([]entity.Event, error) {
	d := 7 * 24 * time.Hour
	dateStart := week.Truncate(d)
	dateEnd := dateStart.Add(d - time.Nanosecond)

	events, err := e.repo.GetForRange(ctx, userID, dateStart, dateEnd)
	if err != nil {
		return nil, &InternalError{err}
	}

	return events, nil
}

// GetForMonth возвращает []Event по его userID и месяцу month, валидируя входные данные.
func (e eventV1) GetForMonth(ctx context.Context, userID string, month time.Time) ([]entity.Event, error) {
	dateStart := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	dateEnd := dateStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	events, err := e.repo.GetForRange(ctx, userID, dateStart, dateEnd)
	if err != nil {
		return nil, &InternalError{err}
	}

	return events, nil
}

// Create валидирует входные данные, создает новый Event и возвращает его.
func (e eventV1) Create(ctx context.Context, event entity.Event) (entity.Event, error) {
	if err := event.ValidateCreate(); err != nil {
		return entity.EmptyEvent, &ExternalError{err}
	}

	event, err := e.repo.Create(ctx, event)
	if err != nil {
		return entity.EmptyEvent, &InternalError{err}
	}

	return event, nil
}

// Update валидирует входные данные, обновляет существующий Event и возвращает его.
func (e eventV1) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	if err := event.ValidateUpdate(); err != nil {
		return entity.EmptyEvent, &ExternalError{err}
	}

	event, err := e.repo.Update(ctx, event)
	if err != nil {
		if errors.Is(err, repo.ErrNotExist) {
			return entity.EmptyEvent, &ExternalError{err}
		}
		return entity.EmptyEvent, &InternalError{err}
	}

	return event, nil
}

// Delete удаляет существующий Event по его userID и id.
func (e eventV1) Delete(ctx context.Context, userID string, id string) error {
	err := e.repo.Delete(ctx, userID, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotExist) {
			return &ExternalError{err}
		}
		return &InternalError{err}
	}

	return nil
}
