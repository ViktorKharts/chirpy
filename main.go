package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ViktorKharts/chirpy/internal/database"
	"github.com/go-chi/chi"
)

const port = ":8080"

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileserverHits: 0,
		DB:		*db,	
	}
	r := chi.NewRouter()

	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)
	apiRouter.Get("/chirps", cfg.getChirpsHandler)
	apiRouter.Get("/chirps/{chirpID}", cfg.getChirpHandler)
	apiRouter.Post("/chirps", cfg.createChirpsHandler)

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

func respondWithError(w http.ResponseWriter, s int, m string) {
	type errorResponse struct {
		Error string `json:"error"`	
	}

	respondWithJson(w, s, errorResponse {
		Error: m,
	})
}

func respondWithJson(w http.ResponseWriter, s int, p interface{}) {
	w.Header().Set("Content-Type", "application/json")	
	j, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(s)
	w.Write(j)
}
	
