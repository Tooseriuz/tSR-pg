package gcs

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const fakeGCSEndpoint = "http://localhost:4443"
const testBucketName = "tsr-pg-journey-images-test"

func freshStorage(t *testing.T) *ImageStorage {
	t.Helper()

	if !fakeGCSAvailable() {
		t.Skip("fake-gcs-server is not running at localhost:4443")
	}

	ctx := context.Background()
	storage, err := NewImageStorage(ctx, Config{
		BucketName:        testBucketName,
		Endpoint:          fakeGCSEndpoint,
		PublicBaseURL:     fakeGCSEndpoint + "/" + testBucketName,
		SignedURLAccessID: "test@example.com",
		SignedURLHostname: "localhost:4443",
		SignedURLInsecure: true,
		SignedURLSignBytes: func(value []byte) ([]byte, error) {
			return append([]byte("signed:"), value...), nil
		},
	})
	if err != nil {
		t.Fatalf("create gcs storage: %v", err)
	}

	resetBucket(t, ctx, storage.client)

	return storage
}

func fakeGCSAvailable() bool {
	conn, err := net.DialTimeout("tcp", "localhost:4443", 200*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()

	return true
}

func resetBucket(t *testing.T, ctx context.Context, client *storage.Client) {
	t.Helper()

	bucket := client.Bucket(testBucketName)
	_, err := bucket.Attrs(ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		if err := bucket.Create(ctx, "tsr-pg-test", nil); err != nil {
			t.Fatalf("create test bucket: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatalf("get test bucket attrs: %v", err)
	}

	objects := bucket.Objects(ctx, nil)
	for {
		object, err := objects.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			t.Fatalf("list test bucket objects: %v", err)
		}
		if err := bucket.Object(object.Name).Delete(ctx); err != nil {
			t.Fatalf("delete test bucket object: %v", err)
		}
	}
}

func TestImageStorageUploadsFile(t *testing.T) {
	storage := freshStorage(t)

	image, err := storage.Upload(context.Background(), "photo.png", "image/png", strings.NewReader("png-bytes"))
	if err != nil {
		t.Fatalf("upload image: %v", err)
	}

	attrs, err := storage.client.Bucket(testBucketName).Object(image.Key).Attrs(context.Background())
	if err != nil {
		t.Fatalf("expected uploaded file in storage: %v", err)
	}
	if attrs.ContentType != "image/png" {
		t.Fatalf("expected content type image/png, got %q", attrs.ContentType)
	}
}

func TestImageStorageGetsFileURL(t *testing.T) {
	storage := freshStorage(t)

	fileURL := storage.URL("journey/images/photo.png")

	if fileURL != "http://localhost:4443/tsr-pg-journey-images-test/journey/images/photo.png" {
		t.Fatalf("expected file url, got %q", fileURL)
	}
}

func TestImageStorageGetsSignedFileURL(t *testing.T) {
	storage := freshStorage(t)

	signedURL, err := storage.SignedURL("journey/images/photo.png", time.Minute)
	if err != nil {
		t.Fatalf("get signed url: %v", err)
	}

	parsedURL, err := url.Parse(signedURL)
	if err != nil {
		t.Fatalf("parse signed url: %v", err)
	}

	if parsedURL.Scheme != "http" || parsedURL.Host != "localhost:4443" || parsedURL.Path != "/tsr-pg-journey-images-test/journey/images/photo.png" {
		t.Fatalf("expected signed storage url, got %q", signedURL)
	}
	if parsedURL.Query().Get("X-Goog-Expires") == "" {
		t.Fatal("expected signed url expires query")
	}
	if parsedURL.Query().Get("X-Goog-Signature") == "" {
		t.Fatal("expected signed url signature query")
	}
}
