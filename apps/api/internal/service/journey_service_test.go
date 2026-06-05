package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

type stubJourneyRepository struct {
	journeys       []domain.Journey
	journeyContent domain.JourneyContent
	createdID      int64
	createdJourney *domain.CreateJourney
}

func (r *stubJourneyRepository) List(context.Context) ([]domain.Journey, error) {
	return r.journeys, nil
}

func (r *stubJourneyRepository) Get(context.Context, int64) (domain.JourneyContent, error) {
	return r.journeyContent, nil
}

func (r *stubJourneyRepository) Create(_ context.Context, journey domain.CreateJourney) (int64, error) {
	r.createdJourney = &journey
	return r.createdID, nil
}

type stubJourneyImageRepository struct {
	createdPath string
	createdID   string
	images      map[string]domain.JourneyImage
}

func (r *stubJourneyImageRepository) CreateImage(_ context.Context, path string) (string, error) {
	r.createdPath = path
	return r.createdID, nil
}

func (r *stubJourneyImageRepository) GetImage(_ context.Context, id string) (domain.JourneyImage, error) {
	image, ok := r.images[id]
	if !ok {
		return domain.JourneyImage{}, errors.New("image not found")
	}

	return image, nil
}

type stubImageStorage struct {
	uploadedPath string
	uploadErr    error
	signedURLs   map[string]string
	signedErr    error
}

func (s stubImageStorage) Upload(context.Context, string, string, io.Reader) (StoredImage, error) {
	if s.uploadErr != nil {
		return StoredImage{}, s.uploadErr
	}

	return StoredImage{Key: s.uploadedPath}, nil
}

func (s stubImageStorage) URL(key string) string {
	return key
}

func (s stubImageStorage) SignedURL(key string, expiresIn time.Duration) (string, error) {
	if s.signedErr != nil {
		return "", s.signedErr
	}

	return s.signedURLs[key], nil
}

func TestJourneyServiceList(t *testing.T) {
	timestamp := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC)
	thumbnail := "https://example.com/journey.jpg"
	service := NewJourneyService(&stubJourneyRepository{
		journeys: []domain.Journey{
			{
				Name:      "Build the journey endpoint",
				Timestamp: timestamp,
				Location:  "Bangkok",
				Thumbnail: &thumbnail,
			},
		},
	})

	journeys, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list journeys: %v", err)
	}

	if len(journeys) != 1 {
		t.Fatalf("expected 1 journey, got %d", len(journeys))
	}

	if journeys[0].Name != "Build the journey endpoint" {
		t.Fatalf("expected journey name, got %q", journeys[0].Name)
	}
}

func TestJourneyServiceGet(t *testing.T) {
	timestamp := time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC)
	service := NewJourneyService(&stubJourneyRepository{
		journeyContent: domain.JourneyContent{
			Name:      "Build the journey endpoint",
			Timestamp: timestamp,
			CreatedAt: timestamp,
			Content:   "# Start",
		},
	})

	journey, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get journey: %v", err)
	}

	if journey.Content != "# Start" {
		t.Fatalf("expected journey content, got %q", journey.Content)
	}
}

func TestJourneyServiceCreate(t *testing.T) {
	repository := &stubJourneyRepository{createdID: 42}
	service := NewJourneyService(repository)

	id, err := service.Create(context.Background(), domain.CreateJourney{
		Name:      " Build the journey endpoint ",
		Timestamp: time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC),
		Location:  " Bangkok ",
		Content:   " # Start ",
		Highlight: true,
	})
	if err != nil {
		t.Fatalf("create journey: %v", err)
	}

	if id != 42 {
		t.Fatalf("expected created journey id 42, got %d", id)
	}
	if repository.createdJourney == nil {
		t.Fatal("expected journey row to be created")
	}
	if repository.createdJourney.Name != "Build the journey endpoint" {
		t.Fatalf("expected trimmed journey name, got %q", repository.createdJourney.Name)
	}
}

func TestJourneyServiceCreateRejectsEmptyContent(t *testing.T) {
	validJourney := domain.CreateJourney{
		Name:      "Build the journey endpoint",
		Timestamp: time.Date(2026, time.May, 29, 0, 0, 0, 0, time.UTC),
		Location:  "Bangkok",
		Content:   "# Start",
	}

	tests := map[string]func(*domain.CreateJourney){
		"name": func(journey *domain.CreateJourney) {
			journey.Name = ""
		},
		"timestamp": func(journey *domain.CreateJourney) {
			journey.Timestamp = time.Time{}
		},
		"location": func(journey *domain.CreateJourney) {
			journey.Location = ""
		},
		"content": func(journey *domain.CreateJourney) {
			journey.Content = ""
		},
	}

	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewJourneyService(&stubJourneyRepository{})
			journey := validJourney
			mutate(&journey)

			_, err := service.Create(context.Background(), journey)
			if err != ErrInvalidJourney {
				t.Fatalf("expected ErrInvalidJourney, got %v", err)
			}
		})
	}
}

