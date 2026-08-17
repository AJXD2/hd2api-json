package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Faction struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type JsonFaction map[string]string

func loadFactionJson() ([]Faction, error) {
	jsonFile, err := os.Open("json-repo/factions.json")

	if err != nil {
		fmt.Println("Error opening", err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading", err)
		return nil, err

	}

	var factions []Faction
	var rawFactions JsonFaction
	err = json.Unmarshal(byteValue, &rawFactions)
	if err != nil {
		fmt.Println("Error unmarshaling", err)
		return nil, err
	}

	for index, faction := range rawFactions {
		id, err := strconv.Atoi(index)
		if err != nil {
			fmt.Println("Error casting to integer", err)
			return nil, err

		}
		fmt.Println(id, faction)
		factions = append(factions, Faction{Index: id, Name: faction})
	}

	return factions, nil
}

func main() {
	factions, err := loadFactionJson()
	if err != nil {
		fmt.Println("Error getting faction json", err)
		return
	}
	r := chi.NewRouter()

	r.Get("/api/faction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(factions)
	})

	http.ListenAndServe(":8080", r)
}
