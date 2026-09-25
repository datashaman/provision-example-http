package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	if response.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("verification Cache-Control = %q; want no-store", response.Header.Get("Cache-Control"))
	}
}

func TestSlowEndpointReportsTheRevisionThatAcceptedTheRequest(t *testing.T) {
	server := httptest.NewServer(handler("example-v2"))
	t.Cleanup(server.Close)

	started := time.Now()
	response, err := http.Get(server.URL + "/slow?seconds=1")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed < time.Second {
		t.Fatalf("slow response returned after %s; want at least 1s", elapsed)
	}
	if response.StatusCode != http.StatusOK || string(body) != "{\"revision\":\"example-v2\"}\n" {
		t.Fatalf("slow response = %d %q", response.StatusCode, body)
	}
	if response.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("slow response Cache-Control = %q; want no-store", response.Header.Get("Cache-Control"))
	}
}

func TestSlowEndpointRejectsInvalidDurations(t *testing.T) {
	server := httptest.NewServer(handler("example-v2"))
	t.Cleanup(server.Close)

	for _, query := range []string{"", "0", "31", "nope"} {
		response, err := http.Get(server.URL + "/slow?seconds=" + query)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusBadRequest {
			t.Fatalf("seconds=%q status = %d; want 400", query, response.StatusCode)
		}
	}
}
