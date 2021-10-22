// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package validations

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(valid, ignoreCase))
}

func IntBetween(min int, max int) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntBetween(min, max))
}

func IntAtLeast(min int) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntAtLeast(min))
}
