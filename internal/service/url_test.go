package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kerbatek/url-shortener/internal/model"
	"github.com/kerbatek/url-shortener/internal/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestShorten_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *model.URL) error {
			u.ID = "550e8400-e29b-41d4-a716-446655440000"
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return nil
		})

	result, err := svc.Shorten(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected original URL https://example.com, got %s", result.OriginalURL)
	}
	if len(result.Code) != codeLength {
		t.Errorf("expected code length %d, got %d", codeLength, len(result.Code))
	}
	if result.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("expected ID 550e8400-e29b-41d4-a716-446655440000, got %s", result.ID)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	_, err := svc.Shorten(context.Background(), "not-a-url")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestShorten_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("db error"))

	_, err := svc.Shorten(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolve_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	expected := &model.URL{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Code:        "abc1234",
		OriginalURL: "https://example.com",
	}

	mockRepo.EXPECT().
		GetByCode(gomock.Any(), "abc1234").
		Return(expected, nil)

	result, err := svc.Resolve(context.Background(), "abc1234")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", result.OriginalURL)
	}
}

func TestResolve_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	mockRepo.EXPECT().
		GetByCode(gomock.Any(), "missing").
		Return(nil, fmt.Errorf("not found"))

	_, err := svc.Resolve(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	mockRepo.EXPECT().
		Delete(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
		Return(nil)

	err := svc.Delete(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockURLRepository(ctrl)
	svc := NewURLService(mockRepo)

	mockRepo.EXPECT().
		Delete(gomock.Any(), "00000000-0000-0000-0000-000000000000").
		Return(fmt.Errorf("url with id 00000000-0000-0000-0000-000000000000 not found"))

	err := svc.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGenerateCode(t *testing.T) {
	code, err := generateCode()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(code) != codeLength {
		t.Errorf("expected length %d, got %d", codeLength, len(code))
	}

	// Verify all characters are in charset
	for _, c := range code {
		found := false
		for _, ch := range charset {
			if c == ch {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("character %c not in charset", c)
		}
	}

	// Verify uniqueness (two codes should differ)
	code2, _ := generateCode()
	if code == code2 {
		t.Error("expected different codes, got identical")
	}
}
