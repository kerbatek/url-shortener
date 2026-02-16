package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerbatek/url-shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	GetByID(ctx context.Context, id string) (*model.URL, error)
	Delete(ctx context.Context, id string) error
}

type postgresURLRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresURLRepository(pool *pgxpool.Pool) URLRepository {
	return &postgresURLRepository{pool: pool}
}

func (r *postgresURLRepository) Create(ctx context.Context, url *model.URL) error {
	return r.pool.QueryRow(ctx,
		"INSERT INTO urls (code, original_url) VALUES ($1, $2) RETURNING id, created_at, updated_at",
		url.Code, url.OriginalURL,
	).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)
}

func (r *postgresURLRepository) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	var url model.URL
	err := r.pool.QueryRow(ctx,
		"SELECT id, code, original_url, created_at, updated_at FROM urls WHERE code = $1",
		code,
	).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *postgresURLRepository) GetByID(ctx context.Context, id string) (*model.URL, error) {
	var url model.URL
	err := r.pool.QueryRow(ctx,
		"SELECT id, code, original_url, created_at, updated_at FROM urls WHERE id = $1",
		id,
	).Scan(&url.ID, &url.Code, &url.OriginalURL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *postgresURLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("url with id %s not found", id)
	}
	return nil
}
