package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/go-chi/chi/v5"
)

func getFileJson[T any](filePath string) (T, error) {
	var zero T

	jsonFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error opening", err)
		return zero, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading", err)
		return zero, err

	}

	var returnData T
	err = json.Unmarshal(byteValue, &returnData)
	if err != nil {
		fmt.Println("Error unmarshaling", err)
		return zero, err
	}

	return returnData, nil
}

func loadFactionJson() ([]Faction, error) {

	rawFactions, err := getFileJson[map[string]string]("json-repo/factions.json")
	if err != nil {
		fmt.Println("Error getting json", err)
		return nil, err
	}

	factions, err := flattenNamed(rawFactions, func(id int, name string) Faction {
		return Faction{Index: id, Name: name}
	})

	sort.Slice(factions, func(i, j int) bool {
		return factions[i].Index < factions[j].Index
	})

	return factions, nil
}

func loadWarbondJson() ([]Warbond, error) {
	rawWarbonds, err := getFileJson[map[string]Warbond]("json-repo/warbonds.json")
	if err != nil {
		fmt.Println("Error getting json", err)
		return nil, err
	}

	warbonds, err := withIndex(rawWarbonds, func(item *Warbond, index int) {
		item.Index = index
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(warbonds, func(i, j int) bool {
		return warbonds[i].Index < warbonds[j].Index
	})

	return warbonds, nil
}

func getEnviromentalJson() ([]Environment, error) {
	rawEnviroments, err := getFileJson[map[string]Environment]("json-repo/planets/environmentals.json")

	if err != nil {
		fmt.Println("Error getting json", err)
		return nil, err
	}

	enviroments, err := withKey(rawEnviroments, func(item *Environment, key string) {
		item.ID = key
	})
	if err != nil {
		return nil, err
	}

	return enviroments, nil
}

func loadBiomeJson() ([]Biome, error) {
	rawBiomes, err := getFileJson[map[string]Biome]("json-repo/planets/biomes.json")
	if err != nil {
		fmt.Println("Error getting json", err)
		return nil, err
	}

	biomes, err := withKey(rawBiomes, func(item *Biome, key string) {
		item.ID = key
	})
	if err != nil {
		return nil, err
	}

	return biomes, nil
}

func loadPlanetsJson() ([]RawPlanet, error) {
	rawPlanets, err := getFileJson[map[string]RawPlanet]("json-repo/planets/planets.json")
	if err != nil {
		fmt.Println("Error getting json", err)
		return nil, err
	}

	planets, err := withIndex(rawPlanets, func(item *RawPlanet, index int) {
		item.Index = index
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(planets, func(i, j int) bool {
		return planets[i].Index < planets[j].Index
	})

	return planets, nil
}

func hydratePlanets(raw []RawPlanet, envIndex map[string]Environment, biomeIndex map[string]Biome) ([]Planet, error) {
	planets := make([]Planet, len(raw))
	for i, rp := range raw {
		envs := make([]Environment, 0, len(rp.EnviromentalIDs))
		for _, id := range rp.EnviromentalIDs {
			env, ok := envIndex[id]
			if !ok {
				return nil, fmt.Errorf("planet %q references unknown environment %q", rp.Name, id)
			}
			envs = append(envs, env)
		}

		biome, ok := biomeIndex[rp.BiomeID]
		if !ok {
			return nil, fmt.Errorf("planet %q references unknown biome %q", rp.Name, rp.BiomeID)
		}

		planets[i] = Planet{
			Index:            rp.Index,
			Name:             rp.Name,
			Sector:           rp.Sector,
			Biome:            biome,
			Names:            rp.Names,
			Type:             rp.Type,
			Environments:     envs,
			WeatherEffectIDs: rp.WeatherEffectIDs,
		}
	}
	return planets, nil
}

func staticJSON(body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Write(body)
	}
}

func main() {
	factions := must(loadFactionJson())
	warbonds := must(loadWarbondJson())
	environments := must(getEnviromentalJson())
	biomes := must(loadBiomeJson())
	rawPlanets := must(loadPlanetsJson())
	planets := must(hydratePlanets(rawPlanets,
		indexBy(environments, func(e Environment) string { return e.ID }),
		indexBy(biomes, func(b Biome) string { return b.ID }),
	))

	factionsJSON := mustMarshal(factions)
	warbondsJSON := mustMarshal(warbonds)
	planetsJSON := mustMarshal(planets)

	r := chi.NewRouter()
	r.Get("/api/faction", staticJSON(factionsJSON))
	r.Get("/api/warbonds", staticJSON(warbondsJSON))
	r.Get("/api/planets", staticJSON(planetsJSON))

	log.Fatal(http.ListenAndServe(":8080", r))
}
