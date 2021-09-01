// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	testutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/test-utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
)

func main() {
	// Read config file for acceptance test if TF_ACC sets
	utils.ReadAccConfig(".")

	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := &plugin.ServeOpts{
		ProviderFunc: testutils.ProviderFunc(),
	}
	if debugMode {
		err := plugin.Debug(context.Background(), "terraform.example.com/vmaas/hpegl", opts)
		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)
}
