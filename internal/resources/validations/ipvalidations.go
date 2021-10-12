package validations

import (
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
