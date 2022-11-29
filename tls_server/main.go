package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", middle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello world"))
	})))
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func middle(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		// とりあえずヘッダの有無だけ
		if auth == "" {
			http.Error(w, "Unauthorize request", http.StatusUnauthorized)
			return
		}
	}

	return http.HandlerFunc(fn)
}
