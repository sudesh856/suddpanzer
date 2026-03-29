package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/jwt", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"error":"missing auth"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"got":"%s"}`, auth)
	})

	http.HandleFunc("/apikey", func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"error":"missing api key"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"got":"%s"}`, key)
	})

	log.Println("auth test server listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
