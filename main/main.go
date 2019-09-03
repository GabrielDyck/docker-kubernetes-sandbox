package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/health-check",healthCheck)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9290"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/health-check" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "OK!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
