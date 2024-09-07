// Пакет repo предоставляет репозитории для работы с данными сущностей.
// Другие пакеты должны принимать интерфейсы репозиториев вместо конкретных реализаций.

package repo

import (
	"context"
	"dev11/app/entity"
	"errors"
	"time"
)

// Ошибки репозитория.
var (
	ErrExists   = errors.New("already exists")
	ErrNotExist = errors.New("does not exist")
)

// Интерфейс репозитория для сущности "событие".
type Event interface {
	GetByID(ctx context.Context, userID string, id string) (entity.Event, error)
	GetForRange(ctx context.Context, userID string, dateStart, dateEnd time.Time) ([]entity.Event, error)
	Create(ctx context.Context, event entity.Event) (entity.Event, error)
	Update(ctx context.Context, event entity.Event) (entity.Event, error)
	Delete(ctx context.Context, userID string, id string) error
}
