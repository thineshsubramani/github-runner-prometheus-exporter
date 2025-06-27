package server

import (
	"log"
	"net/http"
)

// Start runs the HTTP server with /metrics and optional /health
func Start(listenAddr string, handler http.Handler) {
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle("/metrics", handler)

	// health check (Optional)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Printf("🚀 Exporter listening at http://localhost:%s/metrics", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("❌ HTTP server error: %v", err)
	}
}
