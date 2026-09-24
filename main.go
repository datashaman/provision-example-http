package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	address := os.Getenv("PROVISION_HTTP_LISTEN")
	if address == "" {
		address = "127.0.0.1:18081"
	}
	revision := os.Getenv("PROVISION_REVISION")
	if revision == "" {
		revision = "unconfigured"
	}
	server := &http.Server{
		Addr:              address,
		Handler:           handler(revision),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Provision HTTP example listening on %s as revision %s", address, revision)
	log.Fatal(server.ListenAndServe())
}

func handler(revision string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /verify", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"revision": revision})
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Provision HTTP example\n"))
	})
	return mux
}
