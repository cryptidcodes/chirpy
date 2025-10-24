package main

import (
	"sync/atomic"
	"fmt"
	"net/http"
)

type apiConfig struct {
		fileserverHits atomic.Int32
	}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// init config struct
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// create a new http.ServeMux to handle requests
	mux := http.NewServeMux()

	// handle /app/ with middleware to count hits
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// additional endpoint handlers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	// admin endpoint handlers
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	
	// create a new http.Server struct
	server := &http.Server{
		Addr:	":" + port,
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	// start the server
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}

}