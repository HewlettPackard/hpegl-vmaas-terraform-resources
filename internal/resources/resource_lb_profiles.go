// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
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
			"serviceType": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile",
					 "LBClientSslProfile","LBServerSslProfile",
					 "LBCookiePersistenceProfile","LBGenericPersistenceProfile"
				}, false),
				Required:    true,
				Description: "Network Loadbalancer Supported values are `LBHttpProfile`, 
				`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`,
				 `LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "profile Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profileType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "application-profile", "ssl-profile", "persistence-profile",
				            }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `application-profile`, 
				            `ssl-profile`, `persistence-profile`",
						},
						"requestHeaderSize": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "requestHeaderSize for Network Load balancer Profile",
						},
						"responseHeaderSize": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "responseHeaderSize for Network Load balancer Profile",
						},
						"httpIdleTimeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "httpIdleTimeout for Network Load balancer Profile",
						},
						"fastTcpIdleTimeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "fastTcpIdleTimeout for Network Load balancer Profile",
						},
						"connectionCloseTimeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "connectionCloseTimeout for Network Load balancer Profile",
						},
						"haFlowMirroring": {
							Type:        schema.TypeBool,
							Required:    true,
							Default:     true,
							Description: "haFlowMirroring for Network Load balancer Profile",
						},
						"sslSuite": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
				        	    "BALANCED", "HIGH_SECURITY", "HIGH_COMPATIBILITY", "CUSTOM"
				             }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `BALANCED`, 
				                `HIGH_SECURITY`, `HIGH_COMPATIBILITY`,`CUSTOM`",
						},
						"cookieMode": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "INSERT", "PREFIX", "REWRITE"
				            }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `INSERT`, 
				                `PREFIX`, `REWRITE`",
						},
						"cookieName": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cookieName for Network Load balancer Profile",
						},
						"cookieType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "LBPersistenceCookieTime", "LBSessionCookieTime"
				            }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `LBPersistenceCookieTime`, 
				            `LBSessionCookieTime`",
						},
						"cookieFallback": {
							Type:        schema.TypeBool,
							Required:    true,
							Default:     true,
							Description: "cookieFallback for Network Load balancer Profile",
						},
						"cookieGarbling": {
							Type:        schema.TypeBool,
							Required:    true,
							Default:     true,
							Description: "cookieGarbling for Network Load balancer Profile",
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerProfileReadContext,
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
	if err := c.CmpClient.ResLoadBalancerProfiles.Read(ctx, data, meta); err != nil {
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
	if err := c.CmpClient.ResLoadBalancerProfiles.Create(ctx, data, meta); err != nil {
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
	if err := c.CmpClient.ResLoadBalancerProfiles.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
