package integration

import (
	"bytes"
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

func TestAdminVerifyHandler(t *testing.T) {
	router := httpadapter.NewRouter(httpadapter.WithAdminToken("secret-token"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin-verify", bytes.NewBufferString(`{"token":"secret-token"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if len(recorder.Result().Cookies()) == 0 {
		t.Fatal("expected admin session cookie")
	}
}

func TestAdminVerifyHandlerRejectsInvalidToken(t *testing.T) {
	router := httpadapter.NewRouter(httpadapter.WithAdminToken("secret-token"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/admin-verify", bytes.NewBufferString(`{"token":"wrong-token"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestGetAdminVerifyHandlerAcceptsSessionCookie(t *testing.T) {
	router := httpadapter.NewRouter(httpadapter.WithAdminToken("secret-token"))
	verifyRecorder := httptest.NewRecorder()
	verifyRequest := httptest.NewRequest(http.MethodPost, "/admin-verify", bytes.NewBufferString(`{"token":"secret-token"}`))
	verifyRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(verifyRecorder, verifyRequest)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/admin-verify", nil)
	for _, cookie := range verifyRecorder.Result().Cookies() {
		request.AddCookie(cookie)
	}

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestCreateJourneyHandler(t *testing.T) {
	router := httpadapter.NewRouter(httpadapter.WithAdminToken("secret-token"))
	verifyRecorder := httptest.NewRecorder()
	verifyRequest := httptest.NewRequest(http.MethodPost, "/admin-verify", bytes.NewBufferString(`{"token":"secret-token"}`))
	verifyRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(verifyRecorder, verifyRequest)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/journey", bytes.NewBufferString(`{"name":"Build the journey endpoint","timestamp":"2026-05-29","location":"Bangkok","content":"# Start","highlight":true}`))
	request.Header.Set("Content-Type", "application/json")
	for _, cookie := range verifyRecorder.Result().Cookies() {
		request.AddCookie(cookie)
	}

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response openapi.CreateJourneyResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Id < 1 {
		t.Fatalf("expected created journey id, got %d", response.Id)
	}
}

func TestCreateJourneyHandlerRequiresAdminSession(t *testing.T) {
	router := httpadapter.NewRouter(httpadapter.WithAdminToken("secret-token"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/journey", bytes.NewBufferString(`{"name":"Build the journey endpoint","timestamp":"2026-05-29","location":"Bangkok","content":"# Start","highlight":true}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestGetJourneyHandler(t *testing.T) {
	router := httpadapter.NewRouter()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/journey/1", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response openapi.JourneyContent
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Name != "Remote desk, first production habits" {
		t.Fatalf("expected journey name, got %q", response.Name)
	}

	if response.Content == "" {
		t.Fatal("expected journey content")
	}

	if response.CreatedAt == "" {
		t.Fatal("expected journey posted date")
	}
}
