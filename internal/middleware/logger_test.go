package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger_CallsNextHandler(t *testing.T) {
	// создаем флажок
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	// создаём fake response and request
	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()

	// обернули next в Logger
	loggerHandler := Logger(next)
	loggerHandler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
}

func TestLogger_KeepsStatusCode(t *testing.T) {
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	logger := Logger(next)
	logger.ServeHTTP(rr, req)

	if !called {
		t.Fatal("handler was not called")
	}

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Result().StatusCode)
	}
}
