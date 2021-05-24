// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ListToStringSlice(src interface{}) ([]string, error) {
	list, ok := src.([]interface{})
	if !ok {
		return nil, typeError("[]interface{}", src)
	}
	dst := make([]string, 0, len(list))
	for _, s := range list {
		d, ok := s.(string)
		if !ok {
			return nil, fmt.Errorf("error unable to convert %v (%T) to string", s, s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func ListToIntSlice(src interface{}) ([]int, error) {
	list, ok := src.([]interface{})
	if !ok {
		return nil, typeError("[]interface{}", src)
	}

	dst := make([]int, 0, len(list))
	for _, s := range list {
		d, ok := s.(int)
		if !ok {
			return nil, typeError("Int", s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func ListToMap(src interface{}) ([]map[string]interface{}, error) {

	list, ok := src.([]interface{})
	if !ok {
		return nil, typeError("[]interface{}", src)
	}
	dst := make([]map[string]interface{}, 0, len(list))
	for _, s := range list {
		d, ok := s.(map[string]interface{})
		if !ok {
			return nil, typeError("Map", s)
		}
		dst = append(dst, d)
	}
	return dst, nil
}

func SetToMap(src interface{}) (map[string]interface{}, error) {
	set, ok := src.(*schema.Set)
	if !ok {
		return nil, typeError("*Set", src)
	}
	if set == nil {
		return nil, fmt.Errorf("error nil set")
	}
	list := set.List()
	return list[0].(map[string]interface{}), nil
}

func MapTopMap(src interface{}) (map[string]interface{}, error) {
	dst, ok := src.(map[string]interface{})
	if !ok {
		return nil, typeError("map[string]interface{}", src)
	}
	return dst, nil
}

func JsonNumber(in interface{}) json.Number {
	if a, ok := in.(int); ok {
		return json.Number(strconv.Itoa(a))
	}
	return json.Number(in.(string))
}

func typeError(wantedType string, in interface{}) error {
	return fmt.Errorf("error unable to convert %v (%T) to %s", in, in, wantedType)
}
