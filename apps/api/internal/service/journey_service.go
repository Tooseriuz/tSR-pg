package service

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

var ErrJourneyNotFound = errors.New("journey not found")
var ErrInvalidJourney = errors.New("invalid journey")
var ErrImageStorageNotConfigured = errors.New("image storage not configured")

type JourneyRepository interface {
	List(ctx context.Context) ([]domain.Journey, error)
	Get(ctx context.Context, id int64) (domain.JourneyContent, error)
	Create(ctx context.Context, journey domain.CreateJourney) (int64, error)
}

type JourneyImageRepository interface {
	CreateImage(ctx context.Context, path string) (string, error)
	GetImage(ctx context.Context, id string) (domain.JourneyImage, error)
}

type JourneyService struct {
	repository JourneyRepository
	images     JourneyImageRepository
	storage    ImageStorage
}

type JourneyServiceOption func(*JourneyService)

func WithJourneyImageRepository(repository JourneyImageRepository) JourneyServiceOption {
	return func(service *JourneyService) {
		service.images = repository
	}
}

func WithJourneyImageStorage(storage ImageStorage) JourneyServiceOption {
	return func(service *JourneyService) {
		service.storage = storage
	}
}

func NewJourneyService(repository JourneyRepository, options ...JourneyServiceOption) JourneyService {
	service := JourneyService{repository: repository}
	for _, option := range options {
		option(&service)
	}

	return service
}

func (s JourneyService) List(ctx context.Context) ([]domain.Journey, error) {
	journeys, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	for index := range journeys {
		journeys[index].Thumbnail = s.signedImageURL(ctx, journeys[index].Thumbnail)
	}

	return journeys, nil
}

func (s JourneyService) Get(ctx context.Context, id int64) (domain.JourneyContent, error) {
	journey, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.JourneyContent{}, err
	}

	journey.Content = s.replaceImageMarkers(ctx, journey.Content)

	return journey, nil
}

func (s JourneyService) Create(ctx context.Context, journey domain.CreateJourney) (int64, error) {
	journey.Name = strings.TrimSpace(journey.Name)
	journey.Location = strings.TrimSpace(journey.Location)
	journey.Content = strings.TrimSpace(journey.Content)
	if journey.Thumbnail != nil {
		trimmedThumbnail := strings.TrimSpace(*journey.Thumbnail)
		if trimmedThumbnail == "" {
			journey.Thumbnail = nil
		} else {
			journey.Thumbnail = &trimmedThumbnail
		}
	}
	if journey.Name == "" || journey.Location == "" || journey.Content == "" || journey.Timestamp.IsZero() {
		return 0, ErrInvalidJourney
	}

	return s.repository.Create(ctx, journey)
}

func (s JourneyService) UploadImage(ctx context.Context, fileName string, contentType string, content io.Reader) (string, error) {
	if s.storage == nil || s.images == nil {
		return "", ErrImageStorageNotConfigured
	}

	image, err := s.storage.Upload(ctx, fileName, contentType, content)
	if err != nil {
		return "", err
	}

	return s.images.CreateImage(ctx, image.Key)
}

var imageMarkerPattern = regexp.MustCompile(`\[image/([0-9a-fA-F-]{36})\]`)

func (s JourneyService) replaceImageMarkers(ctx context.Context, content string) string {
	if s.images == nil || s.storage == nil {
		return content
	}

	return imageMarkerPattern.ReplaceAllStringFunc(content, func(marker string) string {
		matches := imageMarkerPattern.FindStringSubmatch(marker)
		if len(matches) != 2 {
			return marker
		}

		url := s.signedImageURLValue(ctx, matches[1])
		if url == "" {
			return ""
		}

		return "![image](" + url + ")"
	})
}

func (s JourneyService) signedImageURL(ctx context.Context, id *string) *string {
	if id == nil {
		return nil
	}

	value := s.signedImageURLValue(ctx, *id)
	return &value
}

func (s JourneyService) signedImageURLValue(ctx context.Context, id string) string {
	if s.images == nil || s.storage == nil {
		return id
	}

	image, err := s.images.GetImage(ctx, id)
	if err != nil {
		return ""
	}

	url, err := s.storage.SignedURL(image.Path, 15*time.Minute)
	if err != nil {
		println("failed to generate signed URL for image", id, ":", err.Error())
		return ""
	}

	return url
}
