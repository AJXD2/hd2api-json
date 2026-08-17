package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Faction struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

func main() {
	r := chi.NewRouter()

	r.Get("/api/faction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode([]Faction{
			{Name: "Terminids", Index: 1},
			{Name: "Automatons", Index: 2},
		})
	})

	http.ListenAndServe(":8080", r)
}
