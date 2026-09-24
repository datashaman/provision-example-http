package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthAndRevisionEndpoints(t *testing.T) {
	server := httptest.NewServer(handler("example-v1"))
	t.Cleanup(server.Close)

	for _, path := range []string{"/live", "/ready"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusNoContent {
			t.Fatalf("GET %s status = %d; want 204", path, response.StatusCode)
		}
	}

	response, err := http.Get(server.URL + "/verify")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || string(body) != "{\"revision\":\"example-v1\"}\n" {
		t.Fatalf("verification response = %d %q", response.StatusCode, body)
	}
}
