package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platfrom       string
	secret         string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	config := apiConfig{}
	config.dbQueries = dbQueries
	config.platfrom = os.Getenv("PLATFORM")
	config.secret = os.Getenv("SECRET")

	handler := http.NewServeMux()
	fileHandler := config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	handler.Handle("/app/", fileHandler)
	handler.HandleFunc("GET /api/healthz", handleHealthz)
	handler.HandleFunc("POST /api/chirps", config.chirpPostHandler)
	handler.HandleFunc("GET /api/chirps", config.chirpGetHandler)
	handler.HandleFunc("GET /api/chirps/{chirpID}", config.chirpGetSingleHandler)
	handler.HandleFunc("POST /api/users", config.handleUsersPost)
	handler.HandleFunc("PUT /api/users", config.handleUsersPut)
	handler.HandleFunc("POST /api/login", config.handleLogin)
	handler.HandleFunc("POST /api/revoke", config.revokeHandler)
	handler.HandleFunc("POST /api/refresh", config.refreshHandler)
	handler.HandleFunc("GET /admin/metrics", config.handleMetrics)
	handler.HandleFunc("POST /admin/reset", config.handleReset)

	server := http.Server{}
	server.Handler = handler
	server.Addr = ":8080"

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
