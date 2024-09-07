// Пакет service предоставляет реализации бизнес-логики для работы с сущностями.
// Другие пакеты должны принимать интерфейсы бизнес-логики вместо конкретных реализаций.

package service

import (
	"context"
	"dev11/app/entity"
	"errors"
	"fmt"
	"time"
)

// Структура внешней ошибки бизнес-логики.
type ExternalError struct{ Err error }

func (e *ExternalError) Error() string { return fmt.Sprintf("external: %s", e.Err) }

// Структура внутренней ошибки бизнес-логики.
type InternalError struct{ Err error }

func (e *InternalError) Error() string { return fmt.Sprintf("internal: %s", e.Err) }

// Ошибки бизнес-логики.
var (
	ErrInvalidRange error = &ExternalError{errors.New("invalid date range")}
)

// Интерфейс сервиса (бизнес-логики) для сущности "событие".
type Event interface {
	GetForRange(ctx context.Context, userID string, dateStart, dateEnd time.Time) ([]entity.Event, error)
	GetForDay(ctx context.Context, userID string, day time.Time) ([]entity.Event, error)
	GetForWeek(ctx context.Context, userID string, week time.Time) ([]entity.Event, error)
	GetForMonth(ctx context.Context, userID string, month time.Time) ([]entity.Event, error)
	Create(ctx context.Context, event entity.Event) (entity.Event, error)
	Update(ctx context.Context, event entity.Event) (entity.Event, error)
	Delete(ctx context.Context, userID string, id string) error
}
