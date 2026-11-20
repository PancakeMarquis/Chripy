package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	serv := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}
	serverMux.HandleFunc("/healthz", handlerReadiness)
	serverMux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("/reset", apiCfg.handlerReset)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())

}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(text))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
