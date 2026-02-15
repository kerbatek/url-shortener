package service

import (
	"context"

	"github.com/kerbatek/url-shortener/internal/model"
	"github.com/kerbatek/url-shortener/internal/repository"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Shorten(ctx context.Context, originalURL string) (*model.URL, error) {
	// TODO: generate short code and persist via repo
	return nil, nil
}

func (s *URLService) Resolve(ctx context.Context, code string) (*model.URL, error) {
	// TODO: look up original URL by code
	return nil, nil
}

func (s *URLService) Delete(ctx context.Context, id int64) error {
	// TODO: delete URL by ID
	return nil
}
