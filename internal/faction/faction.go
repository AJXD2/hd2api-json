package faction

import (
	"fmt"
	"sort"

	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
)

type Faction struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

func Load(repoDir string) ([]Faction, error) {
	raw, err := jsonfile.Read[map[string]string](repoDir + "/factions.json")
	if err != nil {
		return nil, fmt.Errorf("loading factions: %w", err)
	}

	factions, err := jsonfile.FlattenNamed(raw, func(id int, name string) Faction {
		return Faction{Index: id, Name: name}
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(factions, func(i, j int) bool {
		return factions[i].Index < factions[j].Index
	})

	return factions, nil
}
