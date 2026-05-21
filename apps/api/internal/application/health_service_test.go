package application

import "testing"

func TestHealthServiceCheck(t *testing.T) {
	service := NewHealthService()

	response := service.Check()

	if response.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Status)
	}
}
