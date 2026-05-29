package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/http"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
)

func TestHealthHandler(t *testing.T) {
	router := httpadapter.NewRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response openapi.HealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}

func TestCorsPreflightAllowsWebsiteOrigin(t *testing.T) {
	router := httpadapter.NewRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/journeys", nil)
	request.Header.Set("Origin", "https://www.tooseriuz.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodGet)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != "https://www.tooseriuz.com" {
		t.Fatalf("expected website origin, got %q", origin)
	}

	if credentials := recorder.Header().Get("Access-Control-Allow-Credentials"); credentials != "true" {
		t.Fatalf("expected credentials true, got %q", credentials)
	}
}
