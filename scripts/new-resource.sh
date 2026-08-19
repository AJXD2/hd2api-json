#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 2 ]; then
	echo "usage: $0 <resource-name> <source.json>" >&2
	echo "example: $0 stratagem stratagems.json" >&2
	exit 1
fi

name="$1"
source_json="$2"
pkg_dir="internal/$name"
type_name="$(tr '[:lower:]' '[:upper:]' <<<"${name:0:1}")${name:1}"

if [ -e "$pkg_dir" ]; then
	echo "error: $pkg_dir already exists" >&2
	exit 1
fi

mkdir -p "$pkg_dir"

cat >"$pkg_dir/$name.go" <<EOF
package $name

import (
	"fmt"
	"sort"

	"github.com/ajxd2/helldivers-json-api/internal/jsonfile"
)

type $type_name struct {
	Index int    \`json:"index"\`
	Name  string \`json:"name"\`
}

func Load(repoDir string) ([]$type_name, error) {
	raw, err := jsonfile.Read[map[string]$type_name](repoDir + "/$source_json")
	if err != nil {
		return nil, fmt.Errorf("loading $name: %w", err)
	}

	items, err := jsonfile.WithIndex(raw, func(item *$type_name, index int) {
		item.Index = index
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Index < items[j].Index
	})

	return items, nil
}
EOF

echo "wrote $pkg_dir/$name.go"
echo "next: fill in the $type_name struct fields to match $source_json, then wire it into main.go"
