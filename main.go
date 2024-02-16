package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	cfg := apiConfig{0}
	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMatricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("/reset", cfg.resetMetricsHandler)

	corsMux := middlewareCors(mux)

	s := http.Server{
		Addr:	 port,
		Handler: corsMux,
	}

	fmt.Printf("Server started on port%s", port)
	log.Fatal(s.ListenAndServe())
}

