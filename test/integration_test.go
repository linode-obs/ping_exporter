package integrationtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/linode-obs/ping_exporter/internal/server"
)

const expectedStatusCode = 200

func setupTestServer() *httptest.Server {
	handler := server.SetupServer()
	return httptest.NewServer(handler)
}

func validateResponse(t *testing.T, resp *http.Response, expectedBodyContents ...string) {
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status %d, got: %d", expectedStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	for _, content := range expectedBodyContents {
		if !strings.Contains(string(body), content) {
			t.Fatalf("Expected to find %s in response, but not found. Full content: %v", content, string(body))
		}
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

	resp, err := http.Get(server.URL + "/probe?target=127.0.0.1&packet=udp") // UDP so this test can run un-privileged
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "ping_success 1", "ping_timeout 0")
}

func TestPingExporterProbeTimeout(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// this request should always timeout without succeeding
	resp, err := http.Get(server.URL + "/probe?target=localhost&packet=udp&timeout=1s&count=1000")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "ping_success 0", "ping_timeout 1")
}

func TestPingExporterDNSFailure(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/probe?target=invalidhostnamethatdoesntresolve&packet=udp&count=1&timeout=1s")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	validateResponse(t, resp, "ping_success 0")
}

func BenchmarkPingExporterProbeEndpoint(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(server.URL + "/probe?target=127.0.0.1&packet=udp")
		if err != nil {
			b.Fatalf("Failed to send GET request: %v", err)
		}
		resp.Body.Close()
	}
}
