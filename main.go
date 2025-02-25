package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Main().Render(r.Context(), w)
	})

	router.Get("/clock", func(w http.ResponseWriter, r *http.Request) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		sse := datastar.NewSSE(w, r)
		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				sse.MergeFragmentTempl(Clock())
			}
		}
	})
	log.Printf("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
