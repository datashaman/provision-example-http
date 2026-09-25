package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		Handler:           handler(revision, address),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Provision HTTP example listening on %s as revision %s", address, revision)
	log.Fatal(server.ListenAndServe())
}

func handler(revision, directHost string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /live", func(w http.ResponseWriter, r *http.Request) {
		if rejectStableHealth(w, r, revision, directHost) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		if rejectStableHealth(w, r, revision, directHost) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /verify", func(w http.ResponseWriter, r *http.Request) {
		if rejectStableHealth(w, r, revision, directHost) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"revision": revision})
	})
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		delay, err := strconv.Atoi(r.URL.Query().Get("seconds"))
		if err != nil || delay < 1 || delay > 30 {
			http.Error(w, "seconds must be an integer from 1 to 30", http.StatusBadRequest)
			return
		}
		timer := time.NewTimer(time.Duration(delay) * time.Second)
		defer timer.Stop()
		select {
		case <-timer.C:
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"revision": revision})
		case <-r.Context().Done():
			return
		}
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Provision HTTP example\n"))
	})
	return mux
}

func rejectStableHealth(w http.ResponseWriter, r *http.Request, revision, directHost string) bool {
	if !strings.HasSuffix(revision, "-fail-stable") || r.Host == directHost {
		return false
	}
	w.Header().Set("Cache-Control", "no-store")
	http.Error(w, "deliberate post-switch health failure", http.StatusServiceUnavailable)
	return true
}
