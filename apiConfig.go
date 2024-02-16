package main

import (
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

func (c *apiConfig) middlewareMatricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits += 1
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(c.fileserverHits)))
}

func (c *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
	
