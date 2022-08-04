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

func LoadBalancerProfiles() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Profile Name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Creating the Network Load balancer Profile",
				ForceNew:    true,
			},
			"service_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile", "LBClientSslProfile", "LBServerSslProfile", "LBCookiePersistenceProfile", "LBGenericPersistenceProfile"}, false),
				Required:         true,
				Description:      "Network Loadbalancer Supported values are `LBHttpProfile`,`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`,`LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`"},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "profile Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile_type": {
							Type: schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{
								"application-profile", "ssl-profile", "persistence-profile",
							}, false),
							Required:    true,
							Description: "Network Loadbalancer Supported values are `application-profile`, `ssl-profile`, `persistence-profile`"},
						"request_header_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "request_header_size for Network Load balancer Profile",
						},
						"response_header_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "response_header_size for Network Load balancer Profile",
						},
						"http_idle_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "http_idle_timeout for Network Load balancer Profile",
						},
						"fast_tcp_idle_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "fast_tcp_idle_timeout for Network Load balancer Profile",
						},
						"connection_close_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "connection_close_timeout for Network Load balancer Profile",
						},
						"ha_flow_mirroring": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "ha_flow_mirroring for Network Load balancer Profile",
						},
						"ssl_suite": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"BALANCED", "HIGH_SECURITY", "HIGH_COMPATIBILITY", "CUSTOM"}, false),
							Optional:         true,
							Description:      "Network Loadbalancer Supported values are `BALANCED`,`HIGH_SECURITY`, `HIGH_COMPATIBILITY`,`CUSTOM`",
						},
						"cookie_mode": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"INSERT", "PREFIX", "REWRITE"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `INSERT`,`PREFIX`, `REWRITE`",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cookie_name for Network Load balancer Profile",
						},
						"cookie_type": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"LBPersistenceCookieTime", "LBSessionCookieTime"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `LBPersistenceCookieTime`,`LBSessionCookieTime`"},
						"cookie_fallback": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "cookie_fallback for Network Load balancer Profile",
						},
						"cookie_garbling": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "cookie_garbling for Network Load balancer Profile",
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerProfileReadContext,
		UpdateContext: loadbalancerProfileReadContext,
		CreateContext: loadbalancerProfileCreateContext,
		DeleteContext: loadbalancerProfileDeleteContext,
		Description: `loadbalancer Profile resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerProfileReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerProfileCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerProfileReadContext(ctx, rd, meta)
}

func loadbalancerProfileDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
