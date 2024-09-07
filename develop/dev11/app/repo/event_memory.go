package repo

import (
	"context"
	"dev11/app/entity"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Структура репозитория для сущности "событие", реализующая интерфейс
// и работающая с данными in-memory.
type eventMemory struct {
	mu     sync.RWMutex
	events map[string]entity.Event
}

// NewEventMemory возвращает in-memory репозиторий, реализующий интерфейс.
func NewEventMemory() Event { return &eventMemory{events: make(map[string]entity.Event)} }

// GetByID возвращает Event по его userID и id или ошибку, если Event не найден.
func (e *eventMemory) GetByID(ctx context.Context, userID string, id string) (entity.Event, error) {
	e.mu.RLock()
	event, ok := e.events[id]
	e.mu.RUnlock()
	if !ok || userID != event.UserID {
		return event, ErrNotExist
	}
	return event, nil
}

// GetForRange возвращает []Event по его userID и диапазону дат.
func (e *eventMemory) GetForRange(ctx context.Context, userID string, dateStart, dateEnd time.Time) ([]entity.Event, error) {
	events := make([]entity.Event, 0)
	e.mu.RLock()
	for _, event := range e.events {
		if event.UserID == userID && event.Date.Compare(dateStart) >= 0 && event.Date.Compare(dateEnd) <= 0 {
			events = append(events, event)
		}
	}
	e.mu.RUnlock()
	return events, nil
}

// Create добавляет новый Event в репозиторий, генерируя для него случайный id.
// Возвращает созданный и добавленный Event.
func (e *eventMemory) Create(ctx context.Context, event entity.Event) (entity.Event, error) {
	event.ID = uuid.NewString()
	e.mu.Lock()
	e.events[event.ID] = event
	e.mu.Unlock()
	return event, nil
}

// Update обновляет Event в репозитории.
// Возвращает обновленный Event, если Event существует, иначе возвращает ошибку.
func (e *eventMemory) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	_, err := e.GetByID(ctx, event.UserID, event.ID)
	if err != nil {
		return entity.EmptyEvent, err
	}
	e.mu.Lock()
	e.events[event.ID] = event
	e.mu.Unlock()
	return event, nil
}

// Delete удаляет Event из репозитория, если Event существует, иначе возвращает ошибку.
func (e *eventMemory) Delete(ctx context.Context, userID string, id string) error {
	_, err := e.GetByID(ctx, userID, id)
	if err != nil {
		return err
	}
	e.mu.Lock()
	delete(e.events, id)
	e.mu.Unlock()
	return nil
}
