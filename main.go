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

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/metrics", cfg.metricsHandler)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)

	r.Mount("/api", apiRouter)

	corsMux := middlewareCors(r)

	s := http.Server{
		Addr:	 port,
		Handler: corsMux,
	}

	fmt.Printf("Server started on port%s", port)
	log.Fatal(s.ListenAndServe())
}

