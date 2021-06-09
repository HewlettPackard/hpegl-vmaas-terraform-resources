// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//go:generate terraform fmt -recursive ./examples/
//go:generate tfplugindocs

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	testutils "github.com/hpe-hcss/vmaas-terraform-resources/internal/test-utils"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := &plugin.ServeOpts{
		ProviderFunc: testutils.ProviderFunc(),
	}
	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		err := plugin.Debug(context.Background(), "terraform.example.com/vmaas/hpegl", opts)
		if err != nil {
			log.Fatal(err.Error())
		}

		return
	}

	plugin.Serve(opts)
}
