package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

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
	apiRouter.Post("/validate_chirp", validateChirpHandler)

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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		Body string `json:"body"`
    	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

    	decoder := json.NewDecoder(r.Body)
    	params := requestPayload{}
    	err := decoder.Decode(&params)
    	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
    	}

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJson(w, http.StatusOK, validResponse{
		CleanedBody: cleanBody(params.Body),
	})
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

func cleanBody(ds string) (cs string) {	
	dw := []string{"kerfuffle", "sharbert", "fornax"}
	s := strings.Split(ds, " ")
	for i, v := range s {
		if slices.Contains[[]string, string](dw, strings.ToLower(v)) {
			s[i] = "****"
		}
	}
	return strings.Join(s, " ")
}
	
