package main

import (
	"fmt"
	"os"
	"strconv"
)

func must[T any](val T, err error) T {
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	return val
}

func withIndex[T any](rawMap map[string]T, setIndex func(item *T, index int)) ([]T, error) {
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

func flattenNamed[T any](rawMap map[string]string, build func(id int, value string) T) ([]T, error) {
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
