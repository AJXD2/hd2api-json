package jsonfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func Read[T any](path string) (T, error) {
	var zero T

	f, err := os.Open(path)
	if err != nil {
		return zero, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return zero, err
	}

	var data T
	if err := json.Unmarshal(b, &data); err != nil {
		return zero, err
	}

	return data, nil
}

func WithIndex[T any](rawMap map[string]T, setIndex func(item *T, index int)) ([]T, error) {
	result := make([]T, 0, len(rawMap))
	for indexStr, obj := range rawMap {
		id, err := strconv.Atoi(indexStr)
		if err != nil {
			return nil, err
		}
		setIndex(&obj, id)
		result = append(result, obj)
	}
	return result, nil
}

func WithKey[T any](rawMap map[string]T, setKey func(item *T, key string)) ([]T, error) {
	result := make([]T, 0, len(rawMap))
	for key, obj := range rawMap {
		setKey(&obj, key)
		result = append(result, obj)
	}
	return result, nil
}

func FlattenNamed[T any](rawMap map[string]string, build func(id int, value string) T) ([]T, error) {
	result := make([]T, 0, len(rawMap))
	for idStr, value := range rawMap {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		result = append(result, build(id, value))
	}
	return result, nil
}

func IndexBy[T any](items []T, key func(T) string) map[string]T {
	idx := make(map[string]T, len(items))
	for _, item := range items {
		idx[key(item)] = item
	}
	return idx
}

func Must[T any](val T, err error) T {
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	return val
}
