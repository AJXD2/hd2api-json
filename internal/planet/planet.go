package planet

import (
	"fmt"
	"sort"

	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
)

type Biome struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Environment struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Planet struct {
	Index            int               `json:"index"`
	Name             string            `json:"name"`
	Sector           string            `json:"sector"`
	Biome            Biome             `json:"biome"`
	Names            map[string]string `json:"names"`
	Type             string            `json:"type"`
	Environments     []Environment     `json:"environments"`
	WeatherEffectIDs []string          `json:"weather_effects"`
}

type rawPlanet struct {
	Index            int               `json:"index"`
	Name             string            `json:"name"`
	Sector           string            `json:"sector"`
	BiomeID          string            `json:"biome"`
	EnvironmentIDs   []string          `json:"environmentals"`
	Names            map[string]string `json:"names"`
	Type             string            `json:"type"`
	WeatherEffectIDs []string          `json:"weather_effects"`
}

func Load(repoDir string) ([]Planet, error) {
	dir := repoDir + "/planets"

	raw, err := loadRaw(dir)
	if err != nil {
		return nil, err
	}

	environments, err := loadEnvironments(dir)
	if err != nil {
		return nil, err
	}

	biomes, err := loadBiomes(dir)
	if err != nil {
		return nil, err
	}

	return hydrate(raw,
		jsonfile.IndexBy(environments, func(e Environment) string { return e.ID }),
		jsonfile.IndexBy(biomes, func(b Biome) string { return b.ID }),
	)
}

func loadRaw(dir string) ([]rawPlanet, error) {
	rawMap, err := jsonfile.Read[map[string]rawPlanet](dir + "/planets.json")
	if err != nil {
		return nil, fmt.Errorf("loading planets: %w", err)
	}

	planets, err := jsonfile.WithIndex(rawMap, func(item *rawPlanet, index int) {
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

func loadEnvironments(dir string) ([]Environment, error) {
	rawMap, err := jsonfile.Read[map[string]Environment](dir + "/environmentals.json")
	if err != nil {
		return nil, fmt.Errorf("loading environments: %w", err)
	}
	return jsonfile.WithKey(rawMap, func(item *Environment, key string) { item.ID = key })
}

func loadBiomes(dir string) ([]Biome, error) {
	rawMap, err := jsonfile.Read[map[string]Biome](dir + "/biomes.json")
	if err != nil {
		return nil, fmt.Errorf("loading biomes: %w", err)
	}
	return jsonfile.WithKey(rawMap, func(item *Biome, key string) { item.ID = key })
}

func hydrate(raw []rawPlanet, envIndex map[string]Environment, biomeIndex map[string]Biome) ([]Planet, error) {
	planets := make([]Planet, len(raw))
	for i, rp := range raw {
		envs := make([]Environment, 0, len(rp.EnvironmentIDs))
		for _, id := range rp.EnvironmentIDs {
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
