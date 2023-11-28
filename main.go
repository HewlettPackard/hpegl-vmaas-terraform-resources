// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	testutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/test-utils"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := &plugin.ServeOpts{
		ProviderFunc: testutils.ProviderFunc(),
	}
	if debugMode {
		// err := plugin.Debug(context.Background(), "terraform.example.com/vmaas/hpegl", opts) //nolint
		opts.Debug = debugMode
		opts.ProviderAddr = "terraform.example.com/vmaas/hpegl"
		plugin.Serve(opts)
		// if err != nil {
		// 	log.Fatal(err.Error())
		// }

		return
	}

	plugin.Serve(opts)
}
