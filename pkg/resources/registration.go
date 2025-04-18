// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hewlettpackard/hpegl-provider-lib/pkg/registration"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
)

// Assert that Registration implements the ServiceRegistration interface
var _ registration.ServiceRegistration = (*Registration)(nil)

type Registration struct{}

func (r Registration) Name() string {
	return constants.ServiceName
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.DSNetwork:                   resources.NetworkData(),
		resources.DSNetworkType:               resources.NetworkTypeData(),
		resources.DSNetworkPool:               resources.NetworkPoolData(),
		resources.DSLayout:                    resources.LayoutData(),
		resources.DSGroup:                     resources.GroupData(),
		resources.DSPlan:                      resources.PlanData(),
		resources.DSCloud:                     resources.CloudData(),
		resources.DSResourcePool:              resources.ResourcePoolData(),
		resources.DSDatastore:                 resources.DatastoreData(),
		resources.DSPowerSchedule:             resources.PowerScheduleData(),
		resources.DSTemplate:                  resources.TemplateData(),
		resources.DSEnvironment:               resources.EnvironmentData(),
		resources.DSNetworkInterface:          resources.NetworkInterfaceData(),
		resources.DSCloudFolder:               resources.CloudFolderData(),
		resources.DSRouter:                    resources.RouterData(),
		resources.DSNetworkDomain:             resources.DomainData(),
		resources.DSNetworkProxy:              resources.NetworkProxyData(),
		resources.DSEdgeCluster:               resources.EdgeClusterData(),
		resources.DSTransportZone:             resources.TransportZoneData(),
		resources.DSLBProfile:                 resources.LBProfileData(),
		resources.DSLBPool:                    resources.LBPoolData(),
		resources.DSLoadBalancer:              resources.LoadBalancerData(),
		resources.DSLBMonitor:                 resources.MonitorData(),
		resources.DSPoolMemeberGroup:          resources.LBPoolMemeberGroupData(),
		resources.DSLBVirtualServerSslCert:    resources.LBVirtualServerSslCertData(),
		resources.DSDhcpServer:                resources.DhcpServerData(),
		resources.DSInstanceStorageType:       resources.ReadInstanceStorageType(),
		resources.DSInstanceStorageController: resources.ReadInstanceStorageController(),
		resources.DSMorpheusDataSource:        resources.MorpheusDetailsBroker(),
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
		resources.ResDhcpServer:                 resources.DhcpServer(),
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
				Description: "It can also be set with the HPEGL_VMAAS_SPACE_NAME env var. When `HPEGL_IAM_VERSION` is `glcs` it refers to IAM Space name of the GL VMaaS Service i.e., Default. When `HPEGL_IAM_VERSION` is `glp` it refers to GLP Workspace ID.",
			},
			constants.BROKERRURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_VMAAS_BROKER_URL", constants.BrokerURL),
				Description: "The URL to use for the VMaaS Broker API, can also be set with the HPEGL_VMAAS_BROKER_URL env var",
			},
			constants.MORPHEUS_URL: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_MORPHEUS_URL", ""),
				Description: "The Morpheus URL, can also be set with the HPEGL_MORPHEUS_URL env var",
			},
			constants.MORPHEUS_TOKEN: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HPEGL_MORPHEUS_TOKEN", ""),
				Description: "The Morpheus token, can also be set with the HPEGL_MORPHEUS_TOKEN env var",
			},
			constants.INSECURE: {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INSECURE", false),
				Description: "Not to be used in production. To perform client connection ignoring TLS, it can also be set with the INSECURE env var",
			},
		},
	}
}
