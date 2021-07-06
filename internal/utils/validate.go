package utils

// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SkipField() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return old != ""
	}
}
