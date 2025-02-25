package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	datastar "github.com/starfederation/datastar/sdk/go"
)

var output = []string{
	"User 1",
	"User 2",
	"Cat 1",
	"Dog 1",
}

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

	router.Get("/output", func(w http.ResponseWriter, r *http.Request) {
		var signals struct {
			In string `json:"in"`
		}
		datastar.ReadSignals(r, &signals)

		var out []string
		if signals.In != "" {
			for i := range output {
				// filter output
				if strings.HasPrefix(strings.ToLower(output[i]), strings.ToLower(signals.In)) {
					out = append(out, output[i])
				}
			}
		}

		sse := datastar.NewSSE(w, r)
		sse.MergeFragmentTempl(Output(out))
	})
	log.Printf("Server started at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
