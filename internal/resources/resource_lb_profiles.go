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
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent lb ID, lb_id can be obtained by using LB datasource/resource.",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Profile Name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Creating the Network Load balancer Profile",
			},
			"service_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile", "LBClientSslProfile", "LBServerSslProfile", "LBCookiePersistenceProfile", "LBGenericPersistenceProfile", "LBSourceIpPersistenceProfile"}, false),
				Required:         true,
				Description:      "Network Loadbalancer Supported values are `LBHttpProfile`,`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`,`LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`,`LBSourceIpPersistenceProfile`"},
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
							Description: "Network Loadbalancer Supported values are `application-profile`, `ssl-profile`, `persistence-profile`",
						},
						"fast_tcp_idle_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          1800,
							ValidateDiagFunc: validations.IntBetween(1, 2147483647),
							Description:      "http_idle_timeout for Network Load balancer Profile",
						},
						"fast_udp_idle_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          300,
							ValidateDiagFunc: validations.IntBetween(1, 2147483647),
							Description:      "fast_udp_idle_timeout for Network Load balancer Profile",
						},
						"http_idle_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          15,
							ValidateDiagFunc: validations.IntBetween(1, 5400),
							Description:      "http_idle_timeout for Network Load balancer Profile",
						},
						"ha_flow_mirroring": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "ha_flow_mirroring for Network Load balancer Profile",
						},
						"connection_close_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          30,
							ValidateDiagFunc: validations.IntBetween(1, 60),
							Description:      "connection_close_timeout for Network Load balancer Profile",
						},
						"request_header_size": {
							Type:     schema.TypeInt,
							Optional: true,
							//Default:     1024,
							ValidateDiagFunc: validations.IntBetween(1, 65536),
							Description:      "request_header_size for Network Load balancer Profile",
						},
						"response_header_size": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validations.IntBetween(1, 65536),
							//Default:     4096,
							Description: "response_header_size for Network Load balancer Profile",
						},
						"redirection": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "redirection for Network Load balancer Profile",
						},
						"x_forwarded_for": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "x_forwarded_for for Network Load balancer Profile",
						},
						"request_body_size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request_body_size for Network Load balancer Profile",
						},
						"response_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validations.IntBetween(1, 2147483647),
							Description:      "response_timeout for Network Load balancer Profile",
						},
						"ntlm_authentication": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "ntlm_authentication for Network Load balancer Profile",
						},
						"share_persistence": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "ntlm_authentication for Network Load balancer Profile",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "cookie_name for Network Load balancer Profile",
						},
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
						"cookie_mode": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"INSERT", "PREFIX", "REWRITE"}, false),
							Optional:         true,
							Default:          "INSERT",
							Description:      "Network Loadbalancer Supported values are `INSERT`,`PREFIX`, `REWRITE`",
						},
						"cookie_type": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"LBPersistenceCookieTime", "LBSessionCookieTime"}, false),
							Optional:         true,
							Description:      "Network Loadbalancer Supported values are `LBPersistenceCookieTime`,`LBSessionCookieTime`",
						},
						"cookie_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "cookie_domain for Network Load balancer Profile",
						},
						"cookie_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "cookie_path for Network Load balancer Profile",
						},
						"max_idle_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "max_idle_time for Network Load balancer Profile",
						},
						"max_cookie_age": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "max_cookie_age for Network Load balancer Profile",
						},
						"ha_persistence_mirroring": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "ha_persistence_mirroring for Network Load balancer Profile",
						},
						"persistence_entry_timeout": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          300,
							ValidateDiagFunc: validations.IntBetween(1, 2147483647),
							Description:      "persistence_entry_timeout for Network Load balancer Profile",
						},
						"purge_entries_when_full": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "purge_entries_when_full for Network Load balancer Profile",
						},
						"ssl_suite": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"BALANCED", "HIGH_SECURITY", "HIGH_COMPATIBILITY", "CUSTOM"}, false),
							Optional:         true,
							Description:      "Network Loadbalancer Supported values are `BALANCED`,`HIGH_SECURITY`, `HIGH_COMPATIBILITY`,`CUSTOM`",
						},
						"session_cache": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "session_cache for Network Load balancer Profile",
						},
						"session_cache_entry_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "session_cache_entry_timeout for Network Load balancer Profile",
						},
						"prefer_server_cipher": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "prefer_server_cipher for Network Load balancer Profile",
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
				},
			},
		},
		ReadContext:   loadbalancerProfileReadContext,
		UpdateContext: loadbalancerProfileUpdateContext,
		CreateContext: loadbalancerProfileCreateContext,
		DeleteContext: loadbalancerProfileDeleteContext,
		Description: `loadbalancer Profile resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerProfileUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
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
