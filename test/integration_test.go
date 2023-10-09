package integrationtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wbollock/ping_exporter/internal/server"
)

func TestPingExporterEndpoint(t *testing.T) {
	handler := server.SetupServer("/metrics")
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/your_endpoint")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got: %d", resp.StatusCode)
	}
}
