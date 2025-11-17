package main

import (
	"log"
	"net/http"

	"github.com/akharb1192/url-shortener-infracloud/handlers"
	"github.com/akharb1192/url-shortener-infracloud/store"
)

func main() {
	st := store.NewInMemoryStore()
	h := handlers.NewHandler(st)

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenHandler)
	mux.HandleFunc("/s/", h.RedirectHandler)
	mux.HandleFunc("/metrics/top-domains", h.TopDomainsHandler)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
