//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RouterNat struct {
	diff *schema.ResourceDiff
}

func NewRouterNatValidate(diff *schema.ResourceDiff) *RouterNat {
	return &RouterNat{
		diff: diff,
	}
}

func (r *RouterNat) DiffValidate() error {
	return r.validateDandSnat()
}

func (r *RouterNat) validateDandSnat() error {
	actionPath := "config.0.action"
	if r.diff.HasChange(actionPath) {
		action := r.diff.Get(actionPath)
		switch action {
		case "DNAT":
			if r.diff.Get("destination_network") == "" {
				return fmt.Errorf("destination_network should be set for DNAT")
			}
		case "SNAT":
			if r.diff.Get("source_network") == "" {
				return fmt.Errorf("source_network should be set for SNAT")
			}
		}
	}

	return nil
}
