package main 

import (
	"net/http"
	"log"
	"fmt"
	"sync/atomic"
	"database/sql"
	"os"
	"github.com/joho/godotenv"
	"github.com/arey-dev/chirpy/internal/database"
)

import _ "github.com/lib/pq"

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerGetFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	body := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerResetFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Error deleting users", nil)
		return
	}

	err := cfg.db.DeleteUsers(req.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	cfg := apiConfig{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerGetFileserverHits)

	mux.HandleFunc("POST /admin/reset", cfg.handlerResetFileserverHits)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, req *http.Request) {
		handlerCreateUser(&cfg, w, req)
	})

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	fmt.Printf("Server running at port %s\n", server.Addr)

	log.Fatal(server.ListenAndServe())
}
