package main

import (
	"log"
	"net/http"

	"github.com/ajxd2/helldivers-json-api/internal/faction"
	"github.com/ajxd2/helldivers-json-api/internal/httpapi"
	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
	"github.com/ajxd2/helldivers-json-api/internal/planet"
	"github.com/ajxd2/helldivers-json-api/internal/warbond"
	"github.com/go-chi/chi/v5"
)

const repoDir = "json-repo"

func main() {
	resources := []httpapi.Resource{
		jsonfile.Must(httpapi.New("/api/faction", func() ([]faction.Faction, error) { return faction.Load(repoDir) })),
		jsonfile.Must(httpapi.New("/api/warbonds", func() ([]warbond.Warbond, error) { return warbond.Load(repoDir) })),
		jsonfile.Must(httpapi.New("/api/planets", func() ([]planet.Planet, error) { return planet.Load(repoDir) })),
	}

	r := chi.NewRouter()
	for _, res := range resources {
		res.Register(r)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
