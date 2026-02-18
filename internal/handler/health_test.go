package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubPinger struct{ err error }

func (s *stubPinger) Ping(_ context.Context) error { return s.err }

func setupHealthRouter(db DBPinger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	hh := NewHealthHandler(db)
	r.GET("/health", hh.Liveness)
	r.GET("/ready", hh.Readiness)
	return r
}

func TestLiveness(t *testing.T) {
	r := setupHealthRouter(&stubPinger{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestReadiness_DBUp(t *testing.T) {
	r := setupHealthRouter(&stubPinger{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestReadiness_DBDown(t *testing.T) {
	r := setupHealthRouter(&stubPinger{err: errors.New("connection refused")})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ready", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}
