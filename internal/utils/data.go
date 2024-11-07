// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Data struct {
	d *schema.ResourceData
	// errors will hold list of errors for each attrib
	// attrib name will be the key
	errors map[string][]string
}

// NewData returns new Data instance
func NewData(d *schema.ResourceData) *Data {
	return &Data{
		d:      d,
		errors: make(map[string][]string),
	}
}

func (d *Data) Error() error {
	if d.hasError() {
		errStr := ""
		for k, v := range d.errors {
			errStr += k + " : " + strings.Join(v, ",") + "\n"
		}

		return errors.New(errStr)
	}

	return nil
}

func (d *Data) hasError() bool {
	return len(d.errors) > 0
}

func (d *Data) err(key, msg string) {
	d.errors[key] = append(d.errors[key], msg)
}

func GetlistMap(src interface{}) []map[string]interface{} {
	list, ok := src.([]interface{})
	if !ok {
		return GetSMap(src)
	}
	dst := make([]map[string]interface{}, 0, len(list))
	for _, s := range list {
		ds, ok := s.(map[string]interface{})
		if ok {
			dst = append(dst, ds)
		}
	}

	return dst
}

// GetListMap take key as parameter and returns []map[string]interfac{}.
// This function can be used for retrieving list of map or list of set
func (d *Data) GetListMap(key string) []map[string]interface{} {
	src := d.get(key)
	if src == nil {
		return nil
	}

	return GetlistMap(src)
}

func (d *Data) GetChangedListMap(key string) ([]map[string]interface{}, []map[string]interface{}) {
	org, new := d.d.GetChange(key)
	var orgmap, newmap []map[string]interface{}
	if org != nil {
		orgmap = GetlistMap(org)
	}
	if new != nil {
		newmap = GetlistMap(new)
	}

	return orgmap, newmap
}

func (d *Data) HasChanged(key string) bool {
	src := d.d.HasChange(key)

	return src
}

func (d *Data) GetChangedMap(key string) (map[string]interface{}, map[string]interface{}) {
	org, new := d.d.GetChange(key)

	orgmap, ok := org.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	newmap, ok := new.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return orgmap, newmap
}

func (d *Data) get(key string) interface{} {
	return d.d.Get(key)
}

func (d *Data) GetID() int {
	id, err := ParseInt(d.d.Id())
	if err != nil {
		d.err("id", ErrInvalidType)

		return NAN
	}

	return int(id)
}
func (d *Data) GetID64() int64 {
	id, err := ParseInt(d.d.Id())
	if err != nil {
		d.err("id", ErrInvalidType)

		return NAN
	}

	return int64(id)
}

// GetIDString returns ID as string
func (d *Data) GetIDString() string {
	return d.d.Id()
}

// SetID should either be int or string
func (d *Data) SetID(v interface{}) {
	var stringID string
	switch x := v.(type) {
	case int:
		stringID = strconv.Itoa(x)
	case string:
		stringID = x
	default:
		panic("Invalid data on SetID")
	}

	d.d.SetId(stringID)
}

// nolint
func (d *Data) SetId(v string) {
	d.d.SetId(v)
}

// nolint
func (d *Data) Id() string {
	return d.d.Id()
}

func (d *Data) set(key string, value interface{}) error {
	return d.d.Set(key, value)
}

// GetStringList returns list of string
func (d *Data) GetStringList(key string, ignore ...bool) []string {
	src, ok := d.getOk(key, ignore)
	if !ok {
		return nil
	}
	list, ok := src.([]interface{})
	if !ok {
		return nil
	}
	dst := make([]string, 0, len(list))
	for _, s := range list {
		ds, ok := s.(string)
		if ok {
			dst = append(dst, ds)
		}
	}

	return dst
}

func (d *Data) GetInt(key string, ignore ...bool) int {
	valInter, ok := d.getOk(key, ignore)
	if !ok {
		return NAN
	}
	valInt, ok := valInter.(int)
	var err error
	if !ok {
		valString, ok := valInter.(string)
		if ok {
			valInt, err = strconv.Atoi(valString)
			if err != nil {
				d.err(key, ErrInvalidType+valString)

				return NAN
			}

			return valInt
		}
		d.err(key, ErrInvalidType)

		return NAN
	}

	return valInt
}

func (d *Data) GetInt64(key string, ignore ...bool) int64 {
	valInter, ok := d.getOk(key, ignore)
	if !ok {
		return NAN
	}
	valInt, ok := valInter.(int64)
	var err error
	if !ok {
		valString, ok := valInter.(string)
		if ok {
			valInt, err = ParseInt(valString)
			if err != nil {
				d.err(key, ErrInvalidType+valString)

				return NAN
			}

			return valInt
		}
		d.err(key, ErrInvalidType)

		return NAN
	}

	return valInt
}

func (d *Data) getOk(key string, ignore []bool) (interface{}, bool) {
	val, ok := d.d.GetOk(key)
	if len(ignore) != 0 && !ignore[0] {
		if !ok {
			d.err(key, ErrKeyNotDefined)
		}
	}

	return val, ok
}

func (d *Data) GetOk(key string) (interface{}, bool) {
	return d.d.GetOk(key)
}

// GetSMap for get map for a Set
func GetSMap(src interface{}) []map[string]interface{} {
	set, ok := src.(*schema.Set)
	if !ok {
		return nil
	}
	if set == nil {
		return nil
	}
	list := set.List()
	if len(list) == 0 {
		return nil
	}
	mapList := make([]map[string]interface{}, len(list))
	for i, l := range list {
		val := l.(map[string]interface{})
		mapList[i] = val
	}

	return mapList
}

func (d *Data) GetMap(key string, ignore ...bool) map[string]interface{} {
	src, ok := d.getOk(key, ignore)
	if !ok {
		return nil
	}
	dst, ok := src.(map[string]interface{})
	if !ok {
		return nil
	}

	return dst
}

func (d *Data) GetString(key string, ignore ...bool) string {
	val, ok := d.getOk(key, ignore)
	if !ok {
		return ""
	}
	if val != nil {
		return val.(string)
	}
	d.err(key, ErrInvalidType)

	return ""
}

func (d *Data) GetJSONNumber(key string, ignore ...bool) json.Number {
	in, ok := d.getOk(key, ignore)
	if !ok {
		return "0"
	}

	return JSONNumber(in)
}

func (d *Data) GetBool(key string) bool {
	if val := d.get(key); val != nil {
		return val.(bool)
	}
	d.err(key, ErrInvalidType)

	return false
}

func (d *Data) SetString(key string, value string) {
	if err := d.set(key, value); err != nil {
		d.err(key, ErrSet+" : "+err.Error())
	}
}

func (d *Data) ListToIntSlice(key string, ignore ...bool) []int {
	src, ok := d.getOk(key, ignore)
	if !ok {
		return nil
	}
	list, ok := src.([]interface{})
	if !ok {
		return nil
	}

	dst := make([]int, 0, len(list))
	for i, s := range list {
		ds, ok := s.(int)
		if !ok {
			d.err(key, ErrInvalidType+" at index "+strconv.Itoa(i))
		} else {
			dst = append(dst, ds)
		}
	}

	return dst
}

func (d *Data) Set(key string, val interface{}) error {
	err := d.d.Set(key, val)
	if err != nil {
		d.err(key, ErrSet+", "+err.Error())
	}

	return err
}
