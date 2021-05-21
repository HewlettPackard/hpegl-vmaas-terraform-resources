// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"fmt"
)

func ListToStringSlice(src []interface{}) ([]string, error) {
	dst := make([]string, 0, len(src))
	for _, s := range src {
		d, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("unale to convert %v (%T) to string", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func ListToIntSlice(src []interface{}) ([]int, error) {
	dst := make([]int, 0, len(src))
	for _, s := range src {
		d, ok := s.(int)
		if !ok {
			return nil, fmt.Errorf("unale to convert %v (%T) to int", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func ListToMap(src []interface{}) ([]map[string]interface{}, error) {
	var dst []map[string]interface{}

	for _, s := range src {
		d, ok := s.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unale to convert %v (%T) to int", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}
