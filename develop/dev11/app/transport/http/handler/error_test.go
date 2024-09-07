package handler

import (
	"dev11/app/service"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteValue(t *testing.T) {
	w := httptest.NewRecorder()

	WriteValue(w, http.StatusCreated, map[string]int{"test": 1})

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	want := "{\"test\":1}\n"
	if got := string(body); got != want {
		t.Errorf("WriteValue() got = %v, want %v", got, want)
	}

	want1 := http.StatusCreated
	if got1 := w.Code; got1 != want1 {
		t.Errorf("WriteValue() got1 = %v, want1 %v", got1, want1)
	}
}

func TestWriteResult(t *testing.T) {
	w := httptest.NewRecorder()

	WriteResult(w, http.StatusCreated, 1)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	want := "{\"result\":1}\n"
	if got := string(body); got != want {
		t.Errorf("WriteResult() got = %v, want %v", got, want)
	}

	want1 := http.StatusCreated
	if got1 := w.Code; got1 != want1 {
		t.Errorf("WriteResult() got1 = %v, want1 %v", got1, want1)
	}
}

func TestWriteError(t *testing.T) {
	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()

		WriteError(w, http.StatusBadRequest, fmt.Errorf("test"))

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		want := "{\"error\":\"test\"}\n"
		if got := string(body); got != want {
			t.Errorf("WriteError() got = %v, want %v", got, want)
		}

		want1 := http.StatusBadRequest
		if got1 := w.Code; got1 != want1 {
			t.Errorf("WriteError() got1 = %v, want1 %v", got1, want1)
		}
	})

	t.Run("", func(t *testing.T) {
		w := httptest.NewRecorder()

		WriteError(w, http.StatusBadRequest, nil)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		want := "{\"error\":\"\"}\n"
		if got := string(body); got != want {
			t.Errorf("WriteError() got = %v, want %v", got, want)
		}
	})
}

func TestHandleServiceError(t *testing.T) {
	t.Run("NilError", func(t *testing.T) {
		HandleServiceError(nil, nil)
	})

	t.Run("ExternalError", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := &service.ExternalError{Err: fmt.Errorf("test")}
		HandleServiceError(w, err)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		want := "{\"error\":\"test\"}\n"
		if got := string(body); got != want {
			t.Errorf("WriteError() got = %v, want %v", got, want)
		}

		want1 := http.StatusServiceUnavailable
		if got1 := w.Code; got1 != want1 {
			t.Errorf("WriteError() got1 = %v, want1 %v", got1, want1)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		want := fmt.Errorf("test")
		defer func() {
			r := recover()
			err, _ := r.(*service.InternalError)

			if got := err.Err; got != want {
				t.Errorf("HandleServiceError() got = %v, want %v", got, want)
			}
		}()
		err := &service.InternalError{Err: want}
		HandleServiceError(nil, err)
	})
}
