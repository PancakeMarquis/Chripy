package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/PancakeMarquis/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}
	dbQueries := database.New(db)

	plat := os.Getenv("PLATFORM")
	if plat == "" {
		log.Fatal("PLATFORM must be set")
	}

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, db: dbQueries, platform: plat}
	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	serv := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())
}
