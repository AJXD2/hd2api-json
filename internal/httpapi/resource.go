package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Resource struct {
	path  string
	body  []byte
	byKey map[string][]byte
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

func NewIndexed[T any](path string, load func() ([]T, error), key func(T) string) (Resource, error) {
	items, err := load()
	if err != nil {
		return Resource{}, fmt.Errorf("loading %s: %w", path, err)
	}

	body, err := json.Marshal(items)
	if err != nil {
		return Resource{}, fmt.Errorf("marshaling %s: %w", path, err)
	}

	byKey := make(map[string][]byte, len(items))
	for _, item := range items {
		b, err := json.Marshal(item)
		if err != nil {
			return Resource{}, fmt.Errorf("marshaling %s item: %w", path, err)
		}
		byKey[key(item)] = b
	}

	return Resource{path: path, body: body, byKey: byKey}, nil
}

func (r Resource) Register(router chi.Router) {
	router.Get(r.path, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, r.body)
	})

	if r.byKey == nil {
		return
	}

	router.Get(r.path+"/{key}", func(w http.ResponseWriter, req *http.Request) {
		b, ok := r.byKey[chi.URLParam(req, "key")]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, b)
	})
}

func writeJSON(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Write(body)
}