func TestJourneyServiceUploadImageCreatesImageRow(t *testing.T) {
	images := &stubJourneyImageRepository{createdID: "6dce8aef-4828-4b2c-b15a-78412b05d913"}
	service := NewJourneyService(
		&stubJourneyRepository{},
		WithJourneyImageRepository(images),
		WithJourneyImageStorage(stubImageStorage{uploadedPath: "journey/images/photo.png"}),
	)

	id, err := service.UploadImage(context.Background(), "photo.png", "image/png", strings.NewReader("png-bytes"))
	if err != nil {
		t.Fatalf("upload image: %v", err)
	}

	if id != "6dce8aef-4828-4b2c-b15a-78412b05d913" {
		t.Fatalf("expected image id, got %q", id)
	}
	if images.createdPath != "journey/images/photo.png" {
		t.Fatalf("expected image row path, got %q", images.createdPath)
	}
}

func TestJourneyServiceUploadImageReturnsUploadError(t *testing.T) {
	uploadErr := errors.New("upload failed")
	images := &stubJourneyImageRepository{createdID: "6dce8aef-4828-4b2c-b15a-78412b05d913"}
	service := NewJourneyService(
		&stubJourneyRepository{},
		WithJourneyImageRepository(images),
		WithJourneyImageStorage(stubImageStorage{uploadErr: uploadErr}),
	)

	_, err := service.UploadImage(context.Background(), "photo.png", "image/png", strings.NewReader("png-bytes"))
	if !errors.Is(err, uploadErr) {
		t.Fatalf("expected upload error, got %v", err)
	}
	if images.createdPath != "" {
		t.Fatalf("expected no image row, got path %q", images.createdPath)
	}
}

func TestJourneyServiceGetConvertsImageMarkerToSignedURL(t *testing.T) {
	imageID := "6dce8aef-4828-4b2c-b15a-78412b05d913"
	service := NewJourneyService(
		&stubJourneyRepository{journeyContent: domain.JourneyContent{Content: "before\n\n[image/" + imageID + "]\n\nafter"}},
		WithJourneyImageRepository(&stubJourneyImageRepository{images: map[string]domain.JourneyImage{
			imageID: {ID: imageID, Path: "journey/images/photo.png"},
		}}),
		WithJourneyImageStorage(stubImageStorage{signedURLs: map[string]string{
			"journey/images/photo.png": "https://signed.example.com/photo.png",
		}}),
	)

	journey, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get journey: %v", err)
	}

	if !strings.Contains(journey.Content, "![image](https://signed.example.com/photo.png)") {
		t.Fatalf("expected signed image markdown, got %q", journey.Content)
	}
}

func TestJourneyServiceGetLeavesImageMarkerEmptyWhenSignedURLFails(t *testing.T) {
	imageID := "6dce8aef-4828-4b2c-b15a-78412b05d913"
	service := NewJourneyService(
		&stubJourneyRepository{journeyContent: domain.JourneyContent{Content: "before [image/" + imageID + "] after"}},
		WithJourneyImageRepository(&stubJourneyImageRepository{images: map[string]domain.JourneyImage{
			imageID: {ID: imageID, Path: "journey/images/photo.png"},
		}}),
		WithJourneyImageStorage(stubImageStorage{signedErr: errors.New("sign failed")}),
	)

	journey, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("get journey: %v", err)
	}

	if journey.Content != "before  after" {
		t.Fatalf("expected empty image marker, got %q", journey.Content)
	}
}

func TestJourneyServiceListGetsSignedThumbnailURL(t *testing.T) {
	imageID := "6dce8aef-4828-4b2c-b15a-78412b05d913"
	service := NewJourneyService(
		&stubJourneyRepository{journeys: []domain.Journey{{Thumbnail: &imageID}}},
		WithJourneyImageRepository(&stubJourneyImageRepository{images: map[string]domain.JourneyImage{
			imageID: {ID: imageID, Path: "journey/images/thumbnail.png"},
		}}),
		WithJourneyImageStorage(stubImageStorage{signedURLs: map[string]string{
			"journey/images/thumbnail.png": "https://signed.example.com/thumbnail.png",
		}}),
	)

	journeys, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list journeys: %v", err)
	}

	if journeys[0].Thumbnail == nil || *journeys[0].Thumbnail != "https://signed.example.com/thumbnail.png" {
		t.Fatalf("expected signed thumbnail url, got %v", journeys[0].Thumbnail)
	}
}

func TestJourneyServiceListLeavesThumbnailEmptyWhenSignedURLFails(t *testing.T) {
	imageID := "6dce8aef-4828-4b2c-b15a-78412b05d913"
	service := NewJourneyService(
		&stubJourneyRepository{journeys: []domain.Journey{{Thumbnail: &imageID}}},
		WithJourneyImageRepository(&stubJourneyImageRepository{images: map[string]domain.JourneyImage{
			imageID: {ID: imageID, Path: "journey/images/thumbnail.png"},
		}}),
		WithJourneyImageStorage(stubImageStorage{signedErr: errors.New("sign failed")}),
	)

	journeys, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("list journeys: %v", err)
	}

	if journeys[0].Thumbnail == nil || *journeys[0].Thumbnail != "" {
		t.Fatalf("expected empty thumbnail, got %v", journeys[0].Thumbnail)
	}
}
