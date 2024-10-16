package utils

// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SkipField() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return old != ""
	}
}

// SkipEmptyField it to skip diff check when user didn't set any attribute in tf file but state file is updated
func SkipEmptyField() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return new == "" || new == "0"
	}
}
