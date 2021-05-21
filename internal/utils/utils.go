// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import "fmt"

func listToStringSlice(src []interface{}) ([]string, error) {
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

func listToIntSlice(src []interface{}) ([]int, error) {
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
