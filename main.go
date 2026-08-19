package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ajxd2/helldivers-json-api/internal/faction"
	"github.com/ajxd2/helldivers-json-api/internal/httpapi"
	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
	"github.com/ajxd2/helldivers-json-api/internal/planet"
	"github.com/ajxd2/helldivers-json-api/internal/region"
	"github.com/ajxd2/helldivers-json-api/internal/warbond"
	"github.com/go-chi/chi/v5"
)

const repoDir = "json-repo"

func main() {
	resources := []httpapi.Resource{
		jsonfile.Must(httpapi.New("/api/faction", func() ([]faction.Faction, error) { return faction.Load(repoDir) })),
		jsonfile.Must(httpapi.NewIndexed("/api/warbonds",
			func() ([]warbond.Warbond, error) { return warbond.Load(repoDir) },
			func(w warbond.Warbond) string { return strconv.Itoa(w.Index) },
		)),
		jsonfile.Must(httpapi.NewIndexed("/api/planets",
			func() ([]planet.Planet, error) { return planet.Load(repoDir) },
			func(p planet.Planet) string { return strconv.Itoa(p.Index) },
		)),
		jsonfile.Must(httpapi.NewIndexed("/api/regions",
			func() ([]region.Region, error) { return region.Load(repoDir) },
			func(r region.Region) string { return strconv.Itoa(r.ID) },
		)),
	}

	r := chi.NewRouter()
	for _, res := range resources {
		res.Register(r)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
