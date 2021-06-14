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
	return map[string]*schema.Resource{
		resources.DSNetwork:       resources.NetworkData(),
		resources.DSLayout:        resources.LayoutData(),
		resources.DSGroup:         resources.GroupData(),
		resources.DSPlan:          resources.PlanData(),
		resources.DSCloud:         resources.CloudData(),
		resources.DSResourcePool:  resources.ResourcePoolData(),
		resources.DSDatastore:     resources.DatastoreData(),
		resources.DSPowerSchedule: resources.PowerScheduleData(),
		resources.DSTemplate:      resources.TemplateData(),
		resources.DSEnvironment:   resources.EnvironmentData(),
	}
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"hpegl_vmaas_instance": resources.Instances(),
		"hpegl_vmaas_snapshot": resources.Snapshots(),
	}
}

func (r Registration) ProviderSchemaEntry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			constants.LOCATION: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_LOCATION", ""),
				Description: "Location of GL VMaaS Service, can also be set with the HPEGL_VMAAS_LOCATION env var.",
			},
			constants.SPACENAME: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_SPACE_NAME", ""),
				Description: "IAM Space name of the GL VMaaS Service, can also be set with the HPEGL_VMAAS_SPACE_NAME env var.",
			},
		},
	}
}
