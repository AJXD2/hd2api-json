package main

type Faction struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type Warbond struct {
	Index           int    `json:"index"`
	Name            string `json:"name"`
	Id              string `json:"id"`
	CreditsToUnlock int    `json:"credits_to_unlock"`
}
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

type RawPlanet struct {
	Index            int               `json:"index"`
	Name             string            `json:"name"`
	Sector           string            `json:"sector"`
	BiomeID          string            `json:"biome"`
	EnviromentalIDs  []string          `json:"environmentals"`
	Names            map[string]string `json:"names"`
	Type             string            `json:"type"`
	WeatherEffectIDs []string          `json:"weather_effects"`
}
