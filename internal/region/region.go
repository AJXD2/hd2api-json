package region

import (
	"fmt"

	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
)

type Region struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Faction     string `json:"region_faction"`
	Type        string `json:"region_type"`
}

func Load(repoDir string) ([]Region, error) {
	raw, err := jsonfile.Read[map[string]Region](repoDir + "/planets/planetRegion.json")
	if err != nil {
		return nil, fmt.Errorf("loading regions: %w", err)
	}

	regions, err := jsonfile.WithIndex(raw, func(item *Region, index int) {
		item.ID = index
	})
	if err != nil {
		return nil, err
	}

	return regions, nil
}
