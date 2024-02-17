package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

const port = ":8080"

func main() {
	cfg := apiConfig{0}
	r := chi.NewRouter()

	fsHandler := cfg.middlewareMatricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	r.Get("/healthz", readinessHandler)
	r.Get("/metrics", cfg.metricsHandler)
	r.Get("/reset", cfg.resetMetricsHandler)

	corsMux := middlewareCors(r)

	s := http.Server{
		Addr:	 port,
		Handler: corsMux,
	}

	fmt.Printf("Server started on port%s", port)
	log.Fatal(s.ListenAndServe())
}

