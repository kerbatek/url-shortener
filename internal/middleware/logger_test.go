package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouterWithLogger(buf *bytes.Buffer) *gin.Engine {
	logger := zerolog.New(buf).With().Timestamp().Logger()
	router := gin.New()
	router.Use(Logger(logger))
	return router
}

func TestLogger_LogsRequestFields(t *testing.T) {
	var buf bytes.Buffer
	router := setupRouterWithLogger(&buf)
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if logEntry["method"] != "GET" {
		t.Errorf("expected method GET, got %v", logEntry["method"])
	}
	if logEntry["path"] != "/test" {
		t.Errorf("expected path /test, got %v", logEntry["path"])
	}
	if logEntry["status"] != float64(200) {
		t.Errorf("expected status 200, got %v", logEntry["status"])
	}
	if logEntry["user_agent"] != "test-agent" {
		t.Errorf("expected user_agent test-agent, got %v", logEntry["user_agent"])
	}
	if logEntry["message"] != "request" {
		t.Errorf("expected message request, got %v", logEntry["message"])
	}
	if _, ok := logEntry["ip"]; !ok {
		t.Error("expected ip field to be present")
	}
	if _, ok := logEntry["latency"]; !ok {
		t.Error("expected latency field to be present")
	}
}

func TestLogger_LogsCorrectStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		method     string
		path       string
		wantStatus float64
	}{
		{"not found", http.StatusNotFound, http.MethodGet, "/missing", 404},
		{"created", http.StatusCreated, http.MethodPost, "/create", 201},
		{"no content", http.StatusNoContent, http.MethodDelete, "/remove", 204},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			router := setupRouterWithLogger(&buf)
			router.Handle(tt.method, tt.path, func(c *gin.Context) {
				c.Status(tt.status)
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var logEntry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("failed to parse log entry: %v", err)
			}

			if logEntry["status"] != tt.wantStatus {
				t.Errorf("expected status %v, got %v", tt.wantStatus, logEntry["status"])
			}
			if logEntry["method"] != tt.method {
				t.Errorf("expected method %s, got %v", tt.method, logEntry["method"])
			}
			if logEntry["path"] != tt.path {
				t.Errorf("expected path %s, got %v", tt.path, logEntry["path"])
			}
		})
	}
}
