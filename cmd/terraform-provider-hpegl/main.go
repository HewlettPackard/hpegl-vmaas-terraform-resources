// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

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
