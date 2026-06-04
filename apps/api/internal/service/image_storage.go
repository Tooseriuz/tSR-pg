package service

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrInvalidImage = errors.New("invalid image")

type StoredImage struct {
	Key string
	URL string
}

type ImageStorage interface {
	Upload(ctx context.Context, fileName string, contentType string, content io.Reader) (StoredImage, error)
	URL(key string) string
	SignedURL(key string, expiresIn time.Duration) (string, error)
}
