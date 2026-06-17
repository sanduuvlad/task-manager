package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	recoveredHandler := Recovery(handler)
	recoveredHandler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Result().StatusCode)
	}
}
