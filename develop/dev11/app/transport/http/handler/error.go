package handler

import (
	"dev11/app/service"
	"encoding/json"
	"errors"
	"net/http"
)

// WriteValue записывает произвольное значение v в w в формате JSON,
// добавляет http-код code к ответу.
func WriteValue(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.Encode(v)
}

// WriteResult записывает результат res в w вместе с кодом code.
func WriteResult(w http.ResponseWriter, code int, res any) {
	WriteValue(w, code, map[string]any{
		"result": res,
	})
}

// WriteError записывает ошибку err в w вместе с кодом code.
func WriteError(w http.ResponseWriter, code int, err error) {
	var s string
	if err != nil {
		s = err.Error()
	}
	WriteValue(w, code, map[string]string{
		"error": s,
	})
}

// HandleServiceError обрабатывает ошибку сервиса err и записывает в w, если она внешняя,
// иначе вызывает панику для обработки промежуточным слоем.
func HandleServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// Возвращаем пользователю только внешние ошибки бизнес-логики.
	var externalErr *service.ExternalError
	if errors.As(err, &externalErr) {
		WriteError(w, http.StatusServiceUnavailable, externalErr.Err)
		return
	}

	// Паника будет обработана RecovererMiddleware.
	panic(err)
}
