package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_loggerWriter_Write(t *testing.T) {
	rec := httptest.NewRecorder()
	w := loggerWriter{ResponseWriter: rec}
	got, err := w.Write(make([]byte, 1))

	wantErr := false
	if (err != nil) != wantErr {
		t.Errorf("loggerWriter.Write() error = %v, wantErr %v", err, wantErr)
		return
	}

	want := 1
	if got != want {
		t.Errorf("loggerWriter.Write() got = %v, want %v", got, want)
	}

	want1 := make([]byte, 1)
	if got1, _ := io.ReadAll(rec.Result().Body); !reflect.DeepEqual(got1, want1) {
		t.Errorf("loggerWriter.Write() got1 = %v, want1 %v", got1, want1)
	}
}

func Test_loggerWriter_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	w := loggerWriter{ResponseWriter: rec}

	want := http.StatusCreated
	w.WriteHeader(want)

	if got := rec.Result().StatusCode; got != want {
		t.Errorf("loggerWriter.WriteHeader() got = %v, want %v", got, want)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := LoggerMiddleware(logger)

	handler := middleware(http.NotFoundHandler())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	var m map[string]any
	if err := json.NewDecoder(&buf).Decode(&m); err != nil {
		t.Fatal(err)
	}

	want := float64(http.StatusNotFound)
	if got := m["code"]; got != want {
		t.Errorf("LoggerMiddleware() got = %v, want %v", got, want)
	}

	want1 := http.MethodGet
	if got1 := m["method"]; got1 != want1 {
		t.Errorf("LoggerMiddleware() got1 = %v, want1 %v", got1, want1)
	}
}

func TestRecovererMiddleware(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	middleware := RecovererMiddleware(logger)

	want := "test"
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic(fmt.Errorf(want)) }))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)

	var m map[string]any
	if err := json.NewDecoder(&buf).Decode(&m); err != nil {
		t.Fatal(err)
	}

	if got := m["err"]; got != want {
		t.Errorf("RecovererMiddleware() got = %v, want %v", got, want)
	}
}
