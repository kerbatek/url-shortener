package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerbatek/url-shortener/internal/model"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Skip integration tests if no database is available
		os.Exit(0)
	}

	ctx := context.Background()
	var err error
	testPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	defer testPool.Close()

	// Run migrations
	_, err = testPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
			code        VARCHAR(20)  NOT NULL UNIQUE,
			original_url TEXT        NOT NULL,
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_urls_code ON urls (code);
	`)
	if err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	os.Exit(m.Run())
}

func cleanupURLs(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "DELETE FROM urls")
	if err != nil {
		t.Fatalf("failed to clean urls table: %v", err)
	}
}

func TestCreate(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)
	ctx := context.Background()

	url := &model.URL{
		Code:        "test123",
		OriginalURL: "https://example.com",
	}

	err := repo.Create(ctx, url)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if url.ID == "" {
		t.Error("expected ID to be set")
	}
	if url.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if url.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestCreate_DuplicateCode(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)
	ctx := context.Background()

	url1 := &model.URL{Code: "dup1234", OriginalURL: "https://example.com"}
	url2 := &model.URL{Code: "dup1234", OriginalURL: "https://other.com"}

	if err := repo.Create(ctx, url1); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	err := repo.Create(ctx, url2)
	if err == nil {
		t.Fatal("expected error for duplicate code, got nil")
	}
}

func TestGetByCode_Success(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)
	ctx := context.Background()

	original := &model.URL{Code: "find123", OriginalURL: "https://example.com"}
	if err := repo.Create(ctx, original); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	result, err := repo.GetByCode(ctx, "find123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != original.ID {
		t.Errorf("expected ID %s, got %s", original.ID, result.ID)
	}
	if result.Code != "find123" {
		t.Errorf("expected code find123, got %s", result.Code)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected URL https://example.com, got %s", result.OriginalURL)
	}
}

func TestGetByCode_NotFound(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)

	_, err := repo.GetByCode(context.Background(), "nonexist")
	if err == nil {
		t.Fatal("expected error for missing code, got nil")
	}
}

func TestGetByID_Success(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)
	ctx := context.Background()

	original := &model.URL{Code: "byid123", OriginalURL: "https://example.com"}
	if err := repo.Create(ctx, original); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	result, err := repo.GetByID(ctx, original.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Code != "byid123" {
		t.Errorf("expected code byid123, got %s", result.Code)
	}
	if result.OriginalURL != "https://example.com" {
		t.Errorf("expected URL https://example.com, got %s", result.OriginalURL)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)

	_, err := repo.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
}

func TestDelete_Success(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)
	ctx := context.Background()

	url := &model.URL{Code: "del1234", OriginalURL: "https://example.com"}
	if err := repo.Create(ctx, url); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err := repo.Delete(ctx, url.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it's actually gone
	_, err = repo.GetByID(ctx, url.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestDelete_NotFound(t *testing.T) {
	cleanupURLs(t)
	repo := NewPostgresURLRepository(testPool)

	err := repo.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
}
