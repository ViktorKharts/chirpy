package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
		
		    <body>
		        <h1>Welcome, Chirpy Admin</h1>
		        <p>Chirpy has been visited %d times!</p>
		    </body>
		
		</html>
	`, c.fileserverHits)))
}

func (c *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
	
