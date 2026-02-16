package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"

	"github.com/kerbatek/url-shortener/internal/model"
	"github.com/kerbatek/url-shortener/internal/repository"
)

const (
	codeLength = 7
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func generateCode() (string, error) {
	b := make([]byte, codeLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func (s *URLService) Shorten(ctx context.Context, originalURL string) (*model.URL, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	u := &model.URL{
		Code:        code,
		OriginalURL: originalURL,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *URLService) Resolve(ctx context.Context, code string) (*model.URL, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *URLService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
