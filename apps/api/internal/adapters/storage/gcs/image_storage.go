package gcs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
	"google.golang.org/api/option"
)

type Config struct {
	BucketName          string
	Endpoint            string
	PublicBaseURL       string
	SignedURLAccessID   string
	SignedURLPrivateKey []byte
	SignedURLHostname   string
	SignedURLInsecure   bool
	SignedURLSignBytes  func([]byte) ([]byte, error)
}

type ImageStorage struct {
	client *storage.Client
	bucket string
	config Config
}

func NewImageStorage(ctx context.Context, config Config) (*ImageStorage, error) {
	if config.BucketName == "" {
		return nil, fmt.Errorf("gcs bucket name is required")
	}

	options := []option.ClientOption{}
	if config.Endpoint != "" {
		if err := os.Setenv("STORAGE_EMULATOR_HOST", config.Endpoint); err != nil {
			return nil, err
		}
		options = append(options, option.WithoutAuthentication())
	}

	client, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}

	if config.PublicBaseURL == "" {
		config.PublicBaseURL = "https://storage.googleapis.com/" + config.BucketName
	}

	return &ImageStorage{
		client: client,
		bucket: config.BucketName,
		config: config,
	}, nil
}

func (s *ImageStorage) Upload(ctx context.Context, fileName string, contentType string, content io.Reader) (service.StoredImage, error) {
	if !strings.HasPrefix(contentType, "image/") {
		return service.StoredImage{}, service.ErrInvalidImage
	}

	key, err := objectKey(fileName, contentType)
	if err != nil {
		return service.StoredImage{}, err
	}

	writer := s.client.Bucket(s.bucket).Object(key).NewWriter(ctx)
	writer.ContentType = contentType
	if _, err := io.Copy(writer, content); err != nil {
		_ = writer.Close()
		return service.StoredImage{}, err
	}
	if err := writer.Close(); err != nil {
		return service.StoredImage{}, err
	}

	return service.StoredImage{Key: key, URL: s.URL(key)}, nil
}

func (s *ImageStorage) URL(key string) string {
	return strings.TrimRight(s.config.PublicBaseURL, "/") + "/" + strings.TrimLeft(key, "/")
}

func (s *ImageStorage) SignedURL(key string, expiresIn time.Duration) (string, error) {
	options := &storage.SignedURLOptions{
		GoogleAccessID: s.config.SignedURLAccessID,
		PrivateKey:     s.config.SignedURLPrivateKey,
		SignBytes:      s.config.SignedURLSignBytes,
		Method:         "GET",
		Expires:        time.Now().Add(expiresIn),
		Scheme:         storage.SigningSchemeV4,
		Hostname:       s.config.SignedURLHostname,
		Insecure:       s.config.SignedURLInsecure,
	}

	return s.client.Bucket(s.bucket).SignedURL(key, options)
}

func objectKey(fileName string, contentType string) (string, error) {
	extension := strings.ToLower(filepath.Ext(fileName))
	if extension == "" {
		extensions, err := mime.ExtensionsByType(contentType)
		if err == nil && len(extensions) > 0 {
			extension = extensions[0]
		}
	}
	if extension == "" {
		extension = ".img"
	}

	id, err := randomID()
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(filepath.Join("journey", "images", id+extension)), nil
}

func randomID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(value[:]), nil
}
