package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Resource struct {
	path string
	body []byte
}

func New[T any](path string, load func() (T, error)) (Resource, error) {
	data, err := load()
	if err != nil {
		return Resource{}, fmt.Errorf("loading %s: %w", path, err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return Resource{}, fmt.Errorf("marshaling %s: %w", path, err)
	}

	return Resource{path: path, body: body}, nil
}

func (r Resource) Register(router chi.Router) {
	router.Get(r.path, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Write(r.body)
	})
}
