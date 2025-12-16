package main 

import (
	"net/http"
	"log"
	"fmt"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	w.WriteHeader(200)
	cfg.fileserverHits.Swap(0)
}

func main() {
	mux := http.NewServeMux()
	cfg := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerGetFileserverHits)

	mux.HandleFunc("POST /admin/reset", cfg.handlerResetFileserverHits)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	fmt.Printf("Server running at port %s\n", server.Addr)

	log.Fatal(server.ListenAndServe())
}
