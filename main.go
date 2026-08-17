package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"

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
	}

	warbonds, err := withIndex(rawWarbonds, func(item *Warbond, index int) {
		item.Index = index
	})

	sort.Slice(warbonds, func(i, j int) bool {
		return warbonds[i].Index < warbonds[j].Index
	})

	return warbonds, nil
}

func main() {
	factions := must(loadFactionJson())
	warbonds := must(loadWarbondJson())
	r := chi.NewRouter()

	r.Get("/api/faction", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(factions)
	})

	r.Get("/api/warbonds", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(warbonds)
	})

	http.ListenAndServe(":8080", r)
}
