// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//go:generate terraform fmt -recursive ./examples/
//go:generate tfplugindocs

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	testutils "github.com/hpe-hcss/vmaas-terraform-resources/internal/test-utils"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: testutils.ProviderFunc(),
	})
}
