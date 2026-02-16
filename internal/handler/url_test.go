package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kerbatek/url-shortener/internal/model"
	"github.com/kerbatek/url-shortener/internal/repository/mocks"
	"github.com/kerbatek/url-shortener/internal/service"
	"go.uber.org/mock/gomock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(ctrl *gomock.Controller) (*gin.Engine, *mocks.MockURLRepository) {
	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := service.NewURLService(mockRepo)
	h := NewURLHandler(svc)

	router := gin.New()
	router.POST("/shorten", h.ShortenURL)
	router.GET("/:code", h.RedirectURL)
	router.DELETE("/url/:id", h.DeleteURL)

	return router, mockRepo
}

func TestShortenURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockRepo := setupRouter(ctrl)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *model.URL) error {
			u.ID = 1
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		})

	body := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp model.URL
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
	if resp.OriginalURL != "https://example.com" {
		t.Errorf("expected original URL https://example.com, got %s", resp.OriginalURL)
	}
	if resp.Code == "" {
		t.Error("expected non-empty code")
	}
}

func TestShortenURL_MissingURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupRouter(ctrl)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupRouter(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestShortenURL_InvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupRouter(ctrl)

	body := `{"url": "not-a-valid-url"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRedirectURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockRepo := setupRouter(ctrl)

	mockRepo.EXPECT().
		GetByCode(gomock.Any(), "abc1234").
		Return(&model.URL{
			ID:          1,
			Code:        "abc1234",
			OriginalURL: "https://example.com",
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/abc1234", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("expected status 302, got %d", w.Code)
	}
	location := w.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("expected redirect to https://example.com, got %s", location)
	}
}

func TestRedirectURL_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockRepo := setupRouter(ctrl)

	mockRepo.EXPECT().
		GetByCode(gomock.Any(), "missing").
		Return(nil, fmt.Errorf("not found"))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockRepo := setupRouter(ctrl)

	mockRepo.EXPECT().
		Delete(gomock.Any(), int64(1)).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/url/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestDeleteURL_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, mockRepo := setupRouter(ctrl)

	mockRepo.EXPECT().
		Delete(gomock.Any(), int64(999)).
		Return(fmt.Errorf("url with id 999 not found"))

	req := httptest.NewRequest(http.MethodDelete, "/url/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteURL_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	router, _ := setupRouter(ctrl)

	req := httptest.NewRequest(http.MethodDelete, "/url/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
