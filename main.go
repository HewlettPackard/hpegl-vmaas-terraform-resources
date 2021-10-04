// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/spf13/viper"

	testutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/test-utils"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	libUtils "github.com/hewlettpackard/hpegl-provider-lib/pkg/utils"
)

func main() {
	// Read config file for acceptance test if TF_ACC sets
	libUtils.ReadAccConfig(".")
	if viper.GetBool("TF_ACC") {
		pkgutils.ReadSkip()
	}

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
