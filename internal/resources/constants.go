// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package resources

const (
	// datasource key
	DSNetwork          = "hpegl_vmaas_network"
	DSNetworkType      = "hpegl_vmaas_network_type"
	DSNetworkPool      = "hpegl_vmaas_network_pool"
	DSNetworkProxy     = "hpegl_vmaas_network_proxy"
	DSNetworkDomain    = "hpegl_vmaas_network_domain"
	DSLayout           = "hpegl_vmaas_layout"
	DSGroup            = "hpegl_vmaas_group"
	DSPlan             = "hpegl_vmaas_plan"
	DSCloud            = "hpegl_vmaas_cloud"
	DSResourcePool     = "hpegl_vmaas_resource_pool"
	DSDatastore        = "hpegl_vmaas_datastore"
	DSPowerSchedule    = "hpegl_vmaas_power_schedule"
	DSTemplate         = "hpegl_vmaas_template"
	DSEnvironment      = "hpegl_vmaas_environment"
	DSNetworkInterface = "hpegl_vmaas_network_interface"
	DSCloudFolder      = "hpegl_vmaas_cloud_folder"
	DSRouter           = "hpegl_vmaas_router"
	DSEdgeCluster      = "hpegl_vmaas_edge_cluster"
	DSTransportZone    = "hpegl_vmaas_transport_zone"
	DSLoadBalancer     = "hpegl_vmaas_load_balancer"
	DSLBMonitor        = "hpegl_vmaas_load_balancer_monitor"
	DSLBProfile        = "hpegl_vmaas_load_balancer_profile"
	DSLBPool           = "hpegl_vmaas_load_balancer_pool"
	DSLBVirtualServer  = "hpegl_vmaas_load_balancer_virtual_server"

	// resource key
	ResInstance                   = "hpegl_vmaas_instance"
	ResInstanceClone              = "hpegl_vmaas_instance_clone"
	ResNetwork                    = "hpegl_vmaas_network"
	ResRouter                     = "hpegl_vmaas_router"
	ResLoadBalancer               = "hpegl_vmaas_load_balancer"
	ResLoadBalancerMonitors       = "hpegl_vmaas_load_balancer_monitor"
	ResLoadBalancerProfiles       = "hpegl_vmaas_load_balancer_profile"
	ResLoadBalancerPools          = "hpegl_vmaas_load_balancer_pool"
	ResLoadBalancerVirtualServers = "hpegl_vmaas_load_balancer_virtual_server"
	ResRouterNat                  = "hpegl_vmaas_router_nat_rule"
	ResRouterFirewallRuleGroup    = "hpegl_vmaas_router_firewall_rule_group"
	ResRouterRoute                = "hpegl_vmaas_router_route"
	ResRouterBgpNeighbor          = "hpegl_vmaas_router_bgp_neighbor"

	// documentation related constants
	generalNamedesc = "Name of the %s as it appears on HPE GreenLake for private cloud dashboard. " +
		"If there is no %s with this name, a 'NOT FOUND' error will returned."
	generalDDesc = "Unique ID to identify a %s."
)
