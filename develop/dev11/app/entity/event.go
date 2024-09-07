// Пакет entity предоставляет структуры сущностей, с которыми работает бизнес-логика,
// а также методы сериализации этих сущностей.
package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

// Ошибки валидации сущностей.
var (
	ErrTitleEmpty = errors.New("title is empty")
	ErrIdInvalid  = errors.New("id is invalid")
)

var (
	EmptyEvent = Event{}
)

// Структура сущности "событие".
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	UserID      string    `json:"user_id"`
}

// Encode сериализует Event в json и записывает в w.
func (e Event) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(e)
}

// Decode читает r и десериализует json в Event.
func (e *Event) Decode(r io.Reader) error {
	enc := json.NewDecoder(r)
	return enc.Decode(e)
}

// Validate валидирует поля сущности "событие" при создании.
func (e *Event) ValidateCreate() error {
	if e.Title == "" {
		return fmt.Errorf("title: %w", ErrTitleEmpty)
	}

	if err := uuid.Validate(e.UserID); err != nil {
		return fmt.Errorf("user_id: %w", ErrIdInvalid)
	}

	return nil
}

// Validate валидирует поля сущности "событие" при обновлении.
func (e *Event) ValidateUpdate() error {
	err := uuid.Validate(e.ID)
	if err != nil {
		return fmt.Errorf("id: %w", ErrIdInvalid)
	}

	return e.ValidateCreate()
}
