package main 

import (
	"net/http"
	"log"
	"fmt"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	fmt.Printf("Server running at port %s\n", server.Addr)

	log.Fatal(server.ListenAndServe())
}
