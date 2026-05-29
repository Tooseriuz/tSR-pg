package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/tooseriuz/tsr-pg/apps/api/internal/adapters/http"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/dto/openapi"
)

func TestJourneyHandler(t *testing.T) {
	router := httpadapter.NewRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/journeys", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response openapi.JourneysResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 journey, got %d", len(response))
	}

	if response[0].Name != "Remote desk, first production habits" {
		t.Fatalf("expected journey name, got %q", response[0].Name)
	}
}
