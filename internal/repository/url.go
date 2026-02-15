package repository

import (
	"context"

	"github.com/kerbatek/url-shortener/internal/model"
)

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	GetByID(ctx context.Context, id int64) (*model.URL, error)
	Delete(ctx context.Context, id int64) error
}

type postgresURLRepository struct {
	// TODO: add *pgxpool.Pool field
}

func NewPostgresURLRepository() URLRepository {
	// TODO: accept *pgxpool.Pool and return initialized repo
	return &postgresURLRepository{}
}

func (r *postgresURLRepository) Create(ctx context.Context, url *model.URL) error {
	// TODO: implement
	return nil
}

func (r *postgresURLRepository) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	// TODO: implement
	return nil, nil
}

func (r *postgresURLRepository) GetByID(ctx context.Context, id int64) (*model.URL, error) {
	// TODO: implement
	return nil, nil
}

func (r *postgresURLRepository) Delete(ctx context.Context, id int64) error {
	// TODO: implement
	return nil
}
