// Пакет service предоставляет реализации бизнес-логики для работы с сущностями.

// Другие пакеты должны принимать интерфейсы бизнес-логики вместо конкретных реализаций.

package service

import (
	"errors"
	"testing"
)

func TestExternalError_Error(t *testing.T) {
	e := &ExternalError{errors.New("test error")}
	want := "external: test error"
	if got := e.Error(); got != want {
		t.Errorf("ExternalError.Error() = %v, want %v", got, want)
	}
}

func TestInternalError_Error(t *testing.T) {
	e := &InternalError{errors.New("test error")}
	want := "internal: test error"
	if got := e.Error(); got != want {
		t.Errorf("InternalError.Error() = %v, want %v", got, want)
	}
}
