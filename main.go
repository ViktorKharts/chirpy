package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ViktorKharts/chirpy/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

const (
	PORT = ":8080"
	JWT_SECRET = "JWT_SECRET"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileserverHits: 0,
		DB:		*db,	
		jwtSecret:	os.Getenv(JWT_SECRET),
	}
	r := chi.NewRouter()

	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	// /api 
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/reset", cfg.resetMetricsHandler)
	apiRouter.Post("/login", cfg.loginHandler)
	apiRouter.Post("/refresh", cfg.refreshTokenHandler)
	apiRouter.Post("/revoke", cfg.revokeTokenHandler)

	// /api/chirps
	apiRouter.Get("/chirps", cfg.getChirpsHandler)
	apiRouter.Get("/chirps/{chirpID}", cfg.getChirpHandler)
	apiRouter.Post("/chirps", cfg.createChirpsHandler)
	apiRouter.Delete("/chirps/{chirpID}", cfg.deleteChirpHandler)
	
	// /api/users
	apiRouter.Post("/users", cfg.createUsersHandler)
	apiRouter.Put("/users", cfg.updateUsersHandler)

	// /admin
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.metricsHandler)

	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	s := http.Server{
		Addr:	 PORT,
		Handler: corsMux,
	}

	fmt.Printf("Server started on port%s", PORT)
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
	
