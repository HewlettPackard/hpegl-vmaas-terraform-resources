// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/hpegl-provider-lib/pkg/registration"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/resources"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/constants"
)

// Assert that Registration implements the ServiceRegistration interface
var _ registration.ServiceRegistration = (*Registration)(nil)

type Registration struct{}

func (r Registration) Name() string {
	return constants.ServiceName
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return nil
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"hpegl_vmaas_vm": resources.VirtualMachine(),
	}
}

func (r Registration) ProviderSchemaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			constants.APIURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_API_URL", ""),
				Description: "The URL to use for the VMaaS API, can also be set with the HPEGL_VMAAS_API_URL env var",
			},
			constants.LOCATION: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_LOCATION", ""),
				Description: "Location of GL VMaaS Service, can also be set with the HPEGL_VMAAS_LOCATION env var",
			},
		},
	}
}
