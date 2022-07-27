// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/registration"
)

// Assert that Registration implements the ServiceRegistration interface
var _ registration.ServiceRegistration = (*Registration)(nil)

type Registration struct{}

func (r Registration) Name() string {
	return constants.ServiceName
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.DSNetwork:          resources.NetworkData(),
		resources.DSNetworkType:      resources.NetworkTypeData(),
		resources.DSNetworkPool:      resources.NetworkPoolData(),
		resources.DSLayout:           resources.LayoutData(),
		resources.DSGroup:            resources.GroupData(),
		resources.DSPlan:             resources.PlanData(),
		resources.DSCloud:            resources.CloudData(),
		resources.DSResourcePool:     resources.ResourcePoolData(),
		resources.DSDatastore:        resources.DatastoreData(),
		resources.DSPowerSchedule:    resources.PowerScheduleData(),
		resources.DSTemplate:         resources.TemplateData(),
		resources.DSEnvironment:      resources.EnvironmentData(),
		resources.DSNetworkInterface: resources.NetworkInterfaceData(),
		resources.DSCloudFolder:      resources.CloudFolderData(),
		resources.DSRouter:           resources.RouterData(),
		resources.DSNetworkDomain:    resources.DomainData(),
		resources.DSNetworkProxy:     resources.NetworkProxyData(),
		resources.DSEdgeCluster:      resources.EdgeClusterData(),
		resources.DSTransportZone:    resources.TransportZoneData(),
		resources.DSLoadBalancer:     resources.LoadBalancerData(),
	}
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.ResInstance:                   resources.Instances(),
		resources.ResInstanceClone:              resources.InstancesClone(),
		resources.ResNetwork:                    resources.Network(),
		resources.ResRouter:                     resources.Router(),
		resources.ResRouterNat:                  resources.RouterNatRule(),
		resources.ResRouterFirewallRuleGroup:    resources.RouterFirewallRuleGroup(),
		resources.ResRouterRoute:                resources.RouterRoute(),
		resources.ResRouterBgpNeighbor:          resources.RouterBgpNeighbor(),
		resources.ResLoadBalancer:               resources.LoadBalancer(),
		resources.ResLoadBalancerMonitors:       resources.LoadBalancerMonitor(),
		resources.ResLoadBalancerProfiles:       resources.LoadBalancerProfiles(),
		resources.ResLoadBalancerPools:          resources.LoadBalancerPools(),
		resources.ResLoadBalancerVirtualServers: resources.LoadBalancerVirtualServers(),
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
			constants.APIURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_API_URL", constants.ServiceURL),
				Description: "The URL to use for the VMaaS API, can also be set with the HPEGL_VMAAS_API_URL env var",
			},
		},
	}
}
