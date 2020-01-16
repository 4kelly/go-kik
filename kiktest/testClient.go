// kiktest contains helpers for mocking the real kik bot api.
package kiktest

import (
	"github.com/r-kells/go-kik/kik"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClient is a helper that creates a mocked Kik API.
// To mock behaviour of the Kik API, simply add routes to the `mux` in your tests.
// See [../kik_test.go](../kikbot_test.go).
func TestClient(t *testing.T) (*kik.Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c, err := kik.NewKikClient(
		server.URL+"/",
		"test",
		"test",
		&http.Client{},
	)
	if err != nil {
		t.Fatalf("error starting the kiktest client: %s", err)
	}

	return c, mux, func() {
		server.Close()
	}
}
