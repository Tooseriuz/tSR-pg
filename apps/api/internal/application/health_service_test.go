package application

import "testing"

func TestHealthServiceCheck(t *testing.T) {
	service := NewHealthService()

	status := service.Check()

	if status != "ok" {
		t.Fatalf("expected status ok, got %q", status)
	}
}
