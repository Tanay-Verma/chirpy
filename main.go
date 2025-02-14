package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Tanay-Verma/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("DB_URL must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	config := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))),
	)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", config.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", config.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", config.handlerLogin)
	mux.HandleFunc("POST /api/refresh", config.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", config.handlerRevoke)

	mux.HandleFunc("POST /api/chirps", config.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", config.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", config.handlerDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", config.handlerPolkaWebhooks)

	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", config.handlerReset)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
