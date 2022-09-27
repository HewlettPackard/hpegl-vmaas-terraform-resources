// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LoadBalancerPools() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent lb ID, lb_id can be obtained by using LB datasource/resource.",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer pool name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Creating the Network loadbalancer pool.",
			},
			"min_active_members": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The minimum number of members for the pool to be considered active",
			},
			"algorithm": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"ROUND_ROBIN",
					"WEIGHTED_ROUND_ROBIN",
					"LEAST_CONNECTION",
					"WEIGHTED_LEAST_CONNECTION",
					"IP_HASH",
				}, false),
				Required:     true,
				InputDefault: "ROUND_ROBIN",
				Description: "Load balancing pool algorithm controls how the incoming connections" +
					"are distributed among the members",
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "pool Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snat_translation_type": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"LBSnatAutoMap", "LBSnatDisabled", "LBSnatIpPool"}, false),
							Optional:         true,
							Default:          "LBSnatDisabled",
							Description:      "Network Loadbalancer Supported values are `LBSnatAutoMap`,`LBSnatDisabled`, `LBSnatIpPool`",
						},
						"passive_monitor_path": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "Passive Monitor ID, Get the `Id` from " + DSLBMonitor +
								"datasource to obtain the passive monitor ID",
						},
						"active_monitor_paths": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "Active Monitor ID, Get the `Id` from " + DSLBMonitor +
								"datasource to obtain the active monitor ID",
						},
						"tcp_multiplexing": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "With TCP multiplexing, user can use the same TCP connection" +
								"between a load balancer and the server for" +
								"sending multiple client requests from different client TCP connections.",
						},
						"tcp_multiplexing_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  6,
							Description: "The maximum number of TCP connections per pool" +
								"that are idly kept alive for sending future client requests",
						},
						"snat_ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Address of the snat_ip for Network loadbalancer pool",
						},
						"member_group": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "member group",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Pool Member Groups path, get the `externalId` from " + DSPoolMemeberGroup +
											"datasource to obtain the path",
									},
									"max_ip_list_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Description: "It Should only be specified if `limit_ip_list_size` is set to true." +
											"Limits the max number of pool members to the specified value",
									},
									"ip_revision_filter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Ip version filter is used to filter `IPv4` addresses from the grouping object",
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Description: "This is member port, The traffic which enter into VIP will get transfer" +
											"to member groups based on the port specified." +
											"Depends on the application running on the member VM",
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "tags Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag for Network Load balancer Profile",
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "scope for Network Load balancer Profile",
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerPoolReadContext,
		UpdateContext: loadbalancerPoolUpdateContext,
		CreateContext: loadbalancerPoolCreateContext,
		DeleteContext: loadbalancerPoolDeleteContext,
		Description: `loadbalancer Pool resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerPoolUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerPool.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerPoolReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerPool.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerPoolCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerPool.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerPoolReadContext(ctx, rd, meta)
}

func loadbalancerPoolDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerPool.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
