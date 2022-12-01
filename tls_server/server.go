package main

import (
	"net/http"
	"time"
)

func NewServeMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.Handle("/", middle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello world"))
	})))
	mux.HandleFunc("/sleep15", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte("wake up. I was sleeping until 15 seconds."))
	})

	return mux, nil
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
