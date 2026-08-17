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

type Planet struct {
	Index          int               `json:"index"`
	Name           string            `json:"name"`
	Sector         string            `json:"sector"`
	Biome          string            `json:"biome"`
	Enviromentals  []string          `json:"enviromentals"`
	Names          map[string]string `json:"names"`
	Type           string            `json:"type"`
	WeatherEffects []string          `json:"weather_effects"`
}
