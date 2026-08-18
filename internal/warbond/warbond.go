package warbond

import (
	"fmt"
	"sort"

	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
)

type Warbond struct {
	Index           int    `json:"index"`
	Name            string `json:"name"`
	ID              string `json:"id"`
	CreditsToUnlock int    `json:"credits_to_unlock"`
}

func Load(repoDir string) ([]Warbond, error) {
	raw, err := jsonfile.Read[map[string]Warbond](repoDir + "/warbonds.json")
	if err != nil {
		return nil, fmt.Errorf("loading warbonds: %w", err)
	}

	warbonds, err := jsonfile.WithIndex(raw, func(item *Warbond, index int) {
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
