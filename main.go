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

	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.metricsHandler)

	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	s := http.Server{
		Addr:	 port,
		Handler: corsMux,
	}

	fmt.Printf("Server started on port%s", port)
	log.Fatal(s.ListenAndServe())
}

