// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package validations

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func errsTodiags(errs []error) diag.Diagnostics {
	d := make([]diag.Diagnostic, 0, len(errs))
	for _, err := range errs {
		d = append(d, diag.FromErr(err)...)
	}

	return d
}

// ValidateIPAddress Validates IP address
func ValidateIPAddress(i interface{}, p cty.Path) diag.Diagnostics {
	if i == nil {
		return nil
	}

	_, errs := validation.IsIPv4Address(i, "")

	return errsTodiags(errs)
}

// ValidateCidr validate cidr
func ValidateCidr(i interface{}, p cty.Path) diag.Diagnostics {
	if i == nil {
		return nil
	}

	_, errs := validation.IsCIDR(i, "")

	return errsTodiags(errs)
}

// ValidateCidr validate cidr or IP Address
func ValidateIPorCidr(i interface{}, p cty.Path) diag.Diagnostics {
	var errors []error
	if i == nil {
		return nil
	}

	_, errsCidr := validation.IsCIDR(i, "")

	if errsCidr != nil {
		_, errsIpv4 := validation.IsIPv4Address(i, "")
		if errsIpv4 != nil {
			errors = append(errors, fmt.Errorf("expected %s to contain a valid IPv4 address or CIDR", i.(string)))
		}
	}

	return errsTodiags(errors)
}
