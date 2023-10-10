package integrationtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wbollock/ping_exporter/internal/server"
)

const expectedStatusCode = 200

func setupTestServer() *httptest.Server {
	handler := server.SetupServer()
	return httptest.NewServer(handler)
}

func validateResponse(t *testing.T, resp *http.Response, expectedBodyContent string) {
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status %d, got: %d", expectedStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if !strings.Contains(string(body), expectedBodyContent) {
		t.Fatalf("Expected %s, got: %v", expectedBodyContent, string(body))
	}
}

func TestPingExporterRootEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "<head><title>Ping Exporter</title></head>")
}

func TestPingExporterMetricsEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "promhttp_metric_handler_requests_in_flight 1")
}

func TestPingExporterProbeEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/probe?target=127.0.0.1")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "ping_success 1")
}
