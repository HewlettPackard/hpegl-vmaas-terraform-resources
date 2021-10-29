// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package validations

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ValidateUniqueNameInList(i interface{}, p cty.Path) diag.Diagnostics {
	var diagErr diag.Diagnostics
	lMap := utils.GetlistMap(i)
	if lMap == nil {
		return diagErr
	}
	// check for duplicate name
	nameMap := make(map[string]struct{})
	for _, l := range lMap {
		name := l["name"].(string)
		if _, ok := nameMap[name]; ok {
			diagErr = append(diagErr, diag.Errorf("duplicate name entry %s", name)...)
		}
		nameMap[name] = struct{}{}
	}

	return diagErr
}
