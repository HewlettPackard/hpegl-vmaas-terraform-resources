// +build tools

// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package tools

import (
	// document generation
	_ "github.com/golang/mock/mockgen"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
